package shack

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

func TestName(t *testing.T) {
	var a int16 = 1
	ctx1 := Context{}
	fmt.Printf("part1 size: %d, align: %d\n", unsafe.Sizeof(ctx1), unsafe.Alignof(ctx1))
	fmt.Printf("*int16 size: %d\n", unsafe.Sizeof(&a))
	fmt.Printf("int16 size: %d\n", unsafe.Sizeof(int16(1)))
	fmt.Printf("int8 size: %d\n", unsafe.Sizeof(int8(1)))
	fmt.Printf("[]HandlerFunc size: %d\n", unsafe.Sizeof([]HandlerFunc{}))
	fmt.Printf("[]byte size: %d\n", unsafe.Sizeof([]byte{}))
	fmt.Printf("error size: %d\n", unsafe.Sizeof(errors.New("1")))
}


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
		{"foo":{"123"}, "Foo":{"foo"}, "Nil": {"nil"}},
		{"Baz":{"abc"}, "bar":{"123"}, "Nil": {"nil"}},
	}

	type tmpStruct struct {
		Foo string
		Bar string
		Baz int     `form:"foo"`
		Qux float64 `json:"bar"`
		Nil string  `json:"-" form:"-"`
	}

	tests := []struct{
		dst interface{}
		tag string
		result interface{}
	}{
		{&tmpStruct{},"form",&tmpStruct{Foo: "foo", Bar: "", Baz: 123, Qux: 0, Nil: ""}},
		{&tmpStruct{},"json",&tmpStruct{Foo: "", Bar: "", Baz: 0, Qux: 123, Nil: ""}},
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
