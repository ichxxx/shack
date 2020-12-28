package middleware

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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
	ts := httptest.NewServer(r)
	defer ts.Close()
	request(t, ts, "GET", "/panic", nil)
}


func TestAccessLog(t *testing.T) {
	r := shack.NewRouter()
	r.GET("/access", func(ctx *shack.Context) {
		ctx.String("access")
	}).With(AccessLog())
	ts := httptest.NewServer(r)
	defer ts.Close()
	request(t, ts, "GET", "/access", nil)
}


func request(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
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
