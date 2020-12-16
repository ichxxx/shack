package shack

import (
	"net/http"
)


const (
	GET     = "GET"
	POST    = "POST"
	DELETE  = "DELETE"
	PUT     = "PUT"
	PATCH   = "PATCH"
	OPTIONS = "OPTIONS"
	HEAD    = "HEAD"
)

type HandlerFunc func(*Context)

type Router struct {
	trie *trie
}


func NewRouter() *Router {
	return &Router{
		trie: newTrie(),
	}
}


func(r *Router) ServeHTTP(w http.ResponseWriter, _r *http.Request) {
	c := newContext(w, _r)
	r.handler(c)
}


func(r *Router) handler(ctx *Context) {
	handler, params, ok := r.trie.search(ctx.Method, ctx.Path)
	if ok && handler != nil  {
		ctx.Params = params
		handler(ctx)
	} else if ok {
		ctx.Status(http.StatusMethodNotAllowed)
	} else {
		ctx.Status(http.StatusNotFound)
	}
}


func(r *Router) GET(pattern string, handler func(*Context)) {
	r.trie.insert(GET, pattern, handler)
}


func(r *Router) POST(pattern string, handler func(*Context)) {
	r.trie.insert(POST, pattern, handler)
}


func(r *Router) DELETE(pattern string, handler func(*Context)) {
	r.trie.insert(DELETE, pattern, handler)
}


func(r *Router) PUT(pattern string, handler func(*Context)) {
	r.trie.insert(PUT, pattern, handler)
}


/*
func NewRouters() *Router {
	return &Router{
		handler: chi.NewRouter(),
	}
}


func(r *Router) Default() {
	r.handler.Use(middleware.RequestID)
	r.handler.Use(middleware.RealIP)
	r.handler.Use(middleware.Logger)
	r.handler.Use(middleware.Recoverer)
	r.handler.Use(middleware.SetHeader("Content-type", "application/json"))
	r.handler.MethodNotAllowed(r.handler.MethodNotAllowedHandler())
}


func(r *Router) Get(pattern string, handlerFunc http.HandlerFunc) {
	r.handler.Get(pattern, handlerFunc)
}


func(r *Router) Post(pattern string, handlerFunc http.HandlerFunc) {
	r.handler.Post(pattern, handlerFunc)
}


func(r *Router) Put(pattern string, handlerFunc http.HandlerFunc) {
	r.handler.Put(pattern, handlerFunc)
}


func(r *Router) Delete(pattern string, handlerFunc http.HandlerFunc) {
	r.handler.Delete(pattern, handlerFunc)
}


func(r *Router) With(middleware func(http.Handler)http.Handler) {
	r.handler.With(middleware)
}


func(r *Router) Mount(pattern string, handler http.Handler) {
	subRoute := chi.NewRouter()
	r.handler.Mount(pattern, subRoute)
}
*/