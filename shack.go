package shack

import (
	"fmt"
	"net/http"
)


func Run(addr string, router *Router) {
	err := http.ListenAndServe(addr, router)
	if err != nil {
		panic(fmt.Sprint("shack: ", err))
	}
	return
}


// Logger returns a logger by specify a name
func Logger(name string) *logger {
	Log.name = name
	return Log
}


type _Router interface {
	ServeHTTP(w http.ResponseWriter, req *http.Request)

	GET(pattern string, handler HandlerFunc) *trie
	POST(pattern string, handler HandlerFunc) *trie
	DELETE(pattern string, handler HandlerFunc) *trie
	PUT(pattern string, handler HandlerFunc) *trie
	PATCH(pattern string, handler HandlerFunc) *trie
	OPTIONS(pattern string, handler HandlerFunc) *trie
	HEAD(pattern string, handler HandlerFunc) *trie
	Handle(pattern string, handler HandlerFunc)

	With(middleware ...HandlerFunc)
	Use(middlewares ...HandlerFunc)
	Mount(pattern string, router *Router)
	Group(pattern string, fn func(r *Router)) *Router

	NotFound(handler HandlerFunc)
	MethodNotAllowed(handler HandlerFunc)
}

type _Logger interface {
	Enable()

	Level(level int8) *logger
	Encoding(encoding string) *logger
	Output(paths ...string) *logger
	Dev() *logger

	Debug(msg string, keyAndValues ...interface{})
	Info(msg string, keyAndValues ...interface{})
	Warn(msg string, keyAndValues ...interface{})
	Error(msg string, keyAndValues ...interface{})
	Panic(msg string, keyAndValues ...interface{})
	Fatal(msg string, keyAndValues ...interface{})
}

type _Context interface {
	Status(code int) *Context
	Header(key string, value string) *Context
	String(s ...string) *Context
	JSON(data interface{}) *Context
	Data(data []byte) *Context

	Param(key string) string
	Body() *bodyFlow
	Form(key string) *valueFlow
	Forms() *formFlow
	Query(key string) *valueFlow
	RawQuery() *valueFlow

	Set(key string, value string)
	Get(key string) string

	Abort()
	Next()
}


type _Flow interface {
	Value()
	Int() int
	Float64() float64
	Bool() bool
	BindJson(dst interface{}) error
	Bind(dst interface{}, tag ...string) error
}
