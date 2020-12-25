package shack

import (
	"net/http/httptest"
	"testing"
)


func abortHandler(ctx *Context) {
	ctx.String("abort")
	ctx.Abort()
}


func TestContextAbort(t *testing.T) {
	r := NewRouter()
	r.GET("/abort", func(ctx *Context) {
		ctx.String("hello")
	}).With(abortHandler)

	ts := httptest.NewServer(r)
	defer ts.Close()
	if _, body := request(t, ts, "GET", "/abort", nil); body != "abort" {
		t.Fatalf(body)
	}
}
