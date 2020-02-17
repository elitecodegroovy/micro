package template4app

var (
	AppServer = `
package appserver

import (
	"context"
	"flag"
	"fmt"
	"{{.Dir}}/pkg/bus"
	"{{.Dir}}/pkg/infra/localcache"
	"{{.Dir}}/pkg/infra/log"
	"{{.Dir}}/pkg/middleware"
	"{{.Dir}}/pkg/registry"
	"{{.Dir}}/pkg/routing"
	"{{.Dir}}/pkg/server"
	_ "{{.Dir}}/pkg/server"
	_ "{{.Dir}}/pkg/services/auth"
	_ "{{.Dir}}/pkg/services/login"
	_ "{{.Dir}}/pkg/services/notifications"
	_ "{{.Dir}}/pkg/services/quota"
	_ "{{.Dir}}/pkg/services/serverlock"
	_ "{{.Dir}}/pkg/services/sqlstore"
	"{{.Dir}}/pkg/setting"
	"github.com/facebookgo/inject"
	"golang.org/x/sync/errgroup"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

var Version = "1.0.0"
var Commit = "NA"
var BuildBranch = "master"
var Buildstamp string

var ConfigFile = flag.String("config", "", "path to config file")
var HomePath = flag.String("homepath", "", "path to gnetwork install/home path, defaults to working directory")
var PidFile = flag.String("pidfile", "", "path to pid file")
var Packaging = flag.String("packaging", "unknown", "describes the way gnetwork was installed")

type AppServerImpl struct {
	context            context.Context
	shutdownFn         context.CancelFunc
	childRoutines      *errgroup.Group
	log                log.Logger
	cfg                *setting.Cfg
	shutdownReason     string
	shutdownInProgress bool

	RouteRegister routing.RouteRegister ` + "`inject:\"\"`" + `
	HttpServer    *server.HTTPServer    `+  "`inject:\"\"`" + `
}

func NewGNetworkServer() *AppServerImpl {
	//a copy of parent
	rootCtx, shutdownFn := context.WithCancel(context.Background())
	//a new Group
	childRoutines, childCtx := errgroup.WithContext(rootCtx)

	return &AppServerImpl{
		context:       childCtx,
		shutdownFn:    shutdownFn,
		childRoutines: childRoutines,
		log:           log.New("server"),
		cfg:           setting.NewCfg(),
	}
}

func (g *AppServerImpl) Shutdown(reason string) {
	g.log.Info("Shutdown started", "reason", reason)
	g.shutdownReason = reason
	g.shutdownInProgress = true

	// call cancel func on root context
	g.shutdownFn()

	// wait for child routines
	g.childRoutines.Wait()
}

func (g *AppServerImpl) Exit(reason error) int {
	// default exit code is 1
	code := 1

	if reason == context.Canceled && g.shutdownReason != "" {
		reason = fmt.Errorf(g.shutdownReason)
		code = 0
	}

	g.log.Error("Server shutdown", "reason", reason)
	return code
}

func (g *AppServerImpl) loadConfiguration() {
	err := g.cfg.Load(&setting.CommandLineArgs{
		Config:   *ConfigFile,
		HomePath: *HomePath,
		Args:     flag.Args(),
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start grafana. error: %s\n", err.Error())
		os.Exit(1)
	}

	g.log.Info("Starting "+setting.ApplicationName, "version", Version, "commit", Commit, "branch", BuildBranch, "compiled", time.Unix(setting.BuildStamp, 0))
	g.cfg.LogConfigSources()
	g.log.Info("loading configuration successfully!")
}

func (g *AppServerImpl) writePIDFile() {
	if *PidFile == "" {
		g.log.Info("path is empty for the pid file! ")
		return
	}

	// Ensure the required directory structure exists.
	err := os.MkdirAll(filepath.Dir(*PidFile), 0700)
	if err != nil {
		g.log.Error("Failed to verify pid directory", "error", err)
		os.Exit(1)
	}

	// Retrieve the PID and write it.
	pid := strconv.Itoa(os.Getpid())
	if err := ioutil.WriteFile(*PidFile, []byte(pid), 0644); err != nil {
		g.log.Error("Failed to write pidfile", "error", err)
		os.Exit(1)
	}

	g.log.Info("Writing PID file", "path", *PidFile, "pid", pid)
}

func (g *AppServerImpl) Run() error {
	var err error
	g.loadConfiguration()
	g.writePIDFile()

	serviceGraph := inject.Graph{}
	err = serviceGraph.Provide(&inject.Object{Value: bus.GetBus()})
	if err != nil {
		return fmt.Errorf("failed to provide object to the graph: %v", err)
	}
	err = serviceGraph.Provide(&inject.Object{Value: g.cfg})
	if err != nil {
		return fmt.Errorf("failed to provide object to the graph: %v", err)
	}
	//err = serviceGraph.Provide(&inject.Object{Value: routing.NewRouteRegister(middleware.RequestMetrics, middleware.RequestTracing)})
	err = serviceGraph.Provide(&inject.Object{Value: routing.NewRouteRegister(middleware.RequestMetrics)})
	if err != nil {
		return fmt.Errorf("failed to provide object to the graph: %v", err)
	}
	err = serviceGraph.Provide(&inject.Object{Value: localcache.New(5*time.Minute, 10*time.Minute)})
	if err != nil {
		return fmt.Errorf("failed to provide object to the graph: %v", err)
	}

	// self registered services
	services := registry.GetServices()

	// Add all services to dependency graph
	for _, service := range services {
		err = serviceGraph.Provide(&inject.Object{Value: service.Instance})
		if err != nil {
			return fmt.Errorf("failed to provide object to the graph: %v", err)
		}
	}

	err = serviceGraph.Provide(&inject.Object{Value: g})
	if err != nil {
		return fmt.Errorf("failed to provide object to the graph: %v", err)
	}

	// Inject dependencies to services
	if err := serviceGraph.Populate(); err != nil {
		return fmt.Errorf("failed to populate service dependency: %v", err)
	}

	// Init & start services
	for _, service := range services {
		if registry.IsDisabled(service.Instance) {
			continue
		}

		g.log.Info("Initializing " + service.Name)

		if err := service.Instance.Init(); err != nil {
			return fmt.Errorf("service init failed: %v", err)
		}
	}

	// Start background services
	for _, srv := range services {
		// variable needed for accessing loop variable in function callback
		descriptor := srv
		service, ok := srv.Instance.(registry.BackgroundService)
		if !ok {
			continue
		}

		if registry.IsDisabled(descriptor.Instance) {
			continue
		}

		g.childRoutines.Go(func() error {
			// Skip starting new service when shutting down
			// Can happen when service stop/return during startup
			if g.shutdownInProgress {
				return nil
			}

			err := service.Run(g.context)

			// If error is not canceled then the service crashed
			if err != context.Canceled && err != nil {
				g.log.Error("Stopped "+descriptor.Name, "reason", err)
			} else {
				g.log.Info("Stopped "+descriptor.Name, "reason", err)
			}

			// Mark that we are in shutdown mode
			// So more services are not started
			g.shutdownInProgress = true
			return err
		})
	}

	sendSystemdNotification("READY=1")

	return g.childRoutines.Wait()
}

//Setting system signal handling
func ListenToSystemSignals(server *AppServerImpl) {
	signalChan := make(chan os.Signal, 1)
	sighupChan := make(chan os.Signal, 1)

	signal.Notify(sighupChan, syscall.SIGHUP)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-sighupChan:
			log.Reload()
		case sig := <-signalChan:
			server.Shutdown(fmt.Sprintf("System signal: %s", sig))
		}
	}
}

func sendSystemdNotification(state string) error {
	notifySocket := os.Getenv("NOTIFY_SOCKET")

	if notifySocket == "" {
		return fmt.Errorf("NOTIFY_SOCKET environment variable empty or unset")
	}

	socketAddr := &net.UnixAddr{
		Name: notifySocket,
		Net:  "unixgram",
	}

	conn, err := net.DialUnix(socketAddr.Net, nil, socketAddr)

	if err != nil {
		return err
	}

	_, err = conn.Write([]byte(state))

	conn.Close()

	return err
}

`
)
