package serve

// Group maintains data for a router group.
// Middleware assigned to a group is only applied to routes within the group.
// TODO(RM): Validate middleware is limited to group
// All route paths begin with the assigned basePath.
type Group struct {
	w *Serve

	basePath string
}

func (g *Group) GET(pattern string, handler HandlerFunc) {
	g.w.GET(g.basePath+pattern, handler)
}

func (g *Group) POST(pattern string, handler HandlerFunc) {
	g.w.POST(g.basePath+pattern, handler)
}

func (g *Group) PUT(pattern string, handler HandlerFunc) {
	g.w.PUT(g.basePath+pattern, handler)
}

func (g *Group) DELETE(pattern string, handler HandlerFunc) {
	g.w.DELETE(g.basePath+pattern, handler)
}

func (g *Group) Use(middleware ...Middleware) {
	g.w.Use(middleware...)
}

func (g *Group) Group(path string) *Group {
	return g.w.Group(g.basePath + path)
}
