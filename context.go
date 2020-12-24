package shack

import (
	"encoding/json"
	"net/http"
)

type Context struct {
	Router     *Router
	StatusCode int
	Writer     http.ResponseWriter
	Request    *http.Request
	Path       string
	Method     string
	Params     map[string]string
	Query      map[string]string
	handlers   []HandlerFunc
	index      int8
}


func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer : w,
		Request: r,
		Path   : r.URL.Path,
		Method : r.Method,
		index  : -1,
	}
}


func(c *Context) Param(key string) string {
	return c.Params[key]
}


func(c *Context) Status(code int) *Context {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
	return c
}


func (c *Context) Header(key string, value string) *Context {
	c.Writer.Header().Set(key, value)
	return c
}


func (c *Context) String(data string) *Context {
	c.Header("Content-Type", "text/plain")
	c.Writer.Write([]byte(data))
	return c
}


func(c *Context) JSON(data interface{}) *Context {
	c.Header("Content-Type", "application/json")
	b, _ := json.Marshal(data)
	c.Writer.Write(b)
	return c
}


func (c *Context) Data(data []byte) *Context {
	c.Writer.Write(data)
	return c
}


func(c *Context) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}
