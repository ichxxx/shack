package shack

import (
	"encoding/json"
	"errors"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

type (
	rawFlow   string
	valueFlow string
	bodyFlow  []byte
	formFlow  map[string][]string
)


func newRawFlow(value string) rawFlow {
	value, _ = url.QueryUnescape(value)
	return rawFlow(value)
}


func newValueFlow(value string) valueFlow {
	return valueFlow(value)
}


func newBodyFlow(value []byte) bodyFlow {
	return value
}


func newFormFlow(value map[string][]string) formFlow {
	return value
}


// Value returns the raw value of the workflow.
func(f rawFlow) Value() string {
	return string(f)
}


// Value returns the raw value of the workflow.
func(f valueFlow) Value() string {
	return string(f)
}


// Int trans the raw value to int.
func(f valueFlow) Int() int {
	i, _ := strconv.Atoi(f.Value())
	return i
}


// Int64 trans the raw value to int64.
func(f valueFlow) Int64() int64 {
	return int64(f.Int())
}


// Int8 trans the raw value to int8.
func(f valueFlow) Int8() int8 {
	return int8(f.Int())
}


// Float64 trans the raw value to float64.
func(f valueFlow) Float64() float64 {
	f64, _ := strconv.ParseFloat(f.Value(), 64)
	return f64
}


// Bool trans the raw value to bool.
func(f valueFlow) Bool() bool {
	b, _ := strconv.ParseBool(f.Value())
	return b
}


// Value returns the raw value of the workflow.
func(f bodyFlow) Value() []byte {
	return f
}


// Value returns the raw value of the workflow.
func(f formFlow) Value() map[string][]string {
	return f
}


// BindJson binds the passed struct pointer with the raw value parsed to json.
func(f valueFlow) BindJson(dst interface{}) error {
	return json.Unmarshal([]byte(f), dst)
}


// BindJson binds the passed struct pointer with the raw value parsed to json.
func(f bodyFlow) BindJson(dst interface{}) error {
	return json.Unmarshal(f, dst)
}


// Bind binds the passed struct pointer with the raw value parsed by the given tag.
// If the tag isn't given, it will parse according to key's name.
func(f rawFlow) Bind(dst interface{}, tag ...string) error {
	p := reflect.ValueOf(dst)
	if p.Kind() != reflect.Ptr || p.IsNil() {
		return errors.New("dst must be a pointer")
	}

	m := make(map[string]string)
	segments := strings.Split(f.Value(), "&")
	for _, segment := range segments {
		kv := strings.Split(segment, "=")
		if len(kv) > 1 {
			m[kv[0]] = kv[1]
		}
	}

	return mapTo(p.Elem(), m, tag...)
}


// Bind binds the passed struct pointer with the raw value parsed by the given tag.
// If the tag isn't given, it will parse according to key's name.
func(f formFlow) Bind(dst interface{}, tag ...string) error {
	p := reflect.ValueOf(dst)
	if p.Kind() != reflect.Ptr || p.IsNil() {
		return errors.New("dst is not a pointer")
	}

	m := map[string]string{}
	for k, v := range f {
		m[k] = v[0]
	}

	return mapTo(reflect.Indirect(p), m, tag...)
}


func(f rawFlow) reset() {}


func(f valueFlow) reset() {}


func(f bodyFlow) reset() {}


func(f formFlow) reset() {}


func mapTo(rv reflect.Value, m map[string]string, tag ...string) error {
	if rv.Kind() != reflect.Struct && rv.IsNil() {
		return errors.New("dst is nil")
	}

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

			fieldLoop:
			for i := 0; i < size; i++ {
				field := t.Field(i)
				if field.Type.Kind() == reflect.Struct {
					err := mapTo(reflect.ValueOf(field.Type), m, tag...)
					if err != nil {
						return err
					}
				}

				var _tag string
				if len(tag) > 0 {
					_tag = tag[0]
				}
				key := field.Tag.Get(_tag)
				if len(key) == 0 {
					if field.Name == k {
						rv.Field(i).Set(toValue(v, rv.Field(i).Kind()))
					}
					continue fieldLoop
				}

				if key == "-" {
					continue fieldLoop
				}

				key = strings.TrimSuffix(key, ",omitempty")
				if key == k {
					rv.Field(i).Set(toValue(v, rv.Field(i).Kind()))
				}
			}
		}
	}

	return nil
}


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
