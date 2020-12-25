package shack

import (
	"fmt"
	"net/http"
	"strings"
)


const (
	GET     = "GET"
	POST    = "POST"
	DELETE  = "DELETE"
	PUT     = "PUT"
	PATCH   = "PATCH"
	OPTIONS = "OPTIONS"
	HEAD    = "HEAD"
	ALL     = "ALL"
)

type HandlerFunc func(*Context)

type Router struct {
	sub                     map[string]*Router
	trie                    *trie
	groupMiddlewares        []HandlerFunc
	notFountHandler         HandlerFunc
	methodNotAllowedHandler HandlerFunc
}


func NewRouter() *Router {
	return &Router{
		sub : make(map[string]*Router),
		trie: newTrie(),
	}
}


func DefaultRouter() *Router {
	r := NewRouter()
	//r.Use(middleware.Recovery())
	//r.Use(middleware.AccessLog())
	return r
}


func(r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	c.handlers = append(c.handlers, getGroupMiddlewares(r, c.Path)...)
	r.handler(c)
}


func getGroupMiddlewares(r *Router, path string) (middlewares []HandlerFunc) {
	middlewares = append(middlewares, r.groupMiddlewares...)
	for pattern, router := range r.sub {
		if strings.HasPrefix(path, pattern) {
			middlewares = append(middlewares, getGroupMiddlewares(router, strings.TrimPrefix(path, pattern))...)
		}
	}
	return
}


func(r *Router) handler(ctx *Context) {
	handler, params, ok := r.trie.search(ctx.Method, ctx.Path)
	if ok && handler != nil  {
		ctx.Params = params
		ctx.handlers = append(ctx.handlers, handler...)
		ctx.Next()
	} else if ok {
		ctx.Status(http.StatusMethodNotAllowed)
		if r.methodNotAllowedHandler != nil {
			r.methodNotAllowedHandler(ctx)
		}
	} else {
		ctx.Status(http.StatusNotFound)
		if r.notFountHandler != nil {
			r.notFountHandler(ctx)
		}
	}
}


func(r *Router) GET(pattern string, handler func(*Context)) *trie {
	return r.trie.insert(GET, pattern, handler)
}


func(r *Router) POST(pattern string, handler func(*Context)) *trie {
	return r.trie.insert(POST, pattern, handler)
}


func(r *Router) DELETE(pattern string, handler func(*Context)) *trie {
	return r.trie.insert(DELETE, pattern, handler)
}


func(r *Router) PUT(pattern string, handler func(*Context)) *trie {
	return r.trie.insert(PUT, pattern, handler)
}


func(r *Router) PATCH(pattern string, handler func(*Context)) *trie {
	return r.trie.insert(PATCH, pattern, handler)
}


func(r *Router) OPTIONS(pattern string, handler func(*Context)) *trie {
	return r.trie.insert(OPTIONS, pattern, handler)
}


func(r *Router) HEAD(pattern string, handler func(*Context)) *trie {
	return r.trie.insert(HEAD, pattern, handler)
}


func(r *Router) Use(middleware ...HandlerFunc) {
	r.groupMiddlewares = append(r.groupMiddlewares, middleware...)
}


func(r *Router) Mount(pattern string, router *Router) {
	if !isVaildPattern(pattern) {
		panic(fmt.Sprintf("shack: pattern %s is not valid", pattern))
	}

	if router == nil {
		panic(fmt.Sprintf("shack: handler is nil in mounting %s", pattern))
	}

	if r.sub[pattern] != nil {
		panic(fmt.Sprintf("shack: pattern %s to mount is already exist", pattern))
	}

	// todo: 冲突检测
	r.sub[pattern] = router
	r.trie.child[pattern[1:]] = router.trie
}


func(r *Router) Group(pattern string, fn func(r *Router)) *Router {
	if !isVaildPattern(pattern) {
		panic(fmt.Sprintf("shack: pattern %s is not valid", pattern))
	}

	if fn == nil {
		panic(fmt.Sprintf("shack: fn is nil in grouping %s", pattern))
	}

	// todo: 优化
	sub := NewRouter()
	fn(sub)
	if r.trie.child[pattern[1:]] != nil {
		for k, v := range sub.trie.child {
			r.trie.child[pattern[1:]].child[k] = v
		}
		return r
	}

	r.Mount(pattern, sub)
	return r
}


// Handle adds routes for `pattern` that matches all HTTP methods.
func(r *Router) Handle(pattern string, fn HandlerFunc) {
	r.trie.insert(ALL, pattern, fn)
}


// NotFound defines a handler to respond whenever a route could
// not be found.
func(r *Router) NotFound(fn HandlerFunc) {
	r.notFountHandler = fn
}

// MethodNotAllowed defines a handler to respond whenever a method is
// not allowed.
func(r *Router) MethodNotAllowed(fn HandlerFunc) {
	r.methodNotAllowedHandler = fn
}


/*

func(r *Router) Default() {
	r.handler.Use(middleware.RequestID)
	r.handler.Use(middleware.RealIP)
	r.handler.Use(middleware.Logger)
	r.handler.Use(middleware.Recoverer)
	r.handler.Use(middleware.SetHeader("Content-type", "application/json"))
	r.handler.MethodNotAllowed(r.handler.MethodNotAllowedHandler())
}
*/