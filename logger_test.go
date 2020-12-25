package shack

import (
	"net/http/httptest"
	"testing"
)


func TestLogger(t *testing.T) {
	r := NewRouter()
	r.GET("/logger", func(ctx *Context) {
		Log.Info("logger test")
		ctx.String("logger")
	})
	ts := httptest.NewServer(r)
	defer ts.Close()
	request(t, ts, "GET", "/logger", nil)
	Logger("shack").Encoding("JSON").Output("./logs").Enable()
	request(t, ts, "GET", "/logger", nil)
}
