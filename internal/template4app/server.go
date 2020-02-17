package template4app

var (
	ServerAdminUser = `
package server

import (
	"{{.Dir}}/pkg/bus"
	"{{.Dir}}/pkg/infra/metrics"
	"{{.Dir}}/pkg/models"
	"{{.Dir}}/pkg/server/dtos"
	"{{.Dir}}/pkg/util"
)

func AdminCreateUser(c *models.ReqContext, form dtos.AdminCreateUserForm) {
	cmd := models.CreateUserCommand{
		Login:    form.Login,
		Email:    form.Email,
		Password: form.Password,
		Name:     form.Name,
	}

	if len(cmd.Login) == 0 {
		cmd.Login = cmd.Email
		if len(cmd.Login) == 0 {
			c.JsonApiErr(400, "Validation error, need specify either username or email", nil)
			return
		}
	}

	if len(cmd.Password) < 4 {
		c.JsonApiErr(400, "Password is missing or too short", nil)
		return
	}

	if err := bus.Dispatch(&cmd); err != nil {
		c.JsonApiErr(500, "failed to create user", err)
		return
	}

	metrics.M_Api_Admin_User_Create.Inc()

	user := cmd.Result

	result := models.UserIdDTO{
		Message: "User created",
		Id:      user.Id,
	}

	c.JSON(200, result)
}

func AdminUpdateUserPassword(c *models.ReqContext, form dtos.AdminUpdateUserPasswordForm) {
	userID := c.ParamsInt64(":id")

	if len(form.Password) < 4 {
		c.JsonApiErr(400, "New password too short", nil)
		return
	}

	userQuery := models.GetUserByIdQuery{Id: userID}

	if err := bus.Dispatch(&userQuery); err != nil {
		c.JsonApiErr(500, "Could not read user from database", err)
		return
	}

	passwordHashed := util.EncodePassword(form.Password, userQuery.Result.Salt)

	cmd := models.ChangeUserPasswordCommand{
		UserId:      userID,
		NewPassword: passwordHashed,
	}

	if err := bus.Dispatch(&cmd); err != nil {
		c.JsonApiErr(500, "Failed to update user password", err)
		return
	}

	c.JsonOK("User password updated")
}

// PUT /api/admin/users/:id/permissions
func AdminUpdateUserPermissions(c *models.ReqContext, form dtos.AdminUpdateUserPermissionsForm) {
	userID := c.ParamsInt64(":id")

	cmd := models.UpdateUserPermissionsCommand{
		UserId:         userID,
		IsGrafanaAdmin: form.IsGrafanaAdmin,
	}

	if err := bus.Dispatch(&cmd); err != nil {
		if err == models.ErrLastGrafanaAdmin {
			c.JsonApiErr(400, models.ErrLastGrafanaAdmin.Error(), nil)
			return
		}

		c.JsonApiErr(500, "Failed to update user permissions", err)
		return
	}

	c.JsonOK("User permissions updated")
}

func AdminDeleteUser(c *models.ReqContext) {
	userID := c.ParamsInt64(":id")

	cmd := models.DeleteUserCommand{UserId: userID}

	if err := bus.Dispatch(&cmd); err != nil {
		c.JsonApiErr(500, "Failed to delete user", err)
		return
	}

	c.JsonOK("User deleted")
}

// POST /api/admin/users/:id/disable
func (server *HTTPServer) AdminDisableUser(c *models.ReqContext) Response {
	userID := c.ParamsInt64(":id")

	// External users shouldn't be disabled from API
	authInfoQuery := &models.GetAuthInfoQuery{UserId: userID}
	if err := bus.Dispatch(authInfoQuery); err != models.ErrUserNotFound {
		return Error(500, "Could not disable external user", nil)
	}

	disableCmd := models.DisableUserCommand{UserId: userID, IsDisabled: true}
	if err := bus.Dispatch(&disableCmd); err != nil {
		return Error(500, "Failed to disable user", err)
	}

	err := server.AuthTokenService.RevokeAllUserTokens(c.Req.Context(), userID)
	if err != nil {
		return Error(500, "Failed to disable user", err)
	}

	return Success("User disabled")
}

// POST /api/admin/users/:id/enable
func AdminEnableUser(c *models.ReqContext) Response {
	userID := c.ParamsInt64(":id")

	// External users shouldn't be disabled from API
	authInfoQuery := &models.GetAuthInfoQuery{UserId: userID}
	if err := bus.Dispatch(authInfoQuery); err != models.ErrUserNotFound {
		return Error(500, "Could not enable external user", nil)
	}

	disableCmd := models.DisableUserCommand{UserId: userID, IsDisabled: false}
	if err := bus.Dispatch(&disableCmd); err != nil {
		return Error(500, "Failed to enable user", err)
	}

	return Success("User enabled")
}

// POST /api/admin/users/:id/logout
//func (server *HTTPServer) AdminLogoutUser(c *models.ReqContext) Response {
//	userID := c.ParamsInt64(":id")
//
//	if c.UserId == userID {
//		return Error(400, "You cannot logout yourself", nil)
//	}
//
//	return server.logoutUserFromAllDevicesInternal(c.Req.Context(), userID)
//}
//
//// GET /api/admin/users/:id/auth-tokens
//func (server *HTTPServer) AdminGetUserAuthTokens(c *models.ReqContext) Response {
//	userID := c.ParamsInt64(":id")
//	return server.getUserAuthTokensInternal(c, userID)
//}
//
//// POST /api/admin/users/:id/revoke-auth-token
//func (server *HTTPServer) AdminRevokeUserAuthToken(c *models.ReqContext, cmd models.RevokeAuthTokenCmd) Response {
//	userID := c.ParamsInt64(":id")
//	return server.revokeUserAuthTokenInternal(c, userID, cmd)
//}

`
	ServerCommon = `
package server

import (
	"encoding/json"
	m "{{.Dir}}/pkg/models"
	"{{.Dir}}/pkg/setting"
	"gopkg.in/macaron.v1"
	"net/http"
)

var (
	NotFound = func() Response {
		return Error(404, "Not found", nil)
	}
	ServerError = func(err error) Response {
		return Error(500, "Server error", err)
	}
)

type Response interface {
	WriteTo(ctx *m.ReqContext)
}

type NormalResponse struct {
	status     int
	body       []byte
	header     http.Header
	errMessage string
	err        error
}

func Wrap(action interface{}) macaron.Handler {

	return func(c *m.ReqContext) {
		var res Response
		val, err := c.Invoke(action)
		if err == nil && val != nil && len(val) > 0 {
			res = val[0].Interface().(Response)
		} else {
			res = ServerError(err)
		}

		res.WriteTo(c)
	}
}

// Error create a erroneous response
func Error(status int, message string, err error) *NormalResponse {
	data := make(map[string]interface{})

	switch status {
	case 404:
		data["message"] = "Not Found"
	case 500:
		data["message"] = "Internal Server Error"
	}

	if message != "" {
		data["message"] = message
	}

	if err != nil {
		if setting.Env != setting.PROD {
			data["error"] = err.Error()
		}
	}

	resp := JSON(status, data)

	if err != nil {
		resp.errMessage = message
		resp.err = err
	}

	return resp
}

func (r *NormalResponse) WriteTo(ctx *m.ReqContext) {
	if r.err != nil {
		ctx.Logger.Error(r.errMessage, "error", r.err)
	}

	header := ctx.Resp.Header()
	for k, v := range r.header {
		header[k] = v
	}
	ctx.Resp.WriteHeader(r.status)
	ctx.Resp.Write(r.body)
}

func (r *NormalResponse) Cache(ttl string) *NormalResponse {
	return r.Header("Cache-Control", "public,max-age="+ttl)
}

func (r *NormalResponse) Header(key, value string) *NormalResponse {
	r.header.Set(key, value)
	return r
}

// Empty create an empty response
func Empty(status int) *NormalResponse {
	return Respond(status, nil)
}

// JSON create a JSON response
func JSON(status int, body interface{}) *NormalResponse {
	return Respond(status, body).Header("Content-Type", "application/json")
}

// Success create a successful response
func Success(message string) *NormalResponse {
	resp := make(map[string]interface{})
	resp["message"] = message
	return JSON(200, resp)
}

// Respond create a response
func Respond(status int, body interface{}) *NormalResponse {
	var b []byte
	var err error
	switch t := body.(type) {
	case []byte:
		b = t
	case string:
		b = []byte(t)
	default:
		if b, err = json.Marshal(body); err != nil {
			return Error(500, "body json marshal", err)
		}
	}
	return &NormalResponse{
		body:   b,
		status: status,
		header: make(http.Header),
	}
}

`
	ServerHttpServer = `
package server

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"{{.Dir}}/pkg/bus"
	"{{.Dir}}/pkg/components/simplejson"
	"{{.Dir}}/pkg/infra/localcache"
	"{{.Dir}}/pkg/infra/log"
	"{{.Dir}}/pkg/middleware"
	"{{.Dir}}/pkg/models"
	"{{.Dir}}/pkg/registry"
	"{{.Dir}}/pkg/routing"
	httpstatic "{{.Dir}}/pkg/server/static"
	"{{.Dir}}/pkg/services/quota"
	"{{.Dir}}/pkg/services/remotecache"
	"{{.Dir}}/pkg/setting"
	"gopkg.in/macaron.v1"
	"net"
	"net/http"
	"os"
	"path"
	"time"
)

func init() {
	registry.Register(&registry.Descriptor{
		Name:         "HTTPServer",
		Instance:     &HTTPServer{},
		InitPriority: registry.High,
	})
}

type HTTPServer struct {
	log     log.Logger
	macaron *macaron.Macaron
	context context.Context
	//streamManager *live.StreamManager
	httpSrv *http.Server

	RouteRegister routing.RouteRegister `+"`inject:\"\"`"+`
	Bus           bus.Bus               ` +"`inject:\"\"`"+`

	Cfg *setting.Cfg ` +"`inject:\"\"`"+`

	CacheService *localcache.CacheService ` +"`inject:\"\"`"+`

	AuthTokenService   models.UserTokenService  ` +"`inject:\"\"`"+`
	QuotaService       *quota.QuotaService      ` +"`inject:\"\"`"+`
	RemoteCacheService *remotecache.RemoteCache ` +"`inject:\"\"`"+`
	//ProvisioningService ProvisioningService      ` +"`inject:\"\"`"+`
	//Login               *login.LoginService      ` +"`inject:\"\"`"+`
}

func (hs *HTTPServer) Init() error {
	hs.log = log.New("http.server")

	hs.macaron = hs.newMacaron()
	hs.registerRoutes()

	return nil
}

func (hs *HTTPServer) newMacaron() *macaron.Macaron {
	macaron.Env = setting.Env
	m := macaron.New()

	// automatically set HEAD for every GET
	m.SetAutoHead(true)

	return m
}

func (hs *HTTPServer) applyRoutes() {
	// start with middlewares & static routes
	hs.addMiddlewaresAndStaticRoutes()
	// then add view routes & api routes
	hs.RouteRegister.Register(hs.macaron)
	// then custom app proxy routes
	//hs.initAppPluginRoutes(hs.macaron)
	// lastly not found route
	hs.macaron.NotFound(hs.NotFoundHandler)
}

func (hs *HTTPServer) addMiddlewaresAndStaticRoutes() {
	m := hs.macaron

	m.Use(middleware.Logger())

	if setting.EnableGzip {
		m.Use(middleware.Gziper())
	}

	m.Use(middleware.Recovery())

	//plugin
	//for _, route := range plugins.StaticRoutes {
	//	pluginRoute := path.Join("/public/plugins/", route.PluginId)
	//	hs.log.Debug("Plugins: Adding route", "route", pluginRoute, "dir", route.Directory)
	//	hs.mapStatic(hs.macaron, route.Directory, "", pluginRoute)
	//}

	hs.mapStatic(m, setting.StaticRootPath, "build", "public/build")
	hs.mapStatic(m, setting.StaticRootPath, "", "public")
	hs.mapStatic(m, setting.StaticRootPath, "robots.txt", "robots.txt")

	if setting.ImageUploadProvider == "local" {
		hs.mapStatic(m, hs.Cfg.ImagesDir, "", "/public/img/attachments")
	}

	m.Use(middleware.AddDefaultResponseHeaders())

	if setting.ServeFromSubPath && setting.AppSubUrl != "" {
		m.SetURLPrefix(setting.AppSubUrl)
	}

	m.Use(macaron.Renderer(macaron.RenderOptions{
		Directory:  path.Join(setting.StaticRootPath, "views"),
		IndentJSON: macaron.Env != macaron.PROD,
		Delims:     macaron.Delims{Left: "[[", Right: "]]"},
	}))

	//health for DB connection
	m.Use(hs.healthHandler)

	m.Use(middleware.GetContextHandler(
		hs.AuthTokenService,
		hs.RemoteCacheService,
	))
	//m.Use(middleware.OrgRedirect())

	m.Use(middleware.HandleNoCacheHeader())
}

func (hs *HTTPServer) healthHandler(ctx *macaron.Context) {
	notHeadOrGet := ctx.Req.Method != http.MethodGet && ctx.Req.Method != http.MethodHead
	if notHeadOrGet || ctx.Req.URL.Path != "/api/health" {
		return
	}

	data := simplejson.New()
	data.Set("database", "ok")
	data.Set("version", setting.BuildVersion)
	data.Set("commit", setting.BuildCommit)

	if err := bus.Dispatch(&models.GetDBHealthQuery{}); err != nil {
		data.Set("database", "failing")
		ctx.Resp.Header().Set("Content-Type", "application/json; charset=UTF-8")
		ctx.Resp.WriteHeader(503)
	} else {
		ctx.Resp.Header().Set("Content-Type", "application/json; charset=UTF-8")
		ctx.Resp.WriteHeader(200)
	}

	dataBytes, _ := data.EncodePretty()
	ctx.Resp.Write(dataBytes)
}

func (hs *HTTPServer) mapStatic(m *macaron.Macaron, rootDir string, dir string, prefix string) {
	headers := func(c *macaron.Context) {
		c.Resp.Header().Set("Cache-Control", "public, max-age=3600")
	}

	if prefix == "public/build" {
		headers = func(c *macaron.Context) {
			c.Resp.Header().Set("Cache-Control", "public, max-age=31536000")
		}
	}

	if setting.Env == setting.DEV {
		headers = func(c *macaron.Context) {
			c.Resp.Header().Set("Cache-Control", "max-age=0, must-revalidate, no-cache")
		}
	}

	m.Use(httpstatic.Static(
		path.Join(rootDir, dir),
		httpstatic.StaticOptions{
			SkipLogging: true,
			Prefix:      prefix,
			AddHeaders:  headers,
		},
	))
}

func (hs *HTTPServer) metricsEndpointBasicAuthEnabled() bool {
	return hs.Cfg.MetricsEndpointBasicAuthUsername != "" && hs.Cfg.MetricsEndpointBasicAuthPassword != ""
}

func (hs *HTTPServer) listenAndServeTLS(certfile, keyfile string) error {
	if certfile == "" {
		return fmt.Errorf("cert_file cannot be empty when using HTTPS")
	}

	if keyfile == "" {
		return fmt.Errorf("cert_key cannot be empty when using HTTPS")
	}

	if _, err := os.Stat(setting.CertFile); os.IsNotExist(err) {
		return fmt.Errorf(` +"`Cannot find SSL cert_file at %v`"+`, setting.CertFile)
	}

	if _, err := os.Stat(setting.KeyFile); os.IsNotExist(err) {
		return fmt.Errorf(`+"`Cannot find SSL key_file at %v`"+`, setting.KeyFile)
	}

	tlsCfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		},
	}

	hs.httpSrv.TLSConfig = tlsCfg
	hs.httpSrv.TLSNextProto = make(map[string]func(*http.Server, *tls.Conn, http.Handler))

	return hs.httpSrv.ListenAndServeTLS(setting.CertFile, setting.KeyFile)
}

func (hs *HTTPServer) Run(ctx context.Context) error {
	var err error

	hs.context = ctx

	//url mapping handler
	hs.applyRoutes()

	listenAddr := fmt.Sprintf("%s:%s", setting.HttpAddr, setting.HttpPort)
	hs.log.Info("HTTP Server Listen", "address", listenAddr, "protocol", setting.Protocol, "subUrl", setting.AppSubUrl, "socket", setting.SocketPath)

	hs.httpSrv = &http.Server{Addr: listenAddr, Handler: hs.macaron}

	// handle http shutdown on server context done
	go func() {
		<-ctx.Done()
		// Hacky fix for race condition between ListenAndServe and Shutdown
		time.Sleep(time.Millisecond * 100)
		if err := hs.httpSrv.Shutdown(context.Background()); err != nil {
			hs.log.Error("Failed to shutdown server", "error", err)
		}
	}()

	switch setting.Protocol {
	case setting.HTTP:
		err = hs.httpSrv.ListenAndServe()
		if err == http.ErrServerClosed {
			hs.log.Debug("server was shutdown gracefully")
			return nil
		}
	case setting.HTTPS:
		err = hs.listenAndServeTLS(setting.CertFile, setting.KeyFile)
		if err == http.ErrServerClosed {
			hs.log.Debug("server was shutdown gracefully")
			return nil
		}
	case setting.SOCKET:
		ln, err := net.ListenUnix("unix", &net.UnixAddr{Name: setting.SocketPath, Net: "unix"})
		if err != nil {
			hs.log.Debug("server was shutdown gracefully")
			return nil
		}

		// Make socket writable by group
		os.Chmod(setting.SocketPath, 0660)

		err = hs.httpSrv.Serve(ln)
		if err != nil {
			hs.log.Debug("server was shutdown gracefully")
			return nil
		}
	default:
		hs.log.Error("Invalid protocol", "protocol", setting.Protocol)
		err = errors.New("Invalid Protocol")
	}

	hs.log.Info("server starts with ", setting.Protocol, " protocol")
	return err
}

`
	ServerIndex = `
package server

import (
	"fmt"
	"github.com/elitecodegroovy/util"
	"gopkg.in/macaron.v1"
	"io"
	"os"
	"regexp"
	"strings"

	//"{{.Dir}}/pkg/bus"
	m "{{.Dir}}/pkg/models"
	"{{.Dir}}/pkg/setting"
)

func AdminGetSettings(c *m.ReqContext) {
	settings := make(map[string]interface{})

	for _, section := range setting.Raw.Sections() {
		jsonSec := make(map[string]interface{})
		settings[section.Name()] = jsonSec

		for _, key := range section.Keys() {
			keyName := key.Name()
			value := key.Value()
			if strings.Contains(keyName, "secret") ||
				strings.Contains(keyName, "password") ||
				(strings.Contains(keyName, "provider_config")) {
				value = "************"
			}
			if strings.Contains(keyName, "url") {
				var rgx = regexp.MustCompile(` + "`.*:\\/\\/([^:]*):([^@]*)@.*?$`"+`)
				var subs = rgx.FindAllSubmatch([]byte(value), -1)
				if subs != nil && len(subs[0]) == 3 {
					value = strings.Replace(value, string(subs[0][1]), "******", 1)
					value = strings.Replace(value, string(subs[0][2]), "******", 1)
				}
			}

			jsonSec[keyName] = value
		}
	}

	c.JSON(200, settings)
}

func urlHandler(c *macaron.Context) {
	settings := make(map[string]interface{})
	settings["msg"] = "the request path is: " + c.Req.RequestURI
	settings["code"] = 200
	c.JSON(200, settings)
}

func uploadFile(c *m.ReqContext) {
	r := c.Req
	w := c.Resp
	// the FormFile function takes in the POST input id file
	file, header, err := r.FormFile("file")

	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	defer file.Close()

	c.Logger.Info("Uploaded File:", header.Filename)

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	out, err := os.Create("./" + util.GetCurrentTimeNumberISOStrTime() + "-" + header.Filename)
	if err != nil {
		c.JSON(300, err.Error())
		return
	}

	defer out.Close()

	// write the content from POST to the file
	_, err = io.Copy(out, file)
	if err != nil {
		c.JSON(301, err.Error())
	}

	c.JSON(200, "ok")
}

func (hs *HTTPServer) NotFoundHandler(c *m.ReqContext) {
	if c.IsApiRequest() {
		c.JsonApiErr(404, "Not found", nil)
		return
	}

	//data, err := hs.setIndexViewData(c)
	//if err != nil {
	//	c.Handle(500, "Failed to get settings", err)
	//	return
	//}

	c.HTML(404, "index", nil)
}

`
	ServerLogin = `
package server

import (
	"errors"
	"{{.Dir}}/pkg/bus"
	"{{.Dir}}/pkg/infra/metrics"
	"{{.Dir}}/pkg/middleware"
	m "{{.Dir}}/pkg/models"
	"{{.Dir}}/pkg/server/dtos"
	"{{.Dir}}/pkg/setting"
	"net/url"
)

var (
	ErrEmailNotAllowed       = errors.New("Required email domain not fulfilled")
	ErrInvalidCredentials    = errors.New("Invalid Username or Password")
	ErrNoEmail               = errors.New("Login provider didn't return an email address")
	ErrProviderDeniedRequest = errors.New("Login provider denied login request")
	ErrSignUpNotAllowed      = errors.New("Signup is not allowed for this adapter")
	ErrTooManyLoginAttempts  = errors.New("Too many consecutive incorrect login attempts for user. Login for user temporarily blocked")
	ErrPasswordEmpty         = errors.New("No password provided")
	ErrUserDisabled          = errors.New("User is disabled")
)

func (hs *HTTPServer) LoginPost(c *m.ReqContext, cmd dtos.LoginCommand) Response {
	if setting.DisableLoginForm {
		return Error(401, "Login is disabled", nil)
	}

	authQuery := &m.LoginUserQuery{
		ReqContext: c,
		Username:   cmd.User,
		Password:   cmd.Password,
		IpAddress:  c.Req.RemoteAddr,
	}

	if err := bus.Dispatch(authQuery); err != nil {
		if err == ErrInvalidCredentials || err == ErrTooManyLoginAttempts {
			return Error(401, "Invalid username or password", err)
		}

		if err == ErrUserDisabled {
			return Error(401, "User is disabled", err)
		}

		return Error(500, "Error while trying to authenticate user", err)
	}

	user := authQuery.User

	hs.loginUserWithUser(user, c)

	result := map[string]interface{}{
		"message": "Logged in",
	}

	if redirectTo, _ := url.QueryUnescape(c.GetCookie("redirect_to")); len(redirectTo) > 0 {
		result["redirectUrl"] = redirectTo
		c.SetCookie("redirect_to", "", -1, setting.AppSubUrl+"/")
	}

	metrics.M_Api_Login_Post.Inc()

	return JSON(200, result)
}

func (hs *HTTPServer) loginUserWithUser(user *m.User, c *m.ReqContext) {
	if user == nil {
		hs.log.Error("user login with nil user")
	}

	userToken, err := hs.AuthTokenService.CreateToken(c.Req.Context(), user.Id, c.RemoteAddr(), c.Req.UserAgent())
	if err != nil {
		hs.log.Error("failed to create auth token", "error", err)
	}

	middleware.WriteSessionCookie(c, userToken.UnhashedToken, hs.Cfg.LoginMaxLifetimeDays)
}

`
	ServerUrlMapping = `
package server

import (
	"{{.Dir}}/pkg/middleware"
	"{{.Dir}}/pkg/routing"
	"{{.Dir}}/pkg/server/dtos"
	"github.com/go-macaron/binding"
)

func (hs *HTTPServer) registerRoutes() {
	//reqSignedIn := middleware.ReqSignedIn
	reqGrafanaAdmin := middleware.ReqGrafanaAdmin
	//reqEditorRole := middleware.ReqEditorRole
	//reqOrgAdmin := middleware.ReqOrgAdmin
	//reqCanAccessTeams := middleware.AdminOrFeatureEnabled(hs.Cfg.EditorsCanAdmin)
	//redirectFromLegacyDashboardURL := middleware.RedirectFromLegacyDashboardURL()
	//redirectFromLegacyDashboardSoloURL := middleware.RedirectFromLegacyDashboardSoloURL()
	quota := middleware.Quota(hs.QuotaService)
	bind := binding.Bind

	r := hs.RouteRegister

	// not logged in views
	r.Get("/", func() string {
		return "Macaron Web Framework!"
	})
	// admin api
	r.Group("/api/admin", func(adminRoute routing.RouteRegister) {
		adminRoute.Get("/settings", AdminGetSettings)
		adminRoute.Post("/users", bind(dtos.AdminCreateUserForm{}), AdminCreateUser)
	}, reqGrafanaAdmin)

	r.Post("/login", quota("session"), bind(dtos.LoginCommand{}), Wrap(hs.LoginPost))

	r.Get("/urlReq", urlHandler)

	r.Post("/upload", uploadFile)

	r.Get("/setting", AdminGetSettings)
	//r.Get("/api", reqSignedIn)

}

`
)
