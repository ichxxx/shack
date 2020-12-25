package shack

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	param = ':'
	path  = '*'
	wild  = "*"
)


type trie struct {
	handlers   map[string][]HandlerFunc
	isParam    bool
	isPath     bool
	p          string
	child      map[string]*trie
	m          string
}


/*
func(t *trie) merge(other *trie) {
	if t.isParam != t.isParam || t.isPath != other.isPath || t.p != other.p {
		panic("shack: can't merge, trie do not have same root")
	}

	for method, handler := range other.handler {
		if tHandler, found := t.handler[method]; found {
			if reflect.ValueOf(handler).Pointer() != reflect.ValueOf(tHandler).Pointer() {
				panic(fmt.Sprintf("shack: method %s of the two trie conflicts", method))
			}
		} else {
			t.handler[method] = handler
		}
	}

	for key, trie := range other.child {
		if _, found := t.child[key]; found {
			panic(fmt.Sprintf("shack: child %s of the two trie conficts", key))
		} else {
			t.child[key] = trie
		}
	}
}
*/


func newTrie() *trie {
	return &trie{
		handlers: make(map[string][]HandlerFunc, 7),
		child: make(map[string]*trie),
	}
}


func isWild(segment string) bool {
	if len(segment) == 0 {
		return false
	}

	return segment[0] == param || segment[0] == path
}


func isVaildPattern(pattern string) (isVaild bool) {
	isVaild, _ = regexp.MatchString(`^\/[:*.\-\w]*(\/[:*.\-\w]+)*$`, pattern)
	return
}


func(t *trie) With(middleware ...HandlerFunc) {
	// insert from head
	t.handlers[t.m] = append(middleware, t.handlers[t.m]...)
}


func(t *trie) insert(method, pattern string, handler HandlerFunc) *trie {
	if !isVaildPattern(pattern) {
		panic("shack: pattern is not valid")
	}

	segments := strings.Split(pattern, "/")
	for i, segment := range segments {
		if segment == "" {
			continue
		}

		p := segment
		if isWild(segment) {
			segment = wild
		}

		if _, ok := t.child[segment]; !ok {
			t.child[segment] = newTrie()
		}

		t = t.child[segment]
		switch p[0] {
		case param:
			t.isParam = true
			t.p = p[1:]
		case path:
			t.isPath = true
			t.p = p[1:]
			if i != len(segments)-1 {
				panic("shack: *path can only use in the last")
			}
		}
	}

	if handler != nil {
		switch method {
		case ALL:
			if len(t.handlers) > 0 {
				panic("shack: can't route method 'ALL', method duplicated")
			}
		default:
			if t.handlers[ALL] != nil || t.handlers[method] != nil {
				panic(fmt.Sprintf("shack: can't route method '%s', method duplicated", method))
			}
		}
		t.m = method
		t.handlers[method] = append(t.handlers[method], handler)
	}

	return t
}


func(t *trie) search(method, pattern string) (handlers []HandlerFunc, params map[string]string, ok bool) {
	i := 1
	var splitLoc int
	for ; i < len(pattern); i++ {
		if pattern[i] == '/' {
			next := t.next(pattern[splitLoc+1:i])
			if next == nil {
				return
			}

			t = next
			if t.isPath {
				if params == nil {
					params = make(map[string]string)
				}
				params[t.p] = pattern[splitLoc:]
				break
			}

			if t.isParam {
				if params == nil {
					params = make(map[string]string)
				}
				params[t.p] = pattern[splitLoc+1:i]
			}

			splitLoc = i
		}
	}

	if i > 1 && !t.isPath {
		next := t.next(pattern[splitLoc+1:i])
		if next == nil {
			return
		}
		t = next
	}

	if t.isParam {
		if params == nil {
			params = make(map[string]string)
		}
		params[t.p] = pattern[splitLoc+1:i]
	}

	handlers = t.handlers[method]
	if handlers == nil {
		handlers = t.handlers[ALL]
	}
	ok = true
	return
}


func(t *trie) next(segment string) (next *trie) {
	next = t.child[segment]
	if next == nil {
		next = t.child[wild]
	}
	return
}


func(t *trie) print() {
	t.dfs(1)
}


func(t *trie) dfs(count int) {
	for key, trie := range t.child {
		fmt.Println(strings.Repeat("-", count), key)
		trie.dfs(count+1)
	}
}
