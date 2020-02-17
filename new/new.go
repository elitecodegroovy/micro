// Package new generates micro service templates
package new

import (
	"fmt"
	"go/build"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
	"time"

	"github.com/micro/cli"
	tmpl "github.com/micro/micro/internal/template"
	apptmpl "github.com/micro/micro/internal/template4app"
	"github.com/micro/micro/internal/usage"
	"github.com/xlab/treeprint"
)

type config struct {
	// foo
	Alias string
	// micro new example -type
	Command string
	// go.micro
	Namespace string
	// api, srv, web, fnc
	Type string
	// go.micro.srv.foo
	FQDN string
	// github.com/micro/foo
	Dir string
	// $GOPATH/src/github.com/micro/foo
	GoDir string
	// $GOPATH
	GoPath string
	// Files
	Files []file
	// Comments
	Comments []string
	// Plugins registry=etcd:broker=nats
	Plugins []string
}

type file struct {
	Path string
	Tmpl string
}

func write(c config, file, tmpl string) error {
	fn := template.FuncMap{
		"title": strings.Title,
	}

	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	t, err := template.New("f").Funcs(fn).Parse(tmpl)
	if err != nil {
		return err
	}

	return t.Execute(f, c)
}

func create(c config) error {
	// check if dir exists
	if _, err := os.Stat(c.GoDir); !os.IsNotExist(err) {
		return fmt.Errorf("%s already exists", c.GoDir)
	}

	// create usage report
	u := usage.New("new")
	// a single request/service
	u.Metrics.Count["requests"] = uint64(1)
	u.Metrics.Count["services"] = uint64(1)
	// send report
	go usage.Report(u)

	// just wait
	<-time.After(time.Millisecond * 250)

	fmt.Printf("Creating service %s in %s\n\n", c.FQDN, c.GoDir)

	t := treeprint.New()

	nodes := map[string]treeprint.Tree{}
	nodes[c.GoDir] = t

	// write the files
	for _, file := range c.Files {
		f := filepath.Join(c.GoDir, file.Path)
		dir := filepath.Dir(f)

		b, ok := nodes[dir]
		if !ok {
			d, _ := filepath.Rel(c.GoDir, dir)
			b = t.AddBranch(d)
			nodes[dir] = b
		}

		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return err
			}
		}

		p := filepath.Base(f)

		b.AddNode(p)
		if err := write(c, f, file.Tmpl); err != nil {
			return err
		}
	}

	// print tree
	fmt.Println(t.String())

	for _, comment := range c.Comments {
		fmt.Println(comment)
	}

	// just wait
	<-time.After(time.Millisecond * 250)

	return nil
}

