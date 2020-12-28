package shack

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

var (
	vFlowPool = &sync.Pool{New: func() interface{}{return new(valueFlow)}}
	bFlowPool = &sync.Pool{New: func() interface{}{return new(bodyFlow)}}
	fFlowPool = &sync.Pool{New: func() interface{}{return new(formFlow)}}
)

type valueFlow struct {
	value string
}

type bodyFlow struct {
	value []byte
}

type formFlow struct {
	value map[string][]string
}


func newValueFlow(value string) *valueFlow {
	f := vFlowPool.Get().(*valueFlow)
	vFlowPool.Put(f)
	f.reset()
	f.value = value
	return f
}


func newBodyFlow(value []byte) *bodyFlow {
	f := bFlowPool.Get().(*bodyFlow)
	bFlowPool.Put(f)
	f.value = value
	return f
}


func newFormFlow(value map[string][]string) *formFlow {
	f := fFlowPool.Get().(*formFlow)
	fFlowPool.Put(f)
	f.value = value
	return f
}


// Value returns the raw value of the workflow.
func(f *valueFlow) Value() string {
	return f.value
}


// Int trans the raw value to int.
func(f *valueFlow) Int() int {
	i, _ := strconv.Atoi(f.value)
	return i
}


// Float64 trans the raw value to float64.
func(f *valueFlow) Float64() float64 {
	f64, _ := strconv.ParseFloat(f.value, 64)
	return f64
}


// Bool trans the raw value to bool.
func(f *valueFlow) Bool() bool {
	b, _ := strconv.ParseBool(f.value)
	return b
}


// Value returns the raw value of the workflow.
func(f *bodyFlow) Value() []byte {
	return f.value
}


// Value returns the raw value of the workflow.
func(f *formFlow) Value() map[string][]string {
	return f.value
}


// BindJson binds the passed struct pointer with the raw value parsed to json.
func(f *valueFlow) BindJson(dst interface{}) error {
	return json.Unmarshal([]byte(f.value), dst)
}


// BindJson binds the passed struct pointer with the raw value parsed to json.
func(f *bodyFlow) BindJson(dst interface{}) error {
	return json.Unmarshal(f.value, dst)
}


// Bind binds the passed struct pointer with the raw value parsed by the given tag.
// If the tag isn't given, it will parse according to key's name.
func(f *valueFlow) Bind(dst interface{}, tag ...string) error {
	p := reflect.ValueOf(dst)
	if p.Kind() != reflect.Ptr || p.IsNil() {
		return errors.New("dst must be a pointer")
	}

	m := make(map[string]string)
	segments := strings.Split(f.value, "&")
	for _, segment := range segments {
		kv := strings.Split(segment, "=")
		if len(kv) > 1 {
			m[kv[0]] = kv[1]
		}
	}

	rv := p.Elem()
	switch rv.Kind() {
	case reflect.Map:
		kType := rv.Type().Key().Kind()
		vType := rv.Type().Elem().Kind()
		for k, v := range m {
			rv.SetMapIndex(toValue(k, kType), toValue(v, vType))
		}
	case reflect.Struct:
		for k, v := range m {
			t := rv.Type()
			size := rv.NumField()
			if size == 0 {
				return errors.New("dst struct doesn't have any fields")
			}

			for i := 0; i < size; i++ {
				var _tag string
				if len(tag) > 0 {
					_tag = tag[0]
				}
				key := t.Field(i).Tag.Get(_tag)
				if len(key) == 0 {
					key = k
				}
				vType := rv.FieldByName(key).Kind()
				rv.FieldByName(key).Set(toValue(v, vType))
			}
		}
	}
	return nil
}


// Bind binds the passed struct pointer with the raw value parsed by the given tag.
// If the tag isn't given, it will parse according to key's name.
func(f *formFlow) Bind(dst interface{}, tag ...string) error {
	p := reflect.ValueOf(dst)
	if p.Kind() != reflect.Ptr || p.IsNil() {
		return errors.New("dst is not a pointer")
	}

	rv := reflect.Indirect(p)
	if rv.Kind() != reflect.Struct && rv.IsNil() {
		return errors.New("dst is nil")
	}

	switch rv.Kind() {
	case reflect.Map:
		kType := rv.Type().Key().Kind()
		vType := rv.Type().Elem().Kind()
		for k, v := range f.value {
			rv.SetMapIndex(toValue(k, kType), toValue(v[0], vType))
		}
	case reflect.Struct:
		for k, v := range f.value {
			t := rv.Type()
			size := rv.NumField()
			if size == 0 {
				return errors.New("dst struct doesn't have any fields")
			}

			for i := 0; i < size; i++ {
				field := t.Field(i)
				var _tag string
				if len(tag) > 0 {
					_tag = tag[0]
				}
				key := field.Tag.Get(_tag)
				if len(key) == 0 {
					key = field.Name
				}

				if key == k {
					rv.Field(i).Set(toValue(v[0], rv.Field(i).Kind()))
				}
			}
		}
	}
	return nil
}


func(f *valueFlow) reset() {}


func(f *bodyFlow) reset() {}


func(f *formFlow) reset() {}


func toValue(src string, dType reflect.Kind) reflect.Value {
	switch dType {
	case reflect.Bool:
		b, _ := strconv.ParseBool(src)
		return reflect.ValueOf(b)
	case reflect.Int:
		i, _ := strconv.Atoi(src)
		return reflect.ValueOf(i)
	case reflect.Int8:
		i, _ := strconv.Atoi(src)
		return reflect.ValueOf(int8(i))
	case reflect.Int16:
		i, _ := strconv.Atoi(src)
		return reflect.ValueOf(int16(i))
	case reflect.Int32:
		i, _ := strconv.Atoi(src)
		return reflect.ValueOf(int32(i))
	case reflect.Int64:
		i, _ := strconv.Atoi(src)
		return reflect.ValueOf(int64(i))
	case reflect.Uint:
		i, _ := strconv.Atoi(src)
		return reflect.ValueOf(uint(i))
	case reflect.Uint8:
		i, _ := strconv.Atoi(src)
		return reflect.ValueOf(uint8(i))
	case reflect.Uint16:
		i, _ := strconv.Atoi(src)
		return reflect.ValueOf(uint16(i))
	case reflect.Uint32:
		i, _ := strconv.Atoi(src)
		return reflect.ValueOf(uint32(i))
	case reflect.Uint64:
		i, _ := strconv.Atoi(src)
		return reflect.ValueOf(uint64(i))
	case reflect.Float32:
		f, _ := strconv.ParseFloat(src, 32)
		return reflect.ValueOf(float32(f))
	case reflect.Float64:
		f, _ := strconv.ParseFloat(src, 64)
		return reflect.ValueOf(f)
	case reflect.Interface:
		var i interface{} = src
		return reflect.ValueOf(i)

	}
	return reflect.ValueOf(src)
}
