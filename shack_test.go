package shack

import (
	"fmt"
	"io/ioutil"
	"testing"
)


func TestShackRouter(t *testing.T) {
	r := NewRouter()
	r.Handle("/test", func(ctx *Context) {
		b, _ := ioutil.ReadAll(ctx.Request.Body)
		fmt.Println(string(b))
	})
	ListenAndServe(":18080", r)
}
