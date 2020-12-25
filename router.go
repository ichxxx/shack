package shack

import (
	"fmt"
	"net/http"
	"strings"
)


const (
	_GET     = "GET"
	_POST    = "POST"
	_DELETE  = "DELETE"
	_PUT     = "PUT"
	_PATCH   = "PATCH"
	_OPTIONS = "OPTIONS"
	_HEAD    = "HEAD"
	_ALL     = "ALL"
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
	return r.trie.insert(_GET, pattern, handler)
}


func(r *Router) POST(pattern string, handler func(*Context)) *trie {
	return r.trie.insert(_POST, pattern, handler)
}


func(r *Router) DELETE(pattern string, handler func(*Context)) *trie {
	return r.trie.insert(_DELETE, pattern, handler)
}


func(r *Router) PUT(pattern string, handler func(*Context)) *trie {
	return r.trie.insert(_PUT, pattern, handler)
}


func(r *Router) PATCH(pattern string, handler func(*Context)) *trie {
	return r.trie.insert(_PATCH, pattern, handler)
}


func(r *Router) OPTIONS(pattern string, handler func(*Context)) *trie {
	return r.trie.insert(_OPTIONS, pattern, handler)
}


func(r *Router) HEAD(pattern string, handler func(*Context)) *trie {
	return r.trie.insert(_HEAD, pattern, handler)
}


func(r *Router) Use(middleware ...HandlerFunc) {
	r.groupMiddlewares = append(r.groupMiddlewares, middleware...)
}


func(r *Router) Mount(pattern string, router *Router) {
	if !isValidPattern(pattern) {
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
	if !isValidPattern(pattern) {
		panic(fmt.Sprintf("shack: pattern %s is not valid", pattern))
	}

	if fn == nil {
		panic(fmt.Sprintf("shack: fn is nil in grouping %s", pattern))
	}

	sub := NewRouter()
	fn(sub)
	if child, found := r.trie.child[pattern[1:]]; found {
		for p, t := range sub.trie.child {
			child.child[p] = t
		}
		return r
	}

	r.sub[pattern] = sub
	r.trie.child[pattern[1:]] = sub.trie
	return r
}


// Handle adds routes for `pattern` that matches all HTTP methods.
func(r *Router) Handle(pattern string, fn HandlerFunc) {
	r.trie.insert(_ALL, pattern, fn)
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