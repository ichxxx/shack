package shack

import (
	"fmt"
	"testing"
)


func TestTrie(t *testing.T) {
	handler1 := HandlerFunc(func(*Context){})
	handler2 := HandlerFunc(func(*Context){})
	inputs := []struct{
		method  methodType
		pattern string
		handler HandlerFunc
	}{
		{GET,"/foo/:var/bar", handler1},
		{GET,"/foo/bar",      handler2},
	}

	tests := []struct{
		method  methodType
		path    string
		ok      bool
		handler HandlerFunc
		pKey    string
		pValue  string
	}{
		{GET, "/foo/test/bar", true,  handler1, "var", "test"},
		{GET, "/foo/bar",      true,  handler2, "",    ""},
		{GET, "/foo/test",     false, nil,      "var", "test"},
		{GET, "/foo/test/foo", false, nil,      "var", "test"},
		{GET, "/foo/bar/foo",  false, nil,      "",    ""},
		{POST,"/foo/test/bar", true,  nil,      "var", "test"},
		{POST,"/foo/bar/foo",  false, nil,      "",    ""},
	}

	trie := newTrie()
	for _, input := range inputs {
		trie.insert(input.method, input.pattern, input.handler)
	}

	for i, test := range tests {
		handler, param, ok := trie.search(test.method, test.path)
		if fmt.Sprintf("%v", handler) != fmt.Sprintf("%v", test.handler) {
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
