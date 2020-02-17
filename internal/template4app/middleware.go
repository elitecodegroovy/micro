package template4app

var(
	Auth = `
package middleware

import (
	"net/url"
	"strings"

	macaron "gopkg.in/macaron.v1"

	m "{{.Dir}}/pkg/models"
	"{{.Dir}}/pkg/setting"
	"{{.Dir}}/pkg/util"
)

type AuthOptions struct {
	ReqGrafanaAdmin bool
	ReqSignedIn     bool
}

func getApiKey(c *m.ReqContext) string {
	header := c.Req.Header.Get("Authorization")
	parts := strings.SplitN(header, " ", 2)
	if len(parts) == 2 && parts[0] == "Bearer" {
		key := parts[1]
		return key
	}

	username, password, err := util.DecodeBasicAuthHeader(header)
	if err == nil && username == "api_key" {
		return password
	}

	return ""
}

func accessForbidden(c *m.ReqContext) {
	if c.IsApiRequest() {
		c.JsonApiErr(403, "Permission denied", nil)
		return
	}

	c.Redirect(setting.AppSubUrl + "/")
}

func notAuthorized(c *m.ReqContext) {
	if c.IsApiRequest() {
		c.JsonApiErr(401, "Unauthorized", nil)
		return
	}

	c.SetCookie("redirect_to", url.QueryEscape(setting.AppSubUrl+c.Req.RequestURI), 0, setting.AppSubUrl+"/", nil, false, true)

	c.Redirect(setting.AppSubUrl + "/login")
}

func EnsureEditorOrViewerCanEdit(c *m.ReqContext) {
	if !c.SignedInUser.HasRole(m.ROLE_EDITOR) && !setting.ViewersCanEdit {
		accessForbidden(c)
	}
}

func RoleAuth(roles ...m.RoleType) macaron.Handler {
	return func(c *m.ReqContext) {
		ok := false
		for _, role := range roles {
			if role == c.OrgRole {
				ok = true
				break
			}
		}
		if !ok {
			accessForbidden(c)
		}
	}
}

func Auth(options *AuthOptions) macaron.Handler {
	return func(c *m.ReqContext) {
		if !c.IsSignedIn && options.ReqSignedIn && !c.AllowAnonymous {
			notAuthorized(c)
			return
		}

		if !c.IsGrafanaAdmin && options.ReqGrafanaAdmin {
			accessForbidden(c)
			return
		}
	}
}

// AdminOrFeatureEnabled creates a middleware that allows access
// if the signed in user is either an Org Admin or if the
// feature flag is enabled.
// Intended for when feature flags open up access to APIs that
// are otherwise only available to admins.
func AdminOrFeatureEnabled(enabled bool) macaron.Handler {
	return func(c *m.ReqContext) {
		if c.OrgRole == m.ROLE_ADMIN {
			return
		}

		if !enabled {
			accessForbidden(c)
		}
	}
}

`
	Gzipper = `
package middleware

import (
	"strings"

	"github.com/go-macaron/gzip"
	"gopkg.in/macaron.v1"
)

func Gziper() macaron.Handler {
	macaronGziper := gzip.Gziper()

	return func(ctx *macaron.Context) {
		requestPath := ctx.Req.URL.RequestURI()
		// ignore datasource proxy requests
		if strings.HasPrefix(requestPath, "/api/datasources/proxy") {
			return
		}

		if strings.HasPrefix(requestPath, "/api/plugin-proxy/") {
			return
		}

		if strings.HasPrefix(requestPath, "/metrics") {
			return
		}

		ctx.Invoke(macaronGziper)
	}
}

`
	Header = `
package middleware

import (
	m "{{.Dir}}/pkg/models"
	macaron "gopkg.in/macaron.v1"
)

const HeaderNameNoBackendCache = "X-Grafana-NoCache"

func HandleNoCacheHeader() macaron.Handler {
	return func(ctx *m.ReqContext) {
		ctx.SkipCache = ctx.Req.Header.Get(HeaderNameNoBackendCache) == "true"
	}
}

`
	Logger = `
package middleware

import (
	"net/http"
	"time"

	m "{{.Dir}}/pkg/models"
	"{{.Dir}}/pkg/setting"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/macaron.v1"
)

func Logger() macaron.Handler {
	return func(res http.ResponseWriter, req *http.Request, c *macaron.Context) {
		start := time.Now()
		c.Data["perfmon.start"] = start

		rw := res.(macaron.ResponseWriter)
		c.Next()

		timeTakenMs := time.Since(start) / time.Millisecond

		if timer, ok := c.Data["perfmon.timer"]; ok {
			timerTyped := timer.(prometheus.Summary)
			timerTyped.Observe(float64(timeTakenMs))
		}

		status := rw.Status()
		if status == 200 || status == 304 {
			if !setting.RouterLogging {
				return
			}
		}

		if ctx, ok := c.Data["ctx"]; ok {
			ctxTyped := ctx.(*m.ReqContext)
			if status == 500 {
				ctxTyped.Logger.Error("Request Completed", "method", req.Method, "path", req.URL.Path, "status", status, "remote_addr", c.RemoteAddr(), "time_ms", int64(timeTakenMs), "size", rw.Size(), "referer", req.Referer())
			} else {
				ctxTyped.Logger.Info("Request Completed", "method", req.Method, "path", req.URL.Path, "status", status, "remote_addr", c.RemoteAddr(), "time_ms", int64(timeTakenMs), "size", rw.Size(), "referer", req.Referer())
			}
		}
	}
}
`
	Middleware = `
package middleware

import (
	"fmt"
	"{{.Dir}}/pkg/services/remotecache"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	macaron "gopkg.in/macaron.v1"

	"{{.Dir}}/pkg/bus"
	"{{.Dir}}/pkg/components/apikeygen"
	"{{.Dir}}/pkg/infra/log"
	"{{.Dir}}/pkg/models"
	"{{.Dir}}/pkg/setting"
	"{{.Dir}}/pkg/util"
)

var getTime = time.Now

var (
	ReqGrafanaAdmin = Auth(&AuthOptions{ReqSignedIn: true, ReqGrafanaAdmin: true})
	ReqSignedIn     = Auth(&AuthOptions{ReqSignedIn: true})
	ReqEditorRole   = RoleAuth(models.ROLE_EDITOR, models.ROLE_ADMIN)
	ReqOrgAdmin     = RoleAuth(models.ROLE_ADMIN)
)

func GetContextHandler(ats models.UserTokenService, remoteCache *remotecache.RemoteCache) macaron.Handler {
	return func(c *macaron.Context) {
		ctx := &models.ReqContext{
			Context:        c,
			SignedInUser:   &models.SignedInUser{},
			IsSignedIn:     false,
			AllowAnonymous: false,
			SkipCache:      false,
			Logger:         log.New("context"),
		}

		orgId := int64(0)
		orgIdHeader := ctx.Req.Header.Get("X-Grafana-Org-Id")
		if orgIdHeader != "" {
			orgId, _ = strconv.ParseInt(orgIdHeader, 10, 64)
		}

		// the order in which these are tested are important
		// look for api key in Authorization header first
		// then init session and look for userId in session
		// then look for api key in session (special case for render calls via api)
		// then test if anonymous access is enabled
		switch {
		//case initContextWithRenderAuth(ctx):
		case initContextWithApiKey(ctx):
		case initContextWithBasicAuth(ctx, orgId):
		//case initContextWithAuthProxy(remoteCache, ctx, orgId):
		case initContextWithToken(ats, ctx, orgId):
		case initContextWithAnonymousUser(ctx):
		}

		ctx.Logger = log.New("context", "userId", ctx.UserId, "orgId", ctx.OrgId, "uname", ctx.Login)
		ctx.Data["ctx"] = ctx

		c.Map(ctx)

		// update last seen every 5min
		if ctx.ShouldUpdateLastSeenAt() {
			ctx.Logger.Debug("Updating last user_seen_at", "user_id", ctx.UserId)
			if err := bus.Dispatch(&models.UpdateUserLastSeenAtCommand{UserId: ctx.UserId}); err != nil {
				ctx.Logger.Error("Failed to update last_seen_at", "error", err)
			}
		}
	}
}

func initContextWithAnonymousUser(ctx *models.ReqContext) bool {
	if !setting.AnonymousEnabled {
		return false
	}

	orgQuery := models.GetOrgByNameQuery{Name: setting.AnonymousOrgName}
	if err := bus.Dispatch(&orgQuery); err != nil {
		log.Error(3, "Anonymous access organization error: '%s': %s", setting.AnonymousOrgName, err)
		return false
	}

	ctx.IsSignedIn = false
	ctx.AllowAnonymous = true
	ctx.SignedInUser = &models.SignedInUser{IsAnonymous: true}
	ctx.OrgRole = models.RoleType(setting.AnonymousOrgRole)
	ctx.OrgId = orgQuery.Result.Id
	ctx.OrgName = orgQuery.Result.Name
	return true
}

func initContextWithApiKey(ctx *models.ReqContext) bool {
	var keyString string
	if keyString = getApiKey(ctx); keyString == "" {
		return false
	}

	// base64 decode key
	decoded, err := apikeygen.Decode(keyString)
	if err != nil {
		ctx.JsonApiErr(401, "Invalid API key", err)
		return true
	}

	// fetch key
	keyQuery := models.GetApiKeyByNameQuery{KeyName: decoded.Name, OrgId: decoded.OrgId}
	if err := bus.Dispatch(&keyQuery); err != nil {
		ctx.JsonApiErr(401, "Invalid API key", err)
		return true
	}

	apikey := keyQuery.Result

	// validate api key
	if !apikeygen.IsValid(decoded, apikey.Key) {
		ctx.JsonApiErr(401, "Invalid API key", err)
		return true
	}

	// check for expiration
	if apikey.Expires != nil && *apikey.Expires <= getTime().Unix() {
		ctx.JsonApiErr(401, "Expired API key", err)
		return true
	}

	ctx.IsSignedIn = true
	ctx.SignedInUser = &models.SignedInUser{}
	ctx.OrgRole = apikey.Role
	ctx.ApiKeyId = apikey.Id
	ctx.OrgId = apikey.OrgId
	return true
}

func initContextWithBasicAuth(ctx *models.ReqContext, orgId int64) bool {

	if !setting.BasicAuthEnabled {
		return false
	}

	header := ctx.Req.Header.Get("Authorization")
	if header == "" {
		return false
	}

	username, password, err := util.DecodeBasicAuthHeader(header)
	if err != nil {
		ctx.JsonApiErr(401, "Invalid Basic Auth Header", err)
		return true
	}

	loginQuery := models.GetUserByLoginQuery{LoginOrEmail: username}
	if err := bus.Dispatch(&loginQuery); err != nil {
		ctx.JsonApiErr(401, "Basic auth failed", err)
		return true
	}

	user := loginQuery.Result

	loginUserQuery := models.LoginUserQuery{Username: username, Password: password, User: user}
	if err := bus.Dispatch(&loginUserQuery); err != nil {
		ctx.JsonApiErr(401, "Invalid username or password", err)
		return true
	}

	query := models.GetSignedInUserQuery{UserId: user.Id, OrgId: orgId}
	if err := bus.Dispatch(&query); err != nil {
		ctx.JsonApiErr(401, "Authentication error", err)
		return true
	}

	ctx.SignedInUser = query.Result
	ctx.IsSignedIn = true
	return true
}

func initContextWithToken(authTokenService models.UserTokenService, ctx *models.ReqContext, orgID int64) bool {
	rawToken := ctx.GetCookie(setting.LoginCookieName)
	if rawToken == "" {
		return false
	}

	token, err := authTokenService.LookupToken(ctx.Req.Context(), rawToken)
	if err != nil {
		ctx.Logger.Error("failed to look up user based on cookie", "error", err)
		WriteSessionCookie(ctx, "", -1)
		return false
	}

	query := models.GetSignedInUserQuery{UserId: token.UserId, OrgId: orgID}
	if err := bus.Dispatch(&query); err != nil {
		ctx.Logger.Error("failed to get user with id", "userId", token.UserId, "error", err)
		return false
	}

	ctx.SignedInUser = query.Result
	ctx.IsSignedIn = true
	ctx.UserToken = token

	rotated, err := authTokenService.TryRotateToken(ctx.Req.Context(), token, ctx.RemoteAddr(), ctx.Req.UserAgent())
	if err != nil {
		ctx.Logger.Error("failed to rotate token", "error", err)
		return true
	}

	if rotated {
		WriteSessionCookie(ctx, token.UnhashedToken, setting.LoginMaxLifetimeDays)
	}

	return true
}

func WriteSessionCookie(ctx *models.ReqContext, value string, maxLifetimeDays int) {
	if setting.Env == setting.DEV {
		ctx.Logger.Info("new token", "unhashed token", value)
	}

	var maxAge int
	if maxLifetimeDays <= 0 {
		maxAge = -1
	} else {
		maxAgeHours := (time.Duration(setting.LoginMaxLifetimeDays) * 24 * time.Hour) + time.Hour
		maxAge = int(maxAgeHours.Seconds())
	}

	ctx.Resp.Header().Del("Set-Cookie")
	cookie := http.Cookie{
		Name:     setting.LoginCookieName,
		Value:    url.QueryEscape(value),
		HttpOnly: true,
		Path:     setting.AppSubUrl + "/",
		Secure:   setting.CookieSecure,
		MaxAge:   maxAge,
		SameSite: setting.CookieSameSite,
	}

	http.SetCookie(ctx.Resp, &cookie)
}

func AddDefaultResponseHeaders() macaron.Handler {
	return func(ctx *macaron.Context) {
		ctx.Resp.Before(func(w macaron.ResponseWriter) {
			if !strings.HasPrefix(ctx.Req.URL.Path, "/api/datasources/proxy/") {
				AddNoCacheHeaders(ctx.Resp)
			}

			if !setting.AllowEmbedding {
				AddXFrameOptionsDenyHeader(w)
			}

			AddSecurityHeaders(w)
		})
	}
}

// AddSecurityHeaders adds various HTTP(S) response headers that enable various security protections behaviors in the client's browser.
func AddSecurityHeaders(w macaron.ResponseWriter) {
	if setting.Protocol == setting.HTTPS && setting.StrictTransportSecurity {
		strictHeaderValues := []string{fmt.Sprintf("max-age=%v", setting.StrictTransportSecurityMaxAge)}
		if setting.StrictTransportSecurityPreload {
			strictHeaderValues = append(strictHeaderValues, "preload")
		}
		if setting.StrictTransportSecuritySubDomains {
			strictHeaderValues = append(strictHeaderValues, "includeSubDomains")
		}
		w.Header().Add("Strict-Transport-Security", strings.Join(strictHeaderValues, "; "))
	}

	if setting.ContentTypeProtectionHeader {
		w.Header().Add("X-Content-Type-Options", "nosniff")
	}

	if setting.XSSProtectionHeader {
		w.Header().Add("X-XSS-Protection", "1; mode=block")
	}
}

func AddNoCacheHeaders(w macaron.ResponseWriter) {
	w.Header().Add("Cache-Control", "no-cache")
	w.Header().Add("Pragma", "no-cache")
	w.Header().Add("Expires", "-1")
}

func AddXFrameOptionsDenyHeader(w macaron.ResponseWriter) {
	w.Header().Add("X-Frame-Options", "deny")
}

`
	Quota = `
package middleware

import (
	"fmt"

	"gopkg.in/macaron.v1"

	m "{{.Dir}}/pkg/models"
	"{{.Dir}}/pkg/services/quota"
)

// Quota returns a function that returns a function used to call quotaservice based on target name
func Quota(quotaService *quota.QuotaService) func(target string) macaron.Handler {
	//https://open.spotify.com/track/7bZSoBEAEEUsGEuLOf94Jm?si=T1Tdju5qRSmmR0zph_6RBw fuuuuunky
	return func(target string) macaron.Handler {
		return func(c *m.ReqContext) {
			limitReached, err := quotaService.QuotaReached(c, target)
			if err != nil {
				c.JsonApiErr(500, "failed to get quota", err)
				return
			}
			if limitReached {
				c.JsonApiErr(403, fmt.Sprintf("%s Quota reached", target), nil)
				return
			}
		}
	}
}

`
	Recovery = `
// Copyright 2013 Martini Authors
// Copyright 2014 The Macaron Authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package middleware

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"runtime"

	"gopkg.in/macaron.v1"

	"{{.Dir}}/pkg/infra/log"
	m "{{.Dir}}/pkg/models"
	"{{.Dir}}/pkg/setting"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

// stack returns a nicely formatted stack frame, skipping skip frames
func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 {
		name = name[lastslash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}

// Recovery returns a middleware that recovers from any panics and writes a 500 if there was one.
// While Martini is in development mode, Recovery will also output the panic as HTML.
func Recovery() macaron.Handler {
	return func(c *macaron.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := stack(3)

				panicLogger := log.Root
				// try to get request logger
				if ctx, ok := c.Data["ctx"]; ok {
					ctxTyped := ctx.(*m.ReqContext)
					panicLogger = ctxTyped.Logger
				}

				panicLogger.Error("Request error", "error", err, "stack", string(stack))

				c.Data["Title"] = "Server Error"
				c.Data["AppSubUrl"] = setting.AppSubUrl
				c.Data["Theme"] = setting.DefaultTheme

				if setting.Env == setting.DEV {
					if theErr, ok := err.(error); ok {
						c.Data["Title"] = theErr.Error()
					}

					c.Data["ErrorMsg"] = string(stack)
				}

				ctx, ok := c.Data["ctx"].(*m.ReqContext)

				if ok && ctx.IsApiRequest() {
					resp := make(map[string]interface{})
					resp["message"] = "Internal Server Error - Check the Grafana server logs for the detailed error message."

					if c.Data["ErrorMsg"] != nil {
						resp["error"] = fmt.Sprintf("%v - %v", c.Data["Title"], c.Data["ErrorMsg"])
					} else {
						resp["error"] = c.Data["Title"]
					}

					c.JSON(500, resp)
				} else {
					c.HTML(500, setting.ERR_TEMPLATE_NAME)
				}
			}
		}()

		c.Next()
	}
}

`
	RequestMetric = `
package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"{{.Dir}}/pkg/infra/metrics"
	"gopkg.in/macaron.v1"
)

func RequestMetrics(handler string) macaron.Handler {
	return func(res http.ResponseWriter, req *http.Request, c *macaron.Context) {
		rw := res.(macaron.ResponseWriter)
		now := time.Now()
		c.Next()

		status := rw.Status()

		code := sanitizeCode(status)
		method := sanitizeMethod(req.Method)
		metrics.M_Http_Request_Total.WithLabelValues(handler, code, method).Inc()
		duration := time.Since(now).Nanoseconds() / int64(time.Millisecond)
		metrics.M_Http_Request_Summary.WithLabelValues(handler, code, method).Observe(float64(duration))

		if strings.HasPrefix(req.RequestURI, "/api/datasources/proxy") {
			countProxyRequests(status)
		} else if strings.HasPrefix(req.RequestURI, "/api/") {
			countApiRequests(status)
		} else {
			countPageRequests(status)
		}
	}
}

func countApiRequests(status int) {
	switch status {
	case 200:
		metrics.M_Api_Status.WithLabelValues("200").Inc()
	case 404:
		metrics.M_Api_Status.WithLabelValues("404").Inc()
	case 500:
		metrics.M_Api_Status.WithLabelValues("500").Inc()
	default:
		metrics.M_Api_Status.WithLabelValues("unknown").Inc()
	}
}

func countPageRequests(status int) {
	switch status {
	case 200:
		metrics.M_Page_Status.WithLabelValues("200").Inc()
	case 404:
		metrics.M_Page_Status.WithLabelValues("404").Inc()
	case 500:
		metrics.M_Page_Status.WithLabelValues("500").Inc()
	default:
		metrics.M_Page_Status.WithLabelValues("unknown").Inc()
	}
}

func countProxyRequests(status int) {
	switch status {
	case 200:
		metrics.M_Proxy_Status.WithLabelValues("200").Inc()
	case 404:
		metrics.M_Proxy_Status.WithLabelValues("400").Inc()
	case 500:
		metrics.M_Proxy_Status.WithLabelValues("500").Inc()
	default:
		metrics.M_Proxy_Status.WithLabelValues("unknown").Inc()
	}
}

func sanitizeMethod(m string) string {
	switch m {
	case "GET", "get":
		return "get"
	case "PUT", "put":
		return "put"
	case "HEAD", "head":
		return "head"
	case "POST", "post":
		return "post"
	case "DELETE", "delete":
		return "delete"
	case "CONNECT", "connect":
		return "connect"
	case "OPTIONS", "options":
		return "options"
	case "NOTIFY", "notify":
		return "notify"
	default:
		return strings.ToLower(m)
	}
}

// If the wrapped http.Handler has not set a status code, i.e. the value is
// currently 0, santizeCode will return 200, for consistency with behavior in
// the stdlib.
func sanitizeCode(s int) string {
	switch s {
	case 100:
		return "100"
	case 101:
		return "101"

	case 200, 0:
		return "200"
	case 201:
		return "201"
	case 202:
		return "202"
	case 203:
		return "203"
	case 204:
		return "204"
	case 205:
		return "205"
	case 206:
		return "206"

	case 300:
		return "300"
	case 301:
		return "301"
	case 302:
		return "302"
	case 304:
		return "304"
	case 305:
		return "305"
	case 307:
		return "307"

	case 400:
		return "400"
	case 401:
		return "401"
	case 402:
		return "402"
	case 403:
		return "403"
	case 404:
		return "404"
	case 405:
		return "405"
	case 406:
		return "406"
	case 407:
		return "407"
	case 408:
		return "408"
	case 409:
		return "409"
	case 410:
		return "410"
	case 411:
		return "411"
	case 412:
		return "412"
	case 413:
		return "413"
	case 414:
		return "414"
	case 415:
		return "415"
	case 416:
		return "416"
	case 417:
		return "417"
	case 418:
		return "418"

	case 500:
		return "500"
	case 501:
		return "501"
	case 502:
		return "502"
	case 503:
		return "503"
	case 504:
		return "504"
	case 505:
		return "505"

	case 428:
		return "428"
	case 429:
		return "429"
	case 431:
		return "431"
	case 511:
		return "511"

	default:
		return strconv.Itoa(s)
	}
}

`
	RequestTracing = `
package middleware

import (
	"fmt"
	"net/http"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	"gopkg.in/macaron.v1"
)

func RequestTracing(handler string) macaron.Handler {
	return func(res http.ResponseWriter, req *http.Request, c *macaron.Context) {
		rw := res.(macaron.ResponseWriter)

		tracer := opentracing.GlobalTracer()
		wireContext, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
		span := tracer.StartSpan(fmt.Sprintf("HTTP %s", handler), ext.RPCServerOption(wireContext))
		defer span.Finish()

		ctx := opentracing.ContextWithSpan(req.Context(), span)
		c.Req.Request = req.WithContext(ctx)

		c.Next()

		status := rw.Status()

		ext.HTTPStatusCode.Set(span, uint16(status))
		ext.HTTPUrl.Set(span, req.RequestURI)
		ext.HTTPMethod.Set(span, req.Method)
		if status >= 400 {
			ext.Error.Set(span, true)
		}
	}
}

`
)
