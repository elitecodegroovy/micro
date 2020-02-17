package template4app

var (
	RouterRegister = `
package routing

import (
	"net/http"
	"strings"

	"gopkg.in/macaron.v1"
)

type Router interface {
	Handle(method, pattern string, handlers []macaron.Handler) *macaron.Route
	Get(pattern string, handlers ...macaron.Handler) *macaron.Route
}

// RouteRegister allows you to add routes and macaron.Handlers
// that the web server should serve.
type RouteRegister interface {
	// Get adds a list of handlers to a given route with a GET HTTP verb
	Get(string, ...macaron.Handler)

	// Post adds a list of handlers to a given route with a POST HTTP verb
	Post(string, ...macaron.Handler)

	// Delete adds a list of handlers to a given route with a DELETE HTTP verb
	Delete(string, ...macaron.Handler)

	// Put adds a list of handlers to a given route with a PUT HTTP verb
	Put(string, ...macaron.Handler)

	// Patch adds a list of handlers to a given route with a PATCH HTTP verb
	Patch(string, ...macaron.Handler)

	// Any adds a list of handlers to a given route with any HTTP verb
	Any(string, ...macaron.Handler)

	// Group allows you to pass a function that can add multiple routes
	// with a shared prefix route.
	Group(string, func(RouteRegister), ...macaron.Handler)

	// Insert adds more routes to an existing Group.
	Insert(string, func(RouteRegister), ...macaron.Handler)

	// Register iterates over all routes added to the RouteRegister
	// and add them to the `+ `Router`+` pass as an parameter.
	Register(Router)
}

type RegisterNamedMiddleware func(name string) macaron.Handler

// NewRouteRegister creates a new RouteRegister with all middlewares sent as params
func NewRouteRegister(namedMiddleware ...RegisterNamedMiddleware) RouteRegister {
	return &routeRegister{
		prefix:          "",
		routes:          []route{},
		subfixHandlers:  []macaron.Handler{},
		namedMiddleware: namedMiddleware,
	}
}

type route struct {
	method   string
	pattern  string
	handlers []macaron.Handler
}

type routeRegister struct {
	prefix          string
	subfixHandlers  []macaron.Handler
	namedMiddleware []RegisterNamedMiddleware
	routes          []route
	groups          []*routeRegister
}

func (rr *routeRegister) Insert(pattern string, fn func(RouteRegister), handlers ...macaron.Handler) {

	//loop over all groups at current level
	for _, g := range rr.groups {

		// apply routes if the prefix matches the pattern
		if g.prefix == pattern {
			g.Group("", fn)
			break
		}

		// go down one level if the prefix can be find in the pattern
		if strings.HasPrefix(pattern, g.prefix) {
			g.Insert(pattern, fn)
		}
	}
}

func (rr *routeRegister) Group(pattern string, fn func(rr RouteRegister), handlers ...macaron.Handler) {
	group := &routeRegister{
		prefix:          rr.prefix + pattern,
		subfixHandlers:  append(rr.subfixHandlers, handlers...),
		routes:          []route{},
		namedMiddleware: rr.namedMiddleware,
	}

	fn(group)
	rr.groups = append(rr.groups, group)
}

func (rr *routeRegister) Register(router Router) {
	for _, r := range rr.routes {
		// GET requests have to be added to macaron routing using Get()
		// Otherwise HEAD requests will not be allowed.
		// https://github.com/go-macaron/macaron/blob/a325110f8b392bce3e5cdeb8c44bf98078ada3be/router.go#L198
		if r.method == http.MethodGet {
			router.Get(r.pattern, r.handlers...)
		} else {
			router.Handle(r.method, r.pattern, r.handlers)
		}
	}

	for _, g := range rr.groups {
		g.Register(router)
	}
}

func (rr *routeRegister) route(pattern, method string, handlers ...macaron.Handler) {
	h := make([]macaron.Handler, 0)
	for _, fn := range rr.namedMiddleware {
		h = append(h, fn(pattern))
	}

	h = append(h, rr.subfixHandlers...)
	h = append(h, handlers...)

	for _, r := range rr.routes {
		if r.pattern == rr.prefix+pattern && r.method == method {
			panic("cannot add duplicate route")
		}
	}

	rr.routes = append(rr.routes, route{
		method:   method,
		pattern:  rr.prefix + pattern,
		handlers: h,
	})
}

func (rr *routeRegister) Get(pattern string, handlers ...macaron.Handler) {
	rr.route(pattern, http.MethodGet, handlers...)
}

func (rr *routeRegister) Post(pattern string, handlers ...macaron.Handler) {
	rr.route(pattern, http.MethodPost, handlers...)
}

func (rr *routeRegister) Delete(pattern string, handlers ...macaron.Handler) {
	rr.route(pattern, http.MethodDelete, handlers...)
}

func (rr *routeRegister) Put(pattern string, handlers ...macaron.Handler) {
	rr.route(pattern, http.MethodPut, handlers...)
}

func (rr *routeRegister) Patch(pattern string, handlers ...macaron.Handler) {
	rr.route(pattern, http.MethodPatch, handlers...)
}

func (rr *routeRegister) Any(pattern string, handlers ...macaron.Handler) {
	rr.route(pattern, "*", handlers...)
}

`
	RouterRegisterTest = `
package routing

import (
	"net/http"
	"strconv"
	"testing"

	"gopkg.in/macaron.v1"
)

type fakeRouter struct {
	route []route
}

func (fr *fakeRouter) Handle(method, pattern string, handlers []macaron.Handler) *macaron.Route {
	fr.route = append(fr.route, route{
		pattern:  pattern,
		method:   method,
		handlers: handlers,
	})

	return &macaron.Route{}
}

func (fr *fakeRouter) Get(pattern string, handlers ...macaron.Handler) *macaron.Route {
	fr.route = append(fr.route, route{
		pattern:  pattern,
		method:   http.MethodGet,
		handlers: handlers,
	})

	return &macaron.Route{}
}

func emptyHandlers(n int) []macaron.Handler {
	var res []macaron.Handler
	for i := 1; n >= i; i++ {
		res = append(res, emptyHandler(strconv.Itoa(i)))
	}
	return res
}

func emptyHandler(name string) macaron.Handler {
	return struct{ name string }{name: name}
}

func TestRouteSimpleRegister(t *testing.T) {
	testTable := []route{
		{method: "DELETE", pattern: "/admin", handlers: emptyHandlers(2)},
		{method: "GET", pattern: "/down", handlers: emptyHandlers(3)},
	}

	// Setup
	rr := NewRouteRegister(func(name string) macaron.Handler {
		return emptyHandler(name)
	})

	rr.Delete("/admin", emptyHandler("1"))
	rr.Get("/down", emptyHandler("1"), emptyHandler("2"))

	fr := &fakeRouter{}
	rr.Register(fr)

	// Validation
	if len(fr.route) != len(testTable) {
		t.Errorf("want %v routes, got %v", len(testTable), len(fr.route))
	}

	for i := range testTable {
		if testTable[i].method != fr.route[i].method {
			t.Errorf("want %s got %v", testTable[i].method, fr.route[i].method)
		}

		if testTable[i].pattern != fr.route[i].pattern {
			t.Errorf("want %s got %v", testTable[i].pattern, fr.route[i].pattern)
		}

		if len(testTable[i].handlers) != len(fr.route[i].handlers) {
			t.Errorf("want %d handlers got %d handlers \ntestcase: %v\nroute: %v\n",
				len(testTable[i].handlers),
				len(fr.route[i].handlers),
				testTable[i],
				fr.route[i])
		}
	}
}

func TestRouteGroupedRegister(t *testing.T) {
	testTable := []route{
		{method: "DELETE", pattern: "/admin", handlers: emptyHandlers(1)},
		{method: "GET", pattern: "/down", handlers: emptyHandlers(2)},
		{method: "POST", pattern: "/user", handlers: emptyHandlers(1)},
		{method: "PUT", pattern: "/user/friends", handlers: emptyHandlers(1)},
		{method: "DELETE", pattern: "/user/admin", handlers: emptyHandlers(2)},
		{method: "GET", pattern: "/user/admin/all", handlers: emptyHandlers(4)},
	}

	// Setup
	rr := NewRouteRegister()

	rr.Delete("/admin", emptyHandler("1"))
	rr.Get("/down", emptyHandler("1"), emptyHandler("2"))

	rr.Group("/user", func(user RouteRegister) {
		user.Post("", emptyHandler("1"))
		user.Put("/friends", emptyHandler("2"))

		user.Group("/admin", func(admin RouteRegister) {
			admin.Delete("", emptyHandler("3"))
			admin.Get("/all", emptyHandler("3"), emptyHandler("4"), emptyHandler("5"))

		}, emptyHandler("3"))
	})

	fr := &fakeRouter{}
	rr.Register(fr)

	// Validation
	if len(fr.route) != len(testTable) {
		t.Errorf("want %v routes, got %v", len(testTable), len(fr.route))
	}

	for i := range testTable {
		if testTable[i].method != fr.route[i].method {
			t.Errorf("want %s got %v", testTable[i].method, fr.route[i].method)
		}

		if testTable[i].pattern != fr.route[i].pattern {
			t.Errorf("want %s got %v", testTable[i].pattern, fr.route[i].pattern)
		}

		if len(testTable[i].handlers) != len(fr.route[i].handlers) {
			t.Errorf("want %d handlers got %d handlers \ntestcase: %v\nroute: %v\n",
				len(testTable[i].handlers),
				len(fr.route[i].handlers),
				testTable[i],
				fr.route[i])
		}
	}
}
func TestRouteGroupInserting(t *testing.T) {
	testTable := []route{
		{method: http.MethodGet, pattern: "/api/", handlers: emptyHandlers(1)},
		{method: http.MethodPost, pattern: "/api/group/endpoint", handlers: emptyHandlers(1)},

		{method: http.MethodGet, pattern: "/api/group/inserted", handlers: emptyHandlers(1)},
		{method: http.MethodDelete, pattern: "/api/inserted-endpoint", handlers: emptyHandlers(1)},
	}

	// Setup
	rr := NewRouteRegister()

	rr.Group("/api", func(api RouteRegister) {
		api.Get("/", emptyHandler("1"))

		api.Group("/group", func(group RouteRegister) {
			group.Post("/endpoint", emptyHandler("1"))
		})
	})

	rr.Insert("/api", func(api RouteRegister) {
		api.Delete("/inserted-endpoint", emptyHandler("1"))
	})

	rr.Insert("/api/group", func(group RouteRegister) {
		group.Get("/inserted", emptyHandler("1"))
	})

	fr := &fakeRouter{}
	rr.Register(fr)

	// Validation
	if len(fr.route) != len(testTable) {
		t.Fatalf("want %v routes, got %v", len(testTable), len(fr.route))
	}

	for i := range testTable {
		if testTable[i].method != fr.route[i].method {
			t.Errorf("want %s got %v", testTable[i].method, fr.route[i].method)
		}

		if testTable[i].pattern != fr.route[i].pattern {
			t.Errorf("want %s got %v", testTable[i].pattern, fr.route[i].pattern)
		}

		if len(testTable[i].handlers) != len(fr.route[i].handlers) {
			t.Errorf("want %d handlers got %d handlers \ntestcase: %v\nroute: %v\n",
				len(testTable[i].handlers),
				len(fr.route[i].handlers),
				testTable[i],
				fr.route[i])
		}
	}
}

func TestDuplicateRoutShouldPanic(t *testing.T) {
	defer func() {
		if recover() != "cannot add duplicate route" {
			t.Errorf("Should cause panic if duplicate routes are added ")
		}
	}()

	rr := NewRouteRegister(func(name string) macaron.Handler {
		return emptyHandler(name)
	})

	rr.Get("/api", emptyHandler("1"))
	rr.Get("/api", emptyHandler("1"))

	fr := &fakeRouter{}
	rr.Register(fr)
}
func TestNamedMiddlewareRouteRegister(t *testing.T) {
	testTable := []route{
		{method: "DELETE", pattern: "/admin", handlers: emptyHandlers(2)},
		{method: "GET", pattern: "/down", handlers: emptyHandlers(3)},
		{method: "POST", pattern: "/user", handlers: emptyHandlers(2)},
		{method: "PUT", pattern: "/user/friends", handlers: emptyHandlers(2)},
		{method: "DELETE", pattern: "/user/admin", handlers: emptyHandlers(3)},
		{method: "GET", pattern: "/user/admin/all", handlers: emptyHandlers(5)},
	}

	// Setup
	rr := NewRouteRegister(func(name string) macaron.Handler {
		return emptyHandler(name)
	})

	rr.Delete("/admin", emptyHandler("1"))
	rr.Get("/down", emptyHandler("1"), emptyHandler("2"))

	rr.Group("/user", func(user RouteRegister) {
		user.Post("", emptyHandler("1"))
		user.Put("/friends", emptyHandler("2"))

		user.Group("/admin", func(admin RouteRegister) {
			admin.Delete("", emptyHandler("3"))
			admin.Get("/all", emptyHandler("3"), emptyHandler("4"), emptyHandler("5"))

		}, emptyHandler("3"))
	})

	fr := &fakeRouter{}
	rr.Register(fr)

	// Validation
	if len(fr.route) != len(testTable) {
		t.Errorf("want %v routes, got %v", len(testTable), len(fr.route))
	}

	for i := range testTable {
		if testTable[i].method != fr.route[i].method {
			t.Errorf("want %s got %v", testTable[i].method, fr.route[i].method)
		}

		if testTable[i].pattern != fr.route[i].pattern {
			t.Errorf("want %s got %v", testTable[i].pattern, fr.route[i].pattern)
		}

		if len(testTable[i].handlers) != len(fr.route[i].handlers) {
			t.Errorf("want %d handlers got %d handlers \ntestcase: %v\nroute: %v\n",
				len(testTable[i].handlers),
				len(fr.route[i].handlers),
				testTable[i],
				fr.route[i])
		}
	}
}

`


)
