package shack

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	_PARAM = ':'
	_PATH  = '*'
	_WILD  = "*"
)


type trie struct {
	handlers   map[string][]HandlerFunc
	isParam    bool
	isPath     bool
	p          string
	childs     map[string]*trie
	m          string // m means the passed method
}


func newTrie() *trie {
	return &trie{
		handlers: make(map[string][]HandlerFunc, 7),
		childs: make(map[string]*trie),
	}
}


func isWild(segment string) bool {
	if len(segment) == 0 {
		return false
	}

	return segment[0] == _PARAM || segment[0] == _PATH
}


func isValidPath(path string) (isValid bool) {
	isValid, _ = regexp.MatchString(`^\/[:*.\-\w]*(\/[:*.\-\w]+)*$`, path)
	return
}


func isValidPattern(pattern string) (isValid bool) {
	isValid, _ = regexp.MatchString(`^\/[\-\w]*(\/[\-\w]+)*$`, pattern)
	return
}


// With adds one or more middlewares for an endpoint handler.
func(t *trie) With(middleware ...HandlerFunc) {
	// insert from head
	t.handlers[t.m] = append(middleware, t.handlers[t.m]...)
}


func(t *trie) insert(method, path string, handler HandlerFunc) *trie {
	if !isValidPath(path) {
		panic(fmt.Sprintf("shack: path '%s' is not valid", path))
	}

	segments := strings.Split(path, "/")
	for i, segment := range segments {
		if segment == "" {
			continue
		}

		p := segment
		if isWild(segment) {
			segment = _WILD
		}

		if _, ok := t.childs[segment]; !ok {
			t.childs[segment] = newTrie()
		}

		t = t.childs[segment]
		switch p[0] {
		case _PARAM:
			t.isParam = true
			t.p = p[1:]
		case _PATH:
			t.isPath = true
			t.p = p[1:]
			if i != len(segments)-1 {
				panic(fmt.Sprintf("shack: '%s' *path can only use in the last", path))
			}
		}
	}

	if handler != nil {
		switch method {
		case _ALL:
			if len(t.handlers) > 0 {
				panic("shack: can't route method 'ALL', method duplicated")
			}
		default:
			if t.handlers[_ALL] != nil || t.handlers[method] != nil {
				panic(fmt.Sprintf("shack: can't route method '%s', method duplicated", method))
			}
		}
		t.m = method
		t.handlers[method] = append(t.handlers[method], handler)
	}

	return t
}


func(t *trie) search(method, path string) (handlers []HandlerFunc, params map[string]string, ok bool) {
	i := 1
	var splitLoc int
	for ; i < len(path); i++ {
		if path[i] == '/' {
			next := t.next(path[splitLoc+1:i])
			if next == nil {
				return
			}

			t = next
			if t.isPath {
				if params == nil {
					params = make(map[string]string)
				}
				params[t.p] = path[splitLoc:]
				break
			}

			if t.isParam {
				if params == nil {
					params = make(map[string]string)
				}
				params[t.p] = path[splitLoc+1:i]
			}

			splitLoc = i
		}
	}

	if i > 1 && !t.isPath {
		next := t.next(path[splitLoc+1:i])
		if next == nil {
			return
		}
		t = next
	}

	if t.isParam {
		if params == nil {
			params = make(map[string]string)
		}
		params[t.p] = path[splitLoc+1:i]
	}

	handlers = t.handlers[method]
	if handlers == nil {
		handlers = t.handlers[_ALL]
	}
	ok = true
	return
}


func(t *trie) next(segment string) (next *trie) {
	next = t.childs[segment]
	if next == nil {
		next = t.childs[_WILD]
	}
	return
}


func(t *trie) print() {
	t.dfs(1)
}


func(t *trie) dfs(count int) {
	for key, child := range t.childs {
		fmt.Println(strings.Repeat("-", count), key)
		child.dfs(count+1)
	}
}
