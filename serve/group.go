package serve

import "strings"

// Group maintains data for a router group.
// Middleware assigned to a group is only applied to routes within the group.
// TODO(RM): Validate middleware is limited to group
// All route paths begin with the assigned basePath.
type Group struct {
	w *Serve

	basePath    string
	middlewares []Middleware
}

func (g *Group) chain(handler Handler) Handler {
	for i := len(g.middlewares) - 1; i >= 0; i-- {
		handler = g.middlewares[i](handler)
	}
	return handler
}

func (g *Group) Use(middleware ...Middleware) {
	g.middlewares = append(g.middlewares, middleware...)
}

func (g *Group) Handle(pattern string, handler Handler, methods ...string) {
	g.w.Handle(g.basePath+pattern, g.chain(handler), methods...)
}

func (g *Group) GET(pattern string, handler HandlerFunc) {
	g.Handle(pattern, handler, "GET")
}

func (g *Group) POST(pattern string, handler HandlerFunc) {
	g.Handle(pattern, handler, "POST")
}

func (g *Group) PUT(pattern string, handler HandlerFunc) {
	g.Handle(pattern, handler, "PUT")
}

func (g *Group) DELETE(pattern string, handler HandlerFunc) {
	g.Handle(pattern, handler, "DELETE")
}

func (g *Group) Group(path string) *Group {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	child := &Group{
		w:        g.w,
		basePath: g.basePath + path,
	}
	child.middlewares = append(child.middlewares, g.middlewares...)
	return child
}
