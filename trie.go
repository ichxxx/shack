package shack

import (
	"strings"
)

const wild = "*"


type trie struct {
	handler    map[string]HandlerFunc
	isWild     bool
	hasWild    bool
	wildParam  string
	next       map[string]*trie
}


func newTrie() *trie {
	return &trie{
		handler: make(map[string]HandlerFunc, 7),
		next: make(map[string]*trie),
	}
}


func(t *trie) judgeWild(segment string) (isWild bool) {
	if len(segment) == 0 {
		panic("shack: parse pattern error")
	}

	if segment[0] == ':' || segment[0] == '*' {
		t.wildParam = segment[1:]
		t.hasWild = true
		isWild = true
	}
	return
}


func(t *trie) insert(method, pattern string, handler HandlerFunc) {
	segments := strings.Split(pattern, "/")[1:]
	for _, segment := range segments {
		if _, ok := t.next[segment]; !ok {
			if isWild := t.judgeWild(segment); isWild {
				segment = wild
			}

			t.next[segment] = newTrie()
		}

		t = t.next[segment]
		if segment == wild {
			t.isWild = true
		}
	}

	if handler != nil {
		if t.handler[method] != nil {
			panic("shack: pattern duplicated")
		}
		t.handler[method] = handler
	}

	return
}


func(t *trie) search(method, pattern string) (handler HandlerFunc, paramMap map[string]string, ok bool) {
	segments := strings.Split(pattern, "/")[1:]
	paramMap = make(map[string]string)
	for _, segment := range segments {
		if _, _ok := t.next[segment]; !_ok {
			if t.hasWild {
				paramMap[t.wildParam] = segment
				segment = wild
			} else {
				return
			}
		}

		t = t.next[segment]
	}

	handler = t.handler[method]
	if !t.isWild || handler != nil {
		ok = true
	}
	return
}
