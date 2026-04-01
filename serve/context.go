package serve

import (
	"context"
	"net/http"
)

type Context struct {
	context.Context

	r *http.Request
	w http.ResponseWriter

	params []string
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	c := &Context{
		r: r,
		w: w,
	}
	return c
}

func (c *Context) Method() string {
	return c.r.Method
}

func (c *Context) Header(key string) string {
	return c.r.Header.Get(key)
}

func (c *Context) SetHeader(key, value string) {
	c.r.Header.Set(key, value)
}

func (c *Context) WriteHeader(status int) {
	c.w.WriteHeader(status)
}

func (c *Context) Cookie(key string) (*http.Cookie, error) {
	return c.r.Cookie(key)
}

func (c *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.w, cookie)
}

func (c *Context) Query(key string) string {
	return c.r.URL.Query().Get(key)
}

func (c *Context) Param(key string) string {
	return c.r.PathValue(key)
}

func (c *Context) JSON(status int, data J) {
	writeJSON(c.w, status, data)
}

func (c *Context) Bind(obj any) error {
	return bindJSON(c.r, obj)
}
