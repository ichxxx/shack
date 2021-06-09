package shack

import (
	"testing"
)


func TestLogger(t *testing.T) {
	r := NewRouter()
	r.GET("/logger", func(ctx *Context) {
		Log.Info("logger test")
		ctx.String("logger")
	})

	go Run(":8080", r)

	request(t, "127.0.0.1:8080", "GET", "/logger", nil)
	Logger("shack").Encoding("JSON").Output("./logs").Enable()
	request(t, "127.0.0.1:8080", "GET", "/logger", nil)
}
