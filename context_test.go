package shack

import (
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

	go Run(":8080", r)

	if _, body := request(t, "127.0.0.1:8080", "GET", "/abort", nil); body != "abort" {
		t.Fatalf(body)
	}
}
