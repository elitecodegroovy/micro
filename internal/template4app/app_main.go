package template4app

var (
	AppMain = `
package main

import (
	"flag"
	"fmt"
	"{{.Dir}}/pkg/appserver"
	"{{.Dir}}/pkg/infra/log"
	"{{.Dir}}/pkg/setting"
	_ "github.com/go-sql-driver/mysql"
	sLog "log"
	"net/http"
	"os"
	"runtime"
	"runtime/trace"
	"strconv"
	"time"
)

func validPackaging(packaging string) string {
	validTypes := []string{"dev", "deb", "rpm", "docker", "brew", "hosted", "unknown"}
	for _, vt := range validTypes {
		if packaging == vt {
			return packaging
		}
	}
	return "unknown"
}

func main() {
	sLog.SetOutput(os.Stdout)
	sLog.SetFlags(0)

	v := flag.Bool("v", false, "prints current version and exits")
	profile := flag.Bool("profile", false, "Turn on pprof profiling")
	profilePort := flag.Int("profile-port", 6060, "Define custom port for profiling")
	flag.Parse()
	if *v {
		fmt.Printf("Version %s (commit: %s, branch: %s)\n", appserver.Version, appserver.Commit, appserver.BuildBranch)
		os.Exit(0)
	}

	if *profile {
		runtime.SetBlockProfileRate(1)
		go func() {
			err := http.ListenAndServe(fmt.Sprintf("localhost:%d", *profilePort), nil)
			if err != nil {
				panic(err)
			}
		}()

		f, err := os.Create("trace.out")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		err = trace.Start(f)
		if err != nil {
			panic(err)
		}
		defer trace.Stop()
	}

	buildstampInt64, _ := strconv.ParseInt(appserver.Buildstamp, 10, 64)
	if buildstampInt64 == 0 {
		buildstampInt64 = time.Now().Unix()
	}

	setting.BuildVersion = appserver.Version
	setting.BuildCommit = appserver.Commit
	setting.BuildStamp = buildstampInt64
	setting.BuildBranch = appserver.BuildBranch
	setting.Packaging = validPackaging(*appserver.Packaging)
	sLog.Printf("Version: %s, Commit Version: %s, Package Iteration: %s\n", appserver.Version, setting.BuildCommit, setting.BuildBranch)

	server := appserver.NewGNetworkServer()

	go appserver.ListenToSystemSignals(server)

	err := server.Run()

	code := server.Exit(err)
	trace.Stop()
	log.Close()

	os.Exit(code)

}

`
)
