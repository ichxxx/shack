package middleware

import (
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/ichxxx/shack"
)

func panicFunc() {
	panic("panic test")
}

func panicHandler(ctx *shack.Context) {
	panicFunc()
}

func TestRecovery(t *testing.T) {
	r := shack.NewRouter()
	r.GET("/panic", panicHandler).With(Recovery())
	go shack.Run(":8080", r)
	request(t, "http://127.0.0.1:8080", "GET", "/panic", nil)
}

func TestAccessLog(t *testing.T) {
	r := shack.NewRouter()
	r.GET("/access", func(ctx *shack.Context) {
		ctx.Response.String("access")
	}).With(AccessLog())
	go shack.Run(":8080", r)
	request(t, "http://127.0.0.1:8080", "GET", "/access", nil)
}

func request(t *testing.T, url, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, url+path, body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}
	defer resp.Body.Close()

	return resp, string(respBody)
}
