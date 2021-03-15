package util

import (
	"strconv"
	"strings"
)


func Str(i interface{}, _default ...string) (s string) {
	switch i.(type) {
	case string:
		s = i.(string)
	case int:
		s = strconv.Itoa(i.(int))
	case int8:
		s = strconv.Itoa(int(i.(int8)))
	case int16:
		s = strconv.Itoa(int(i.(int16)))
	case int32:
		s = strconv.Itoa(int(i.(int32)))
	case int64:
		s = strconv.Itoa(int(i.(int64)))
	case uint:
		s = strconv.Itoa(int(i.(uint)))
	case uint8:
		s = strconv.Itoa(int(i.(uint8)))
	case uint16:
		s = strconv.Itoa(int(i.(uint16)))
	case uint32:
		s = strconv.Itoa(int(i.(uint32)))
	case uint64:
		s = strconv.Itoa(int(i.(uint64)))
	case float32:
		s = strconv.FormatFloat(float64(i.(float32)), 'f', -1, 32)
	case float64:
		s = strconv.FormatFloat(i.(float64), 'f', -1, 64)
	case bool:
		s = strconv.FormatBool(i.(bool))
	}

	if len(s) == 0 && len(_default) > 0 {
		s = strings.Join(_default, "")
	}

	return
}


func Int(i interface{}, _default ...int) int {
	var n *int
	switch i.(type) {
	case string:
		tmp, _ := strconv.Atoi(i.(string))
		n = &tmp
	case int:
		tmp := i.(int)
		n = &tmp
	case int8:
		tmp := int(i.(int8))
		n = &tmp
	case int16:
		tmp := int(i.(int16))
		n = &tmp
	case int32:
		tmp := int(i.(int32))
		n = &tmp
	case int64:
		tmp := int(i.(int64))
		n = &tmp
	case uint:
		tmp := int(i.(uint))
		n = &tmp
	case uint8:
		tmp := int(i.(uint8))
		n = &tmp
	case uint16:
		tmp := int(i.(uint16))
		n = &tmp
	case uint32:
		tmp := int(i.(uint32))
		n = &tmp
	case uint64:
		tmp := int(i.(uint64))
		n = &tmp
	case float32:
		tmp := int(i.(float32))
		n = &tmp
	case float64:
		tmp := int(i.(float64))
		n = &tmp
	}

	if n == nil {
		if len(_default) == 0 {
			return 0
		}
		return _default[0]
	}

	return *n
}
