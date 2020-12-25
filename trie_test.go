package shack

import (
	"fmt"
	"testing"
)


func TestTrie(t *testing.T) {
	handler1 := HandlerFunc(func(*Context){})
	handler2 := HandlerFunc(func(*Context){})
	handler3 := HandlerFunc(func(*Context){})
	inputs := []struct{
		method  string
		pattern string
		handler HandlerFunc
	}{
		{_GET,"/",             handler1},
		{_GET,"/foo/:var/bar", handler1},
		{_GET,"/foo/bar",      handler2},
		{_GET,"/bar/*path",    handler3},
		{_GET,"/*path",        handler3},
		{_ALL, "/foo/bar/all", handler2},
	}

	tests := []struct{
		method   string
		path     string
		ok       bool
		handler  HandlerFunc
		pKey     string
		pValue   string
	}{
		{_GET,    "/",             true,  handler1, "",    ""},
		{_GET,    "/foo/test/bar", true,  handler1, "var", "test"},
		{_GET,    "/foo/bar",      true,  handler2, "",    ""},
		{_GET,    "/bar/f/o/o",    true,  handler3, "path","/f/o/o"},
		{_GET,    "/f/o/bar.html", true,  handler3, "path","/f/o/bar.html"},
		{_GET,    "/foo/test",     true,  nil,      "var", "test"},
		{_GET,    "/foo/test/foo", false, nil,      "var", "test"},
		{_GET,    "/foo/bar/foo",  false, nil,      "",    ""},
		{_POST,   "/foo/test/bar", true,  nil,      "var", "test"},
		{_POST,   "/foo/bar/foo",  false, nil,      "",    ""},
		{_POST,   "/foo/bar/all",  true,  handler2, "",    ""},
		{_DELETE, "/foo/bar/all",  true,  handler2, "",    ""},
	}

	trie := newTrie()
	for _, input := range inputs {
		trie.insert(input.method, input.pattern, input.handler)
	}

	for i, test := range tests {
		handlers, param, ok := trie.search(test.method, test.path)
		if handler := popHandler(handlers); fmt.Sprintf("%v", handler) != fmt.Sprintf("%v", test.handler) {
			t.Errorf("input [%d]: expecting handler:%v, got:%v", i, test.handler, handler)
		}
		if param[test.pKey] != test.pValue {
			t.Errorf("input [%d]: expecting param %s:%s, got:%s", i, test.pKey, test.pValue, param[test.pKey])
		}
		if ok != test.ok {
			t.Errorf("input [%d]: expecting ok:%v, got:%v", i, test.ok, ok)
		}


	}
}


func popHandler(handlers []HandlerFunc) HandlerFunc {
	if len(handlers) > 0 {
		return handlers[0]
	}
	return nil
}
