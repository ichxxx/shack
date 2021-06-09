package rest

import (
	"strconv"
	"sync"

	"github.com/ichxxx/shack"
	"github.com/ichxxx/shack/util"
	jsoniter "github.com/json-iterator/go"
)

var (
	json     = jsoniter.ConfigCompatibleWithStandardLibrary
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
	return []byte(strconv.Quote(e.Error())), nil
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
	r.syncStatusCode(okC)
	if len(r.M) == 0 {
		r.M = okM
	}

	r.ctx.JSON(r)
}


func(r *resp) Fail() {
	r.syncStatusCode(failC)
	if len(r.M) == 0 {
		r.M = failM
	}

	if r.ctx.Err != nil && r.E == nil {
		r.Error(r.ctx.Err)
	} else if r.E != nil {
		r.ctx.Error(r.E)
	}

	r.ctx.JSON(r)
}


func(r *resp) syncStatusCode(code int) {
	if r.C == nil {
		if r.ctx != nil {
			if r.ctx.StatusCode != nil {
				r.C = r.ctx.StatusCode
			} else {
				r.C = &code
				r.ctx.Status(code)
			}
		} else {
			r.C = &code
		}
	}
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
		m[util.Str(keyAndValues[i-1])] = keyAndValues[i]
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