func run(ctx *cli.Context) {
	namespace := ctx.String("namespace")
	alias := ctx.String("alias")
	fqdn := ctx.String("fqdn")
	atype := ctx.String("type")
	dir := ctx.Args().First()
	useGoPath := ctx.Bool("gopath")
	useGoModule := os.Getenv("GO111MODULE")
	var plugins []string

	if len(dir) == 0 {
		fmt.Println("specify service name")
		return
	}

	if len(namespace) == 0 {
		fmt.Println("namespace not defined")
		return
	}

	if len(atype) == 0 {
		fmt.Println("type not defined")
		return
	}

	// set the command
	command := fmt.Sprintf("micro new %s", dir)
	if len(namespace) > 0 {
		command += " --namespace=" + namespace
	}
	if len(alias) > 0 {
		command += " --alias=" + alias
	}
	if len(fqdn) > 0 {
		command += " --fqdn=" + fqdn
	}
	if len(atype) > 0 {
		command += " --type=" + atype
	}
	if plugins := ctx.StringSlice("plugin"); len(plugins) > 0 {
		command += " --plugin=" + strings.Join(plugins, ":")
	}

	// check if the path is absolute, we don't want this
	// we want to a relative path so we can install in GOPATH
	if path.IsAbs(dir) {
		fmt.Println("require relative path as service will be installed in GOPATH")
		return
	}

	var goPath string
	var goDir string

	// only set gopath if told to use it
	if useGoPath {
		goPath = build.Default.GOPATH

		// don't know GOPATH, runaway....
		if len(goPath) == 0 {
			fmt.Println("unknown GOPATH")
			return
		}

		// attempt to split path if not windows
		if runtime.GOOS == "windows" {
			goPath = strings.Split(goPath, ";")[0]
		} else {
			goPath = strings.Split(goPath, ":")[0]
		}
		goDir = filepath.Join(goPath, "src", path.Clean(dir))
	} else {
		goDir = path.Clean(dir)
	}

	if len(alias) == 0 {
		// set as last part
		alias = filepath.Base(dir)
		// strip hyphens
		parts := strings.Split(alias, "-")
		alias = parts[0]
	}

	if len(fqdn) == 0 {
		fqdn = strings.Join([]string{namespace, atype, alias}, ".")
	}

	for _, plugin := range ctx.StringSlice("plugin") {
		// registry=etcd:broker=nats
		for _, p := range strings.Split(plugin, ":") {
			// registry=etcd
			parts := strings.Split(p, "=")
			if len(parts) < 2 {
				continue
			}
			plugins = append(plugins, path.Join(parts...))
		}
	}

	var c config

	switch atype {
	case "fnc":
		// create srv config
		c = config{
			Alias:     alias,
			Command:   command,
			Namespace: namespace,
			Type:      atype,
			FQDN:      fqdn,
			Dir:       dir,
			GoDir:     goDir,
			GoPath:    goPath,
			Plugins:   plugins,
			Files: []file{
				{"main.go", tmpl.MainFNC},
				{"generate.go", tmpl.GenerateFile},
				{"plugin.go", tmpl.Plugin},
				{"handler/" + alias + ".go", tmpl.HandlerFNC},
				{"subscriber/" + alias + ".go", tmpl.SubscriberFNC},
				{"proto/" + alias + "/" + alias + ".proto", tmpl.ProtoFNC},
				{"Dockerfile", tmpl.DockerFNC},
				{"Makefile", tmpl.Makefile},
				{"README.md", tmpl.ReadmeFNC},
			},
			Comments: []string{
				"\ndownload protobuf for micro:\n",
				"brew install protobuf",
				"go get -u github.com/golang/protobuf/{proto,protoc-gen-go}",
				"go get -u github.com/micro/protoc-gen-micro",
				"\ncompile the proto file " + alias + ".proto:\n",
				"cd " + goDir,
				"protoc --proto_path=.:$GOPATH/src --go_out=. --micro_out=. proto/" + alias + "/" + alias + ".proto\n",
			},
		}
	case "srv":
		// create srv config
		c = config{
			Alias:     alias,
			Command:   command,
			Namespace: namespace,
			Type:      atype,
			FQDN:      fqdn,
			Dir:       dir,
			GoDir:     goDir,
			GoPath:    goPath,
			Plugins:   plugins,
			Files: []file{
				{"main.go", tmpl.MainSRV},
				{"generate.go", tmpl.GenerateFile},
				{"plugin.go", tmpl.Plugin},
				{"handler/" + alias + ".go", tmpl.HandlerSRV},
				{"subscriber/" + alias + ".go", tmpl.SubscriberSRV},
				{"proto/" + alias + "/" + alias + ".proto", tmpl.ProtoSRV},
				{"Dockerfile", tmpl.DockerSRV},
				{"Makefile", tmpl.Makefile},
				{"README.md", tmpl.Readme},
			},
			Comments: []string{
				"\ndownload protobuf for micro:\n",
				"brew install protobuf",
				"go get -u github.com/golang/protobuf/{proto,protoc-gen-go}",
				"go get -u github.com/micro/protoc-gen-micro",
				"\ncompile the proto file " + alias + ".proto:\n",
				"cd " + goDir,
				"protoc --proto_path=.:$GOPATH/src --go_out=. --micro_out=. proto/" + alias + "/" + alias + ".proto\n",
			},
		}
	case "api":
		// create api config
		c = config{
			Alias:     alias,
			Command:   command,
			Namespace: namespace,
			Type:      atype,
			FQDN:      fqdn,
			Dir:       dir,
			GoDir:     goDir,
			GoPath:    goPath,
			Plugins:   plugins,
			Files: []file{
				{"main.go", tmpl.MainAPI},
				{"generate.go", tmpl.GenerateFile},
				{"plugin.go", tmpl.Plugin},
				{"client/" + alias + ".go", tmpl.WrapperAPI},
				{"handler/" + alias + ".go", tmpl.HandlerAPI},
				{"proto/" + alias + "/" + alias + ".proto", tmpl.ProtoAPI},
				{"Makefile", tmpl.Makefile},
				{"Dockerfile", tmpl.DockerSRV},
				{"README.md", tmpl.Readme},
			},
			Comments: []string{
				"\ndownload protobuf for micro:\n",
				"brew install protobuf",
				"go get -u github.com/golang/protobuf/{proto,protoc-gen-go}",
				"go get -u github.com/micro/protoc-gen-micro",
				"\ncompile the proto file " + alias + ".proto:\n",
				"cd " + goDir,
				"protoc --proto_path=.:$GOPATH/src --go_out=. --micro_out=. proto/" + alias + "/" + alias + ".proto\n",
			},
		}
	case "web":
		// create srv config
		c = config{
			Alias:     alias,
			Command:   command,
			Namespace: namespace,
			Type:      atype,
			FQDN:      fqdn,
			Dir:       dir,
			GoDir:     goDir,
			GoPath:    goPath,
			Plugins:   plugins,
			Files: []file{
				{"main.go", tmpl.MainWEB},
				{"plugin.go", tmpl.Plugin},
				{"handler/handler.go", tmpl.HandlerWEB},
				{"html/index.html", tmpl.HTMLWEB},
				{"Dockerfile", tmpl.DockerWEB},
				{"Makefile", tmpl.Makefile},
				{"README.md", tmpl.Readme},
			},
			Comments: []string{},
		}
	case "app":
		// create srv config
		c = config{
			Alias:     alias,
			Command:   command,
			Namespace: namespace,
			Type:      atype,
			FQDN:      fqdn,
			Dir:       dir,
			GoDir:     goDir,
			GoPath:    goPath,
			Plugins:   plugins,
			Files: []file{
				{"README.md", apptmpl.Readme},
				{"conf/config.ini", apptmpl.Config},
				{"pkg/appserver/app_server.go", apptmpl.AppServer},
				{"pkg/bus/bus.go", apptmpl.Bus},
				{"pkg/bus/bus_test.go", apptmpl.BusTest},
				{"pkg/cmd/" + alias +"/main.go", apptmpl.AppMain},
				{"pkg/components/apikeygen/apikeygen.go", apptmpl.Apikeygen},
				{"pkg/components/apikeygen/apikeygen_test.go", apptmpl.ApikeygenTest},
				{"pkg/components/simplejson/simplejson.go", apptmpl.SimpleJson},
				{"pkg/components/simplejson/simplejson_go11.go", apptmpl.SimpleJsonGo11},
				{"pkg/components/simplejson/simplejson_test.go", apptmpl.SimpleJsonTest},
				{"pkg/events/events.go", apptmpl.Events},
				{"pkg/events/events_test.go", apptmpl.EventsTest},
				{"pkg/infra/localcache/cache.go", apptmpl.Cache},
				{"pkg/infra/log/file.go", apptmpl.LogFile},
				{"pkg/infra/log/file_test.go", apptmpl.LogFileTest},
				{"pkg/infra/log/handlers.go", apptmpl.LogHandler},
				{"pkg/infra/log/interface.go", apptmpl.LogInterface},
				{"pkg/infra/log/log.go", apptmpl.LogLog},
				{"pkg/infra/log/Log_writer.go", apptmpl.LogLogWriter},
				{"pkg/infra/log/Log_writer_test.go", apptmpl.LogLogWriterTest},
				{"pkg/infra/log/syslog.go", apptmpl.LogSyslog},
				{"pkg/infra/log/syslog_windows.go", apptmpl.LogSyslogWindows},
				{"pkg/infra/metrics/metric.go", apptmpl.Metric},
				{"pkg/middleware/auth.go", apptmpl.Auth},
				{"pkg/middleware/gzipper.go", apptmpl.Gzipper},
				{"pkg/middleware/middlerware.go", apptmpl.Middleware},
				{"pkg/middleware/logger.go", apptmpl.Logger},
				{"pkg/middleware/header.go", apptmpl.Header},
				{"pkg/middleware/quota.go", apptmpl.Quota},
				{"pkg/middleware/recovery.go", apptmpl.Recovery},
				{"pkg/middleware/request_metric.go", apptmpl.RequestMetric},
				{"pkg/middleware/request_tracing.go", apptmpl.RequestTracing},
				{"pkg/models/apikey.go", apptmpl.ModelApikey},
				{"pkg/models/context.go", apptmpl.ModelContext},
				{"pkg/models/health.go", apptmpl.ModelHealth},
				{"pkg/models/helpFlags.go", apptmpl.ModelHelpFlags},
				{"pkg/models/login_attempt.go", apptmpl.ModelLoginAttempt},
				{"pkg/models/notifications.go", apptmpl.ModelNotifications},
				{"pkg/models/org.go", apptmpl.ModelOrg},
				{"pkg/models/org_user.go", apptmpl.ModelOrgUser},
				{"pkg/models/quotas.go", apptmpl.ModelQuotas},
				{"pkg/models/tags.go", apptmpl.ModelTags},
				{"pkg/models/tags_test.go", apptmpl.ModelTagsTest},
				{"pkg/models/team.go", apptmpl.ModelTeam},
				{"pkg/models/team_member.go", apptmpl.ModelTeamMember},
				{"pkg/models/tempt_user.go", apptmpl.ModelTempUser},
				{"pkg/models/user.go", apptmpl.ModelUser},
				{"pkg/models/user_auth.go", apptmpl.ModelUserAuth},
				{"pkg/models/user_token.go", apptmpl.ModelUserToken},
				{"pkg/registry/registry.go", apptmpl.Registry},
				{"pkg/routing/router_register.go", apptmpl.RouterRegister},
				{"pkg/routing/router_register_test.go", apptmpl.RouterRegisterTest},
				{"pkg/server/dtos/models.go", apptmpl.DtosModels},
				{"pkg/server/dtos/user.go", apptmpl.DtosUser},
				{"pkg/server/static/static.go", apptmpl.ServerStatic},
				{"pkg/server/admin_user.go", apptmpl.ServerAdminUser},
				{"pkg/server/admin_user.go", apptmpl.ServerAdminUser},
				{"pkg/server/common.go", apptmpl.ServerCommon},
				{"pkg/server/http_server.go", apptmpl.ServerHttpServer},
				{"pkg/server/index.go", apptmpl.ServerIndex},
				{"pkg/server/login.go", apptmpl.ServerLogin},
				{"pkg/server/url_mapping.go", apptmpl.ServerUrlMapping},
				{"pkg/services/annotations/annotations.go", apptmpl.ServiceAnnotations},
				{"pkg/services/auth/auth_token.go", apptmpl.ServiceAuthToken},
				{"pkg/services/auth/auth_token_test.go", apptmpl.ServiceAuthTokenTest},
				{"pkg/services/auth/model.go", apptmpl.ServiceModel},
				{"pkg/services/auth/testing.go", apptmpl.ServiceTesting},
				{"pkg/services/auth/token_cleanup.go", apptmpl.ServiceTokenCleanup},
				{"pkg/services/auth/token_cleanup_test.go", apptmpl.ServiceTokenCleanupTest},
				{"pkg/services/login/auth.go", apptmpl.ServiceLoginAuth},
				{"pkg/services/login/brute_force_login_protection.go", apptmpl.ServiceLoginBruceForceLoginProtection},
				{"pkg/services/login/brute_force_login_protection_test.go", apptmpl.ServiceLoginBruceForceLoginProtectionTest},
				{"pkg/services/login/grafana_login.go", apptmpl.ServiceLoginGrafanaLogin},
				{"pkg/services/notifications/codes.go", apptmpl.ServiceNotificationCodes},
				{"pkg/services/notifications/codes_test.go", apptmpl.ServiceNotificationCodesTest},
				{"pkg/services/notifications/email.go", apptmpl.ServiceNotificationEmail},
				{"pkg/services/notifications/mailer.go", apptmpl.ServiceNotificationMailer},
				{"pkg/services/notifications/notifications.go", apptmpl.ServiceNotifications},
				{"pkg/services/notifications/notifications_test.go", apptmpl.ServiceNotificationTest},
				{"pkg/services/notifications/send_email_integration_test.go", apptmpl.ServiceNotificationSendEmailIntegrationTest},
				{"pkg/services/notifications/webhook.go", apptmpl.ServiceNotificationWebhook},
				{"pkg/services/quota/quota.go", apptmpl.ServiceQuota},
				{"pkg/services/remotecache/database_storage.go", apptmpl.ServiceRemotecacheDatabaseStorage},
				{"pkg/services/remotecache/database_storage_test.go", apptmpl.ServiceRemotecacheDatabaseStorageTest},
				{"pkg/services/remotecache/memcache_storage.go", apptmpl.ServiceRemotecacheMemcacheStorage},
				{"pkg/services/remotecache/memcache_storage_test.go", apptmpl.ServiceRemotecacheMemcacheStorageTest},
				{"pkg/services/remotecache/redis_storage.go", apptmpl.ServiceRemotecacheRedisStorage},
				{"pkg/services/remotecache/redis_storage_integratioin_test.go", apptmpl.ServiceRemotecacheRedisStorageIntegrationTest},
				{"pkg/services/remotecache/redis_storage_test.go", apptmpl.ServiceRemotecacheRedisStorageTest},
				{"pkg/services/remotecache/remotecache.go", apptmpl.ServiceRemotecache},
				{"pkg/services/remotecache/remotecache_test.go", apptmpl.ServiceRemotecacheTest},
				{"pkg/services/remotecache/testing.go", apptmpl.ServiceRemotecacheTesting},
				{"pkg/services/serverlock/model.go", apptmpl.ServiceServerlockModel},
				{"pkg/services/serverlock/serverlock.go", apptmpl.ServiceServerlock},
				{"pkg/services/serverlock/serverlock_test.go", apptmpl.ServiceServerlockTest},
				{"pkg/services/serverlock/serverlock_integration_test.go", apptmpl.ServiceServerlockIntegrationTest},
				{"pkg/services/sqlstore/migrator/column.go", apptmpl.ServiceSqlstoreMigratorColumn},
				{"pkg/services/sqlstore/migrator/conditions.go", apptmpl.ServiceSqlStoreMigratorConditions},
				{"pkg/services/sqlstore/migrator/dialect.go", apptmpl.ServiceSqlstoreMigratorDialect},
				{"pkg/services/sqlstore/migrator/migrator.go", apptmpl.ServiceSqlstoreMigrator},
				{"pkg/services/sqlstore/migrator/mysql_dialect.go", apptmpl.ServiceSqlstoreMigratorMySQLDialect},
				{"pkg/services/sqlstore/migrator/sqlite_dialect.go", apptmpl.ServiceSqlstoreMigratorSqliteDialect},
				{"pkg/services/sqlstore/migrator/types.go", apptmpl.ServiceSqlstoreMigratorTypes},
				{"pkg/services/sqlstore/sqlutil/sqlutil.go", apptmpl.ServiceSqlutil},
				{"pkg/services/sqlstore/annotations.go", apptmpl.ServiceSqlstoreAnnotations},
				{"pkg/services/sqlstore/annotations_test.go", apptmpl.ServiceSqlstoreAnnotationsTest},
				{"pkg/services/sqlstore/health.go", apptmpl.ServiceSqlstoreHealth},
				{"pkg/services/sqlstore/logger.go", apptmpl.ServiceSqlstoreLogger},
				{"pkg/services/sqlstore/login_attempt.go", apptmpl.ServiceSqlstoreLoginAttempt},
				{"pkg/services/sqlstore/login_attempt_test.go", apptmpl.ServiceSqlstoreLoginAttemptTest},
				{"pkg/services/sqlstore/session.go", apptmpl.ServiceSqlstoreSession},
				{"pkg/services/sqlstore/sqlstore.go", apptmpl.ServiceSqlstore},
				{"pkg/services/sqlstore/sqlstore_test.go", apptmpl.ServiceSqlstoreTest},
				{"pkg/services/sqlstore/tags.go", apptmpl.ServiceSqlstoreTags},
				{"pkg/services/sqlstore/team.go", apptmpl.ServiceSqlstoreTeam},
				{"pkg/services/sqlstore/tls_mysql.go", apptmpl.ServiceSqlstoreTlsMySQL},
				{"pkg/services/sqlstore/transactions.go", apptmpl.ServiceSqlstoreTransaction},
				{"pkg/services/sqlstore/transactions_test.go", apptmpl.ServiceSqlstoreTransactionTest},
				{"pkg/services/sqlstore/user.go", apptmpl.ServiceSqlstoreUser},
				{"pkg/services/sqlstore/user_auth.go", apptmpl.ServiceSqlstoreUserAuth},
				{"pkg/services/sqlstore/user_test.go", apptmpl.ServiceSqlstoreUserTest},
				{"pkg/setting/testdata/invalid.ini", apptmpl.SettingInvalid},
				{"pkg/setting/testdata/override.ini", apptmpl.SettingOverride},
				{"pkg/setting/testdata/override_windows.ini", apptmpl.SettingOverrideWindows},
				{"pkg/setting/testdata/session.ini", apptmpl.SettingOverrideWindows},
				{"pkg/setting/setting.go", apptmpl.Setting},
				{"pkg/setting/setting_oauth.go", apptmpl.SettingOauth},
				{"pkg/setting/setting_quota.go", apptmpl.SettingQuota},
				{"pkg/setting/setting_session_test.go", apptmpl.SettingSessionTest},
				{"pkg/setting/setting_smtp.go", apptmpl.SettingSmtp},
				{"pkg/setting/setting_test.go", apptmpl.SettingTest},
				{"pkg/util/errutil/errors.go", apptmpl.UtilErrors},
				{"pkg/util/encoding.go", apptmpl.UtilEncoding},
				{"pkg/util/encoding_test.go", apptmpl.UtilEncodingTest},
				{"pkg/util/encryption.go", apptmpl.UtilEncryption},
				{"pkg/util/filepath.go", apptmpl.UtilFilepath},
				{"pkg/util/ip_address.go", apptmpl.UtilIpAddress},
				{"pkg/util/ip_address_test.go", apptmpl.UtilIpAddressTest},
				{"pkg/util/json.go", apptmpl.UtilJson},
				{"pkg/util/math.go", apptmpl.UtilMath},
				{"pkg/util/md5.go", apptmpl.UtilMd5},
				{"pkg/util/shortid_genarator.go", apptmpl.UtilShortIdGenerator},
				{"pkg/util/strings.go", apptmpl.UtilStrings},
				{"pkg/util/url.go", apptmpl.UtilUrl},
				{"pkg/util/url_test.go", apptmpl.UtilUrlTest},
				{"pkg/util/validation.go", apptmpl.UtilValidation},
				{"pkg/util/validation_test.go", apptmpl.UtilValidationTest},
				{"public/emails/alert_notification.html", apptmpl.PublicAlertNotification},
				{"public/emails/alert_notification_example.html", apptmpl.PublicAlertNotificationExample},
				{"public/emails/invited_to_org.html", apptmpl.PublicInvitedToOrg},
				{"public/emails/new_user_invite.html", apptmpl.PublicNewUserInvite},
				{"public/emails/reset_password.html", apptmpl.PublicResetPassword},
				{"public/emails/signup_started.html", apptmpl.PublicSignupStarted},
				{"public/emails/welcome_on_signup.html", apptmpl.PublicWelcomeOnSignup},

			},
			Comments: []string{},
		}
		if useGoModule != "off" {
			c.Files = append(c.Files, file{"go.mod", apptmpl.GoMod})
			c.Files = append(c.Files, file{"go.sum.mod", apptmpl.GoSum})
		}
	default:
		fmt.Println("Unknown type", atype)
		return
	}

	// set gomodule
	if (useGoModule == "on" || useGoModule == "auto") && atype != "app" {
		c.Files = append(c.Files, file{"go.mod", tmpl.Module})
	}

	if err := create(c); err != nil {
		fmt.Println(err)
		return
	}
}

func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:  "new",
			Usage: "Create a service template",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "namespace",
					Usage: "Namespace for the service e.g com.example",
					Value: "go.micro",
				},
				cli.StringFlag{
					Name:  "type",
					Usage: "Type of service e.g api, fnc, srv, web",
					Value: "srv",
				},
				cli.StringFlag{
					Name:  "fqdn",
					Usage: "FQDN of service e.g com.example.srv.service (defaults to namespace.type.alias)",
				},
				cli.StringFlag{
					Name:  "alias",
					Usage: "Alias is the short name used as part of combined name if specified",
				},
				cli.StringSliceFlag{
					Name:  "plugin",
					Usage: "Specify plugins e.g --plugin=registry=etcd:broker=nats or use flag multiple times",
				},
				cli.BoolTFlag{
					Name:  "gopath",
					Usage: "Create the service in the gopath. Defaults to true.",
				},
			},
			Action: func(c *cli.Context) {
				run(c)
			},
		},
	}
}
