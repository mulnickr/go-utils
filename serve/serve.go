/* Package serve is a lightweight net/http wrapper implementing basic routing capabilities and:
*  - Middleware
*   - Nested route groups
*   - Request method mapping
*
* Based largely on Flow: github.com/alexedwards/flow
*
* Example usage:
*
* router := serve.Default()
* router.Use(DefaultAuth) // auth middleware
* api = := router.Group("/api/v1") // API group with base path "/api/v1"
* api.GET("/health", checkHealth) // GET endpoint for "/health"
* router.ListenAndServe(":5000") // Initialize server
*
* func DefaultAuth(next serve.Handler) serve.HandlerFunc {
*     return serve.HandlerFunc(func(c *serve.Context) {
*         c.JSON(http.StatusOK, serve.J{"msg": "everyone is authorized!"})
*         next.ServeHTTP(c) // continue the handler chain
*     })
* }
*
*
 */
package serve

import (
	"net/http"
	"slices"
	"strings"
)

type Serve struct {
	NotFound         Handler
	MethodNotAllowed Handler
	Options          Handler

	routes      *[]route
	middlewares []Middleware
}

// New returns a new `Serve` instance with all specified data
// nf: NotFound - HandlerFunc
// mna: MethodNotAllowed - HandlerFunc
// opt: Options - HandlerFunc
func New(nf, mna, opt HandlerFunc, routes *[]route) *Serve {
	return &Serve{
		NotFound:         nf,
		MethodNotAllowed: mna,
		Options:          opt,
		routes:           routes,
	}
}

func notFound(c *Context) {
	http.Error(c.w, "404: you goofed", http.StatusNotFound)
}

func Default() *Serve {
	return &Serve{
		NotFound: HandlerFunc(notFound),
		MethodNotAllowed: HandlerFunc(func(c *Context) {
			http.Error(c.w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}),
		Options: HandlerFunc(func(c *Context) {
			c.WriteHeader(http.StatusNoContent)
		}),
		routes: &[]route{},
	}
}

func (x *Serve) ListenAndServe(host string) {
	http.ListenAndServe(host, x)
}

func (x *Serve) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	us := strings.Split(r.URL.EscapedPath(), "/")
	am := []string{}

	for _, route := range *x.routes {

		ctx, ok := route.match(r.Context(), r, us)
		if ok {
			if r.Method == route.method {
				route.handler.ServeHTTP(NewContext(w, r.WithContext(ctx)))
				return
			}
			if !slices.Contains(am, route.method) {
				am = append(am, route.method)
			}
		}
	}

	if len(am) > 0 {
		w.Header().Set("Allow", strings.Join(append(am, http.MethodOptions), ", "))
		if r.Method == http.MethodOptions {
			x.chain(x.Options).ServeHTTP(NewContext(w, r))
		} else {
			x.chain(x.MethodNotAllowed).ServeHTTP(NewContext(w, r))
		}
		return
	}

	x.chain(x.NotFound).ServeHTTP(NewContext(w, r))
}

var routesLog = map[string][]string{}

func (x *Serve) Handle(pattern string, handler Handler, methods ...string) {
	if slices.Contains(methods, http.MethodGet) && !slices.Contains(methods, http.MethodHead) {
		methods = append(methods, http.MethodHead)
	}

	if len(methods) == 0 {
		methods = AllMethods
	}

	for _, method := range methods {
		route := route{
			method:   strings.ToUpper(method),
			segments: strings.Split(pattern, "/"),
			handler:  x.chain(handler),
		}
		routesLog[pattern] = append(routesLog[pattern], route.method)
		*x.routes = append(*x.routes, route)
	}
}

func (x *Serve) HandleFunc(pattern string, fn HandlerFunc, methods ...string) {
	x.Handle(pattern, fn, methods...)
}

func (x *Serve) chain(handler Handler) Handler {
	for i := len(x.middlewares) - 1; i >= 0; i-- {
		handler = x.middlewares[i](handler)
	}
	return handler
}

func (x *Serve) Use(middleware ...Middleware) {
	x.middlewares = append(x.middlewares, middleware...)
}

func (x *Serve) Group(path string) *Group {
	group := &Group{w: x, basePath: path}
	return group
}

func (x *Serve) GET(pattern string, handler HandlerFunc) {
	x.HandleFunc(pattern, handler, "GET")
}

func (x *Serve) POST(pattern string, handler HandlerFunc) {
	x.HandleFunc(pattern, handler, "POST")
}

func (x *Serve) PUT(pattern string, handler HandlerFunc) {
	x.HandleFunc(pattern, handler, "PUT")
}

func (x *Serve) DELETE(pattern string, handler HandlerFunc) {
	x.HandleFunc(pattern, handler, "DELETE")
}
