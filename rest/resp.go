package rest

import (
	"strconv"
	"sync"

	"github.com/ichxxx/shack"
	"github.com/spf13/cast"
)

var (
	okC      = 0
	okM      = "success"
	failC    = 1
	failM    = "fail"
	respPool = &sync.Pool{New: func() interface{} { return new(resp) }}
)

type resp struct {
	ctx *shack.Context
	C   *int        `json:"status,omitempty"`
	M   string      `json:"msg,omitempty"`
	E   error       `json:"error,omitempty"`
	D   interface{} `json:"data,omitempty"`
}

type restErr struct {
	err error
}

func (e *restErr) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(e.Error())), nil
}

func (e *restErr) Error() string {
	if e.err == nil {
		return ""
	}
	return e.err.Error()
}

func Resp(ctx *shack.Context) *resp {
	r := respPool.Get().(*resp)
	r.ctx = ctx
	return r
}

func (r *resp) OK() error {
	defer func() {
		r.reset()
		respPool.Put(r)
	}()

	r.C = &okC
	if len(r.M) == 0 {
		r.M = okM
	}
	return r.ctx.Response.JSON(r)
}

func (r *resp) Fail() error {
	defer func() {
		r.reset()
		respPool.Put(r)
	}()

	r.C = &failC
	if len(r.M) == 0 {
		r.M = failM
	}
	return r.ctx.Response.JSON(r)
}

func (r *resp) Code(code int) *resp {
	r.C = &code
	return r
}

func (r *resp) Msg(msg string) *resp {
	r.M = msg
	return r
}

func (r *resp) Error(err error) *resp {
	r.E = &restErr{err: err}
	r.ctx.Error(err)
	return r
}

func (r *resp) Data(keyAndValues ...interface{}) *resp {
	dataLen := len(keyAndValues)
	if dataLen > 1 {
		kv := make(map[string]interface{}, dataLen>>1)
		for i := 1; i < dataLen; i += 2 {
			kv[cast.ToString(keyAndValues[i-1])] = keyAndValues[i]
		}
		r.D = kv
	} else if dataLen == 1 {
		r.D = keyAndValues[0]
	}
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

func (r *resp) reset() {
	r.ctx = nil
	r.C = nil
	r.M = r.M[0:0]
	r.E = nil
	r.D = nil
}
