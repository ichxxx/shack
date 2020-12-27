package rest

import (
	"fmt"
	"reflect"
	"strconv"
	"sync"
)

var (
	okC      = 0
	okM      = "success"
	failC    = 1
	failM    = "fail"
	respPool = &sync.Pool{New: func() interface{}{return new(resp)}}
)

type resp struct {
	C  int          `json:"code,omitempty"`
	M  string       `json:"msg,omitempty"`
	E  interface{}  `json:"error,omitempty"`
	D  map[string]interface{}
}


func Resp() *resp {
	r := respPool.Get().(*resp)
	r.reset()
	respPool.Put(r)
	return r
}


func(r *resp) OK() *resp {
	r.C = okC
	r.M = okM
	return r
}


func(r *resp) Fail() *resp {
	r.C = failC
	r.M = failM
	return r
}


func(r *resp) Code(code int) *resp {
	r.C = code
	return r
}


func(r *resp) Msg(msg string) *resp {
	r.M = msg
	return r
}


func(r *resp) Error(error interface{}) *resp {
	r.E = error
	return r
}


func(r *resp) Data(keyAndValues ...interface{}) *resp {
	l := len(keyAndValues)
	for i := 1; i < l; i+=2 {
		r.D[str(keyAndValues[i-1])] = keyAndValues[i]
	}

	return r
}


func(r *resp) DefaultOkCode(code int) {
	okC = code
}


func(r *resp) DefaultOkMsg(msg string) {
	okM = msg
}


func(r *resp) DefaultFailCode(code int) {
	failC = code
}


func(r *resp) DefaultFailMsg(msg string) {
	failM = msg
}


func(r *resp) reset() {
	r.C = 0
	r.M = r.M[0:0]
	r.E = nil
	r.D = make(map[string]interface{})
}


func str(i interface{}) string {
	switch i.(type) {
	case string:
		return i.(string)
	case int:
		return strconv.Itoa(i.(int))
	case uint:
		return strconv.FormatUint(uint64(i.(uint)), 10)
	case int64:
		return strconv.Itoa(int(i.(int64)))
	case uint64:
		return strconv.FormatUint(i.(uint64), 10)
	case int32:
		return strconv.Itoa(int(i.(int32)))
	case uint32:
		return strconv.FormatUint(uint64(i.(uint32)), 10)
	case int16:
		return strconv.Itoa(int(i.(int16)))
	case uint16:
		return strconv.FormatUint(uint64(i.(uint16)), 10)
	case int8:
		return strconv.Itoa(int(i.(int8)))
	case float64:
		return strconv.FormatFloat(i.(float64), 'E', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(i.(float32)), 'E', -1, 64)
	case byte:
		return string(i.(byte))
	case []byte:
		return string(i.([]byte))
	default:
		panic(fmt.Sprintf("shack: can't convert %v to string", reflect.TypeOf(i)))
	}
}
