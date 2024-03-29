package shack

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

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
			if len(v) == 0 {
				continue
			}

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
					key = field.Name
				}

				if key == "-" {
					continue fieldLoop
				}

				key = strings.Split(key, ",")[0]
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
