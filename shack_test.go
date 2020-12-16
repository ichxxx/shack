package shack

import (
	"testing"
)


func TestShackRouter(t *testing.T) {
	e := New()
	r := NewRouter()
	r.GET("/test", func(ctx *Context) {
		ctx.JSON(map[string]string{"foo":"bar"})
	})
	e.ListenAndServe(":18080", r)
}
