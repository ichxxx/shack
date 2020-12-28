package shack

import (
	"reflect"
	"testing"
)


func TestFormFlowBind(t *testing.T) {
	inputs := []map[string][]string{
		{"foo":{"123"}},
		{"123":{"foo"}},
		{"Foo":{"bar"}},
	}

	tests := []struct{
		dst interface{}
		result interface{}
	}{
		{&map[string]int{},         &map[string]int{"foo": 123}},
		{&map[int]string{},         &map[int]string{123: "foo"}},
		{&map[string]interface{}{}, &map[string]interface{}{"Foo":"bar"}},
	}

	for i, input := range inputs {
		f := newFormFlow(input)
		err := f.Bind(tests[i].dst)
		if err != nil {
			t.Errorf("input [%d]: got err:%v", i, err)
		} else if !reflect.DeepEqual(tests[i].dst, tests[i].result) {
			t.Errorf("input [%d]: expecting result:%v, got:%v", i, tests[i].result, tests[i].dst)
		}
	}
}


func TestFormFlowBindWithTag(t *testing.T) {
	inputs := []map[string][]string{
		{"foo":{"123"}, "Foo":{"foo"}},
		{"Baz":{"abc"}, "bar":{"123"}},
	}

	type tmpStruct struct {
		Foo string
		Bar string
		Baz int     `form:"foo"`
		Qux float64 `json:"bar"`
	}

	tests := []struct{
		dst interface{}
		tag string
		result interface{}
	}{
		{&tmpStruct{},"form",&tmpStruct{Foo: "foo", Baz: 123}},
		{&tmpStruct{},"json",&tmpStruct{Baz: 0, Qux: 123}},
	}

	for i, input := range inputs {
		f := newFormFlow(input)
		err := f.Bind(tests[i].dst, tests[i].tag)
		if err != nil {
			t.Errorf("input [%d]: got err:%v", i, err)
		} else if !reflect.DeepEqual(tests[i].dst, tests[i].result) {
			t.Errorf("input [%d]: expecting result:%v, got:%v", i, tests[i].result, tests[i].dst)
		}
	}
}
