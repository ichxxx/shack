package shack

import (
	"fmt"
	"net/http"
	"strings"
)


const (
	_GET     = http.MethodGet
	_POST    = http.MethodPost
	_DELETE  = http.MethodDelete
	_PUT     = http.MethodPut
	_PATCH   = http.MethodPatch
	_OPTIONS = http.MethodOptions
	_HEAD    = http.MethodHead
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
	if path == "/" {
		return
	}

	path = path[1:]
	for pattern, router := range r.sub {
		if strings.HasPrefix(path, pattern) {
			nextPath := strings.TrimPrefix(path, pattern)
			if nextPath != "" {
				middlewares = append(middlewares, getGroupMiddlewares(router, nextPath)...)
			}
		}
	}
	return
}


func(r *Router) handler(ctx *Context) {
	handlers, params, ok := r.trie.search(ctx.Method, ctx.Path)
	if ok && len(handlers) > 0  {
		ctx.Params = params
		ctx.handlers = append(ctx.handlers, handlers...)
		ctx.Next()
	} else if ok {
		// todo: 如果是模糊节点，且没有handler，会通过判断，待修复
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


func(r *Router) GET(pattern string, handler HandlerFunc) *trie {
	return r.trie.insert(pattern, handler, _GET)
}


func(r *Router) POST(pattern string, handler HandlerFunc) *trie {
	return r.trie.insert(pattern, handler, _POST)
}


func(r *Router) DELETE(pattern string, handler HandlerFunc) *trie {
	return r.trie.insert(pattern, handler, _DELETE)
}


func(r *Router) PUT(pattern string, handler HandlerFunc) *trie {
	return r.trie.insert(pattern, handler, _PUT)
}


func(r *Router) PATCH(pattern string, handler HandlerFunc) *trie {
	return r.trie.insert(pattern, handler, _PATCH)
}


func(r *Router) OPTIONS(pattern string, handler HandlerFunc) *trie {
	return r.trie.insert(pattern, handler, _OPTIONS)
}


func(r *Router) HEAD(pattern string, handler HandlerFunc) *trie {
	return r.trie.insert(pattern, handler, _HEAD)
}


// Use appends one or more middlewares onto the router.
func(r *Router) Use(middlewares ...HandlerFunc) {
	r.groupMiddlewares = append(r.groupMiddlewares, middlewares...)
}


// Mount attaches another router along a `pattern` string.
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

	// todo: should use mergeSubRouter and mergeSubTrie ?
	r.sub[pattern] = router
	r.trie.childs[pattern[1:]] = router.trie
}


// Group adds a sub-Router to the group along a `pattern` string.
func(r *Router) Group(pattern string, fn func(r *Router)) *Router {
	if !isValidPattern(pattern) {
		panic(fmt.Sprintf("shack: pattern %s is not valid", pattern))
	}

	if fn == nil {
		panic(fmt.Sprintf("shack: fn is nil in grouping %s", pattern))
	}

	root := r
	sub := NewRouter()
	fn(sub)

	if pattern == "/" {
		for method, handler := range sub.trie.handlers {
			root.trie.handlers[method] = handler
		}

		for key, ss := range sub.sub {
			mergeSubRouter(root, ss, key)
		}
		for key, st := range sub.trie.childs {
			mergeSubTrie(root.trie, st, key)
		}
		return r
	}

	segments := strings.Split(pattern, "/")
	segmentsLen := len(segments)
	for _, segment := range segments[1:segmentsLen-1] {
		if root.sub[segment] == nil {
			next := NewRouter()
			root.sub[segment] = next
			root.trie.childs[segment] = next.trie
		}
		root = root.sub[segment]
	}

	mergeSubRouter(root, sub, segments[segmentsLen-1])
	mergeSubTrie(root.trie, sub.trie, segments[segmentsLen-1])
	return r
}


// NotFound defines a handler to respond whenever a route could
// not be found.
func(r *Router) NotFound(handler HandlerFunc) {
	r.notFountHandler = handler
}


// MethodNotAllowed defines a handler to respond whenever a method is
// not allowed.
func(r *Router) MethodNotAllowed(handler HandlerFunc) {
	r.methodNotAllowedHandler = handler
}


// Handle adds routes for `pattern` by specify methods.
// If method isn't given, it will handler all methods.
func(r *Router) Handle(pattern string, fn HandlerFunc, methods ...string) {
	if len(methods) == 0 {
		r.trie.insert(pattern, fn, _ALL)
		return
	}
	r.trie.insert(pattern, fn, methods...)
}


func mergeSubRouter(root, sub *Router, pattern string) {
	if root.sub[pattern] != nil {
		root.sub[pattern].groupMiddlewares = append(root.sub[pattern].groupMiddlewares, sub.groupMiddlewares...)
		for key, r := range sub.sub {
			mergeSubRouter(root.sub[pattern], r, key)
		}
		return
	}

	root.sub[pattern] = sub
}


func mergeSubTrie(root, sub *trie, pattern string) {
	if next, found := root.childs[pattern]; found {
		for key, t := range sub.childs {
			mergeSubTrie(next, t, key)
		}
		return
	}

	root.childs[pattern] = sub
}
