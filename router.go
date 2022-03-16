package shack

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/ichxxx/shack/utils"
)

var slashBytes = []byte("/")

type Router struct {
	sub                     map[string]*Router
	trie                    *trie
	middlewares             []Handler
	notFountHandler         Handler
	methodNotAllowedHandler Handler
}

func NewRouter() *Router {
	return &Router{
		sub:  make(map[string]*Router),
		trie: newTrie(),
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := getContext(req, w)
	c.handlers = append(c.handlers, getMiddlewares(r, utils.UnsafeBytes(c.Request.URI()))...)
	r.handler(c)
	releaseContext(c)
}

func getMiddlewares(r *Router, path []byte) (middlewares []Handler) {
	middlewares = append(middlewares, r.middlewares...)
	if len(path) <= 0 || bytes.Equal(path, slashBytes) {
		return
	}

	path = path[1:]

	for pattern, router := range r.sub {
		if patternBytes := utils.UnsafeBytes(pattern); bytes.HasPrefix(path, patternBytes) {
			nextPath := bytes.TrimPrefix(path, patternBytes)
			if len(nextPath) > 0 {
				middlewares = append(middlewares, getMiddlewares(router, nextPath)...)
			}
		}
	}
	return
}

func (r *Router) handler(ctx *Context) {
	handlers, params, ok := r.trie.search(utils.UnsafeBytes(ctx.Request.Method()), utils.UnsafeBytes(ctx.Request.URI()))
	if ok && len(handlers) > 0 {
		ctx.PathParams = params
		ctx.handlers = append(ctx.handlers, handlers...)
		ctx.Next()
	} else if ok {
		// todo: 如果是模糊节点，且没有handler，会通过判断，待修复
		ctx.Response.Status(http.StatusMethodNotAllowed)
		if r.methodNotAllowedHandler != nil {
			r.methodNotAllowedHandler(ctx)
		}
	} else {
		ctx.Response.Status(http.StatusNotFound)
		if r.notFountHandler != nil {
			r.notFountHandler(ctx)
		}
	}
}

func (r *Router) Handle(pattern string, handler Handler) *trie {
	return r.trie.insert(pattern, handler, _ALL)
}

func (r *Router) GET(pattern string, handler Handler) *trie {
	return r.trie.insert(pattern, handler, _GET)
}

func (r *Router) POST(pattern string, handler Handler) *trie {
	return r.trie.insert(pattern, handler, _POST)
}

func (r *Router) DELETE(pattern string, handler Handler) *trie {
	return r.trie.insert(pattern, handler, _DELETE)
}

func (r *Router) PUT(pattern string, handler Handler) *trie {
	return r.trie.insert(pattern, handler, _PUT)
}

func (r *Router) PATCH(pattern string, handler Handler) *trie {
	return r.trie.insert(pattern, handler, _PATCH)
}

func (r *Router) OPTIONS(pattern string, handler Handler) *trie {
	return r.trie.insert(pattern, handler, _OPTIONS)
}

func (r *Router) HEAD(pattern string, handler Handler) *trie {
	return r.trie.insert(pattern, handler, _HEAD)
}

// Use appends one or more middlewares onto the router.
func (r *Router) Use(middlewares ...Handler) {
	r.middlewares = append(r.middlewares, middlewares...)
}

// Mount attaches another router along a `pattern` string.
func (r *Router) Mount(pattern string, router *Router) {
	if !isValidPattern(pattern) {
		panic(fmt.Sprintf("shack: pattern '%s' is not valid", pattern))
	}

	if router == nil {
		panic(fmt.Sprintf("shack: router is nil while mounting '%s'", pattern))
	}

	if r.sub[pattern] != nil {
		panic(fmt.Sprintf("shack: pattern '%s' to mount is already exist", pattern))
	}

	// todo: should use mergeSubRouter and mergeSubTrie ?
	r.sub[pattern] = router
	r.trie.childs[pattern[1:]] = router.trie
}

// Group adds a sub-Router to the group along a `pattern` string.
func (r *Router) Group(pattern string, routeFunc func(r *Router)) *Router {
	if !isValidPattern(pattern) {
		panic(fmt.Sprintf("shack: pattern '%s' is not valid", pattern))
	}
	if routeFunc == nil {
		return r
	}

	root := r
	sub := NewRouter()
	routeFunc(sub)

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
	for _, segment := range segments[1 : segmentsLen-1] {
		if root.sub[segment] == nil {
			next := NewRouter()
			root.sub[segment] = next
			root.trie.childs[segment] = next.trie
		}
		root = root.sub[segment]
	}

	mergeSubRouter(root, sub, segments[segmentsLen-1])
	mergeSubTrie(root.trie, sub.trie, segments[segmentsLen-1])
	return root.sub[segments[segmentsLen-1]]
}

// Add is a shortcut of Group("/", fn).
func (r *Router) Add(routeFunc func(r *Router)) *Router {
	if routeFunc == nil {
		return r
	}

	root := r
	sub := NewRouter()
	routeFunc(sub)
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

// NotFound defines a handler to respond whenever a route could
// not be found.
func (r *Router) NotFound(handler Handler) {
	r.notFountHandler = handler
}

// MethodNotAllowed defines a handler to respond whenever a method is
// not allowed.
func (r *Router) MethodNotAllowed(handler Handler) {
	r.methodNotAllowedHandler = handler
}

func mergeSubRouter(root, sub *Router, pattern string) {
	if root.sub[pattern] != nil {
		root.sub[pattern].middlewares = append(root.sub[pattern].middlewares, sub.middlewares...)
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
