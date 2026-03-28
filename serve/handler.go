package serve

type Handler interface {
	ServeHTTP(c *Context)
}

type HandlerFunc func(c *Context)

func (f HandlerFunc) ServeHTTP(c *Context) {
	f(c)
}
