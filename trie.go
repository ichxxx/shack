package shack

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ichxxx/shack/utils"
)

const (
	_PARAM = ':'
	_PATH  = '*'
	_WILD  = "*"
)

var (
	validPathReg, _    = regexp.Compile(`^\/[:*.\-\w]*(\/[:*.\-\w]+)*$`)
	validPatternReg, _ = regexp.Compile(`^\/[\-\w]*(\/[\-\w]+)*$`)
)

type trie struct {
	handlers map[string][]Handler
	isParam  bool
	isPath   bool
	childs   map[string]*trie
	p        string   // p means param or path
	m        []string // m means the passed methods
}

func newTrie() *trie {
	return &trie{
		handlers: make(map[string][]Handler, 7),
		childs:   make(map[string]*trie),
	}
}

func isWild(segment string) bool {
	if len(segment) == 0 {
		return false
	}
	return segment[0] == _PARAM || segment[0] == _PATH
}

func isValidPath(path string) bool {
	return validPathReg.MatchString(path)
}

func isValidPattern(pattern string) bool {
	return validPatternReg.MatchString(pattern)
}

// With adds one or more middlewares for an endpoint handler.
func (t *trie) With(middleware ...Handler) {
	// insert from head
	for _, method := range t.m {
		t.handlers[method] = append(middleware, t.handlers[method]...)
	}
}

func (t *trie) insert(path string, handler Handler, methods ...string) *trie {
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
				panic(fmt.Sprintf("shack: '*' can only use in the last in path '%s'", path))
			}
		}
	}

	if handler != nil {
		for _, method := range methods {
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
			t.m = methods
			t.handlers[method] = append(t.handlers[method], handler)
		}
	}

	return t
}

func (t *trie) search(method, path []byte) (handlers []Handler, params map[string]string, ok bool) {
	i := 1
	var splitPos int
	for ; i < len(path); i++ {
		if path[i] == '/' {
			next := t.next(utils.UnsafeString(path[splitPos+1 : i]))
			if next == nil {
				return
			}

			t = next
			if t.isPath {
				if params == nil {
					params = make(map[string]string)
				}
				params[t.p] = utils.UnsafeString(path[splitPos:])
				break
			}

			if t.isParam {
				if params == nil {
					params = make(map[string]string)
				}
				params[t.p] = utils.UnsafeString(path[splitPos+1 : i])
			}

			splitPos = i
		}
	}

	if i > 1 && !t.isPath {
		next := t.next(utils.UnsafeString(path[splitPos+1 : i]))
		if next == nil {
			return
		}
		t = next
	}

	if t.isParam {
		if params == nil {
			params = make(map[string]string)
		}
		params[t.p] = utils.UnsafeString(path[splitPos+1 : i])
	}

	handlers = t.handlers[utils.UnsafeString(method)]
	if handlers == nil {
		handlers = t.handlers[_ALL]
	}
	ok = true
	return
}

func (t *trie) next(segment string) (next *trie) {
	next = t.childs[segment]
	if next == nil {
		next = t.childs[_WILD]
	}
	return
}

func (t *trie) print() {
	t.dfsPrint(1)
}

func (t *trie) dfsPrint(count int) {
	for key, child := range t.childs {
		fmt.Println(strings.Repeat("-", count), key)
		child.dfsPrint(count + 1)
	}
}
