package shack

import (
	"encoding/json"
	"net/http"
)

type Context struct {
	Router  *Router
	Writer  http.ResponseWriter
	Request *http.Request
	Path    string
	Method  string
	Params  map[string]string
	Query   map[string]string
}


func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer:  w,
		Request: r,
		Path:    r.URL.Path,
		Method:  r.Method,
	}
}


func(c *Context) Param(key string) string {
	return c.Params[key]
}


func(c *Context) Status(code int) {
	c.Writer.WriteHeader(code)
}


func(c *Context) JSON(data interface{}) {
	b, _ := json.Marshal(data)
	c.Writer.Write(b)
}
