package shack

import (
	"testing"
)


func TestShackRouter(t *testing.T) {
	r := NewRouter()
	r.GET("/test", func(ctx *Context) {
		ctx.JSON(map[string]string{"foo":"bar"})
	})
	ListenAndServe(":18080", r)
}
