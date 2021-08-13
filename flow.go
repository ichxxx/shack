package shack

import (
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
	return Json.Unmarshal([]byte(f), dst)
}


// BindJson binds the passed struct pointer with the raw value parsed to json.
func(f bodyFlow) BindJson(dst interface{}) error {
	return Json.Unmarshal(f, dst)
}


// Bind binds the passed struct pointer with the raw value parsed by the given tag.
// If the tag isn't given, it will parse according to key's name.
func(f rawFlow) Bind(dst interface{}, tag ...string) error {
	p := reflect.ValueOf(dst)
	if p.Kind() != reflect.Ptr || p.IsNil() {
		return errors.New("shack: raw flow bind error, dst must be a pointer")
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
		return errors.New("shack: form flow bind error, dst must be a pointer")
	}

	m := map[string]string{}
	for k, v := range f {
		m[k] = v[0]
	}

	return mapTo(reflect.Indirect(p), m, tag...)
}


func mapTo(rv reflect.Value, m map[string]string, tag ...string) error {
	if rv.Kind() != reflect.Struct && rv.IsNil() {
		return errors.New("shack: map value error, dst is nil")
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
				return errors.New("shack: map error, dst struct doesn't have any fields")
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
					rvField := rv.Field(i)
					if rvField.Kind() == reflect.Ptr {
						rvField.Set(toValuePtr(v, rvField.Type().Elem().Kind()))
					} else {
						rvField.Set(toValue(v, rvField.Kind()))
					}
				}
			}
		}
	}

	return nil
}


func toValue(src string, dstKind reflect.Kind) reflect.Value {
	return reflect.ValueOf(srcTrans(src, dstKind))
}


func srcTrans(src string, dstKind reflect.Kind) interface{} {
	switch dstKind {
	case reflect.Bool:
		b, _ := strconv.ParseBool(src)
		return b
	case reflect.Int:
		i, _ := strconv.Atoi(src)
		return i
	case reflect.Int8:
		i, _ := strconv.Atoi(src)
		return int8(i)
	case reflect.Int16:
		i, _ := strconv.Atoi(src)
		return int16(i)
	case reflect.Int32:
		i, _ := strconv.Atoi(src)
		return int32(i)
	case reflect.Int64:
		i, _ := strconv.Atoi(src)
		return int64(i)
	case reflect.Uint:
		i, _ := strconv.Atoi(src)
		return uint(i)
	case reflect.Uint8:
		i, _ := strconv.Atoi(src)
		return uint8(i)
	case reflect.Uint16:
		i, _ := strconv.Atoi(src)
		return uint16(i)
	case reflect.Uint32:
		i, _ := strconv.Atoi(src)
		return uint32(i)
	case reflect.Uint64:
		i, _ := strconv.Atoi(src)
		return uint64(i)
	case reflect.Float32:
		f, _ := strconv.ParseFloat(src, 32)
		return float32(f)
	case reflect.Float64:
		f, _ := strconv.ParseFloat(src, 64)
		return f
	case reflect.Interface:
		return src
	}

	return src
}

func toValuePtr(src string, dstKind reflect.Kind) reflect.Value {
	switch dstKind {
	case reflect.Bool:
		b, _ := strconv.ParseBool(src)
		return reflect.ValueOf(&b)
	case reflect.Int:
		i, _ := strconv.Atoi(src)
		return reflect.ValueOf(&i)
	case reflect.Int8:
		i, _ := strconv.Atoi(src)
		i8 := int8(i)
		return reflect.ValueOf(&i8)
	case reflect.Int16:
		i, _ := strconv.Atoi(src)
		i16 := int16(i)
		return reflect.ValueOf(&i16)
	case reflect.Int32:
		i, _ := strconv.Atoi(src)
		i32 := int8(i)
		return reflect.ValueOf(&i32)
	case reflect.Int64:
		i, _ := strconv.Atoi(src)
		i64 := int64(i)
		return reflect.ValueOf(&i64)
	case reflect.Uint:
		i, _ := strconv.Atoi(src)
		ui := uint(i)
		return reflect.ValueOf(&ui)
	case reflect.Uint8:
		i, _ := strconv.Atoi(src)
		ui8 := uint8(i)
		return reflect.ValueOf(&ui8)
	case reflect.Uint16:
		i, _ := strconv.Atoi(src)
		ui16 := uint16(i)
		return reflect.ValueOf(&ui16)
	case reflect.Uint32:
		i, _ := strconv.Atoi(src)
		ui32 := uint32(i)
		return reflect.ValueOf(&ui32)
	case reflect.Uint64:
		i, _ := strconv.Atoi(src)
		ui64 := uint64(i)
		return reflect.ValueOf(&ui64)
	case reflect.Float32:
		f, _ := strconv.ParseFloat(src, 32)
		f32 := float32(f)
		return reflect.ValueOf(&f32)
	case reflect.Float64:
		f, _ := strconv.ParseFloat(src, 64)
		f64 := float64(f)
		return reflect.ValueOf(&f64)
	case reflect.Interface:
		var i interface{} = src
		return reflect.ValueOf(&i)
	}

	return reflect.ValueOf(&src)
}