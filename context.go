package shack

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"net/http"
	"strings"
)

type Context struct {
	Router     *Router
	StatusCode int
	Writer     http.ResponseWriter
	Request    *http.Request
	Path       string
	Method     string
	Params     map[string]string
	handlers   []HandlerFunc
	index      int8
}

const abortIndex int8 = math.MaxInt8 / 2


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


func(c *Context) Header(key string, value string) *Context {
	c.Writer.Header().Set(key, value)
	return c
}


func(c *Context) String(s ...string) *Context {
	c.Header("Content-Type", "text/plain")
	c.Writer.Write([]byte(strings.Join(s, "")))
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


func(c *Context) Body() *bodyFlow {
	b, _ := ioutil.ReadAll(c.Request.Body)
	return newBodyFlow(b)
}


func(c *Context) Form(key string) *valueFlow {
	return newValueFlow(c.Request.FormValue(key))
}


func(c *Context) Forms() *formFlow {
	err := c.Request.ParseMultipartForm(1024*1024*1024) // 10Mb
	if err != nil {
		return newFormFlow(nil)
	}
	return newFormFlow(c.Request.MultipartForm.Value)

}


func(c *Context) Query(key string) *valueFlow {
	return newValueFlow(c.Request.URL.Query().Get(key))
}


func(c *Context) RawQuery() *valueFlow {
	f := newValueFlow(c.Request.URL.RawQuery)
	return f
}


func(c *Context) Set(key string, value string) {
	c.Params[key] = value
}


func(c *Context) Get(key string) string {
	return c.Params[key]
}


func(c *Context) Abort() {
	c.index = abortIndex
}


func(c *Context) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}
