package rest

import (
	"fmt"
	"reflect"
	"strconv"
	"sync"

	"github.com/ichxxx/shack"
)

var (
	okC      = 0
	okM      = "success"
	failC    = 1
	failM    = "fail"
	respPool = &sync.Pool{New: func() interface{}{return new(resp)}}
)

type resp struct {
	ctx *shack.Context
	C   *int        `json:"code,omitempty"`
	M   string      `json:"msg,omitempty"`
	E   error       `json:"error,omitempty"`
	D   interface{} `json:"data,omitempty"`
}


type restErr struct {
	err error
}


func(e *restErr) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", e.Error())), nil
}


func(e *restErr) Error() string {
	return e.err.Error()
}


// R is a shortcut of Resp
func R(ctx *shack.Context) *resp {
	return Resp(ctx)
}


func Resp(ctx *shack.Context) *resp {
	r := respPool.Get().(*resp)
	r.reset()
	respPool.Put(r)
	r.ctx = ctx
	return r
}


func(r *resp) OK() {
	if (r.ctx == nil || r.ctx.StatusCode == nil) && r.C == nil {
		r.C = &okC
	} else if r.ctx != nil {
		if r.ctx.StatusCode != nil {
			r.C = r.ctx.StatusCode
		} else {
			r.ctx.Status(okC)
		}
	}

	if len(r.M) == 0 {
		r.M = okM
	}

	r.ctx.JSON(r)
}


func(r *resp) Fail() {
	if (r.ctx == nil || r.ctx.StatusCode == nil) && r.C == nil {
		r.C = &failC
	} else if r.ctx != nil {
		if r.ctx.StatusCode != nil {
			r.C = r.ctx.StatusCode
		} else {
			r.ctx.Status(failC)
		}
	}

	if len(r.M) == 0 {
		r.M = failM
	}

	if r.ctx.Err != nil && r.E == nil {
		r.Error(r.ctx.Err)
	} else if r.E != nil && r.ctx.Err == nil {
		r.ctx.Error(r.E)
	}

	r.ctx.JSON(r)
}


func(r *resp) Write() {
	r.ctx.JSON(r)
}


func(r *resp) Code(code int) *resp {
	r.C = &code
	return r
}


func(r *resp) Msg(msg string) *resp {
	r.M = msg
	return r
}


func(r *resp) Error(err error) *resp {
	r.E = &restErr{err: err}
	return r
}


func(r *resp) Data(keyAndValues ...interface{}) *resp {
	l := len(keyAndValues)
	if l <= 1 {
		if l == 1 {
			r.D = keyAndValues[0]
		}
		return r
	}

	m := make(map[string]interface{})
	for i := 1; i < l; i+=2 {
		m[str(keyAndValues[i-1])] = keyAndValues[i]
	}
	r.D = m
	return r
}


func DefaultOkCode(code int) {
	okC = code
}


func DefaultOkMsg(msg string) {
	okM = msg
}


func DefaultFailCode(code int) {
	failC = code
}


func DefaultFailMsg(msg string) {
	failM = msg
}


func(r *resp) reset() {
	r.ctx = nil
	r.C = nil
	r.M = r.M[0:0]
	r.E = nil
	r.D = nil
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
