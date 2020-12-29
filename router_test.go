package shack

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)


func forAll() HandlerFunc {
	return func(c *Context) {
		t := time.Now()
		c.String("for all")
		log.Printf("%s in %v for all", c.Path, time.Since(t))
	}
}


func onlyForV1() HandlerFunc {
	return func(c *Context) {
		t := time.Now()
		c.String(" and ")
		c.String("only for v1")
		log.Printf("%s in %v for v1", c.Path, time.Since(t))
	}
}



func onlyForV2() HandlerFunc {
	return func(c *Context) {
		t := time.Now()
		c.String(" and ")
		c.String("only for v2")
		log.Printf("%s in %v for group v2", c.Path, time.Since(t))
	}
}


func TestRouterGroup(t *testing.T) {
	r2 := func(r *Router) {
		r.GET("/", func(ctx *Context){})
		r.GET("/v1", func(ctx *Context){}).With(onlyForV1())
		r.Group("/v1", func(r *Router) {
			r.GET("/test", func(ctx *Context){})
		})
		r.Group("/v2", func(r *Router) {
			r.Use(onlyForV2())
			r.GET("/test", func(ctx *Context){})
		})
	}
	r := NewRouter()
	r.Group("/", r2).Use(forAll())

	ts := httptest.NewServer(r)
	defer ts.Close()

	if _, body := request(t, ts, "GET", "/", nil); body != "for all" {
		t.Fatalf(body)
	}
	if _, body := request(t, ts, "GET", "/v1", nil); body != "for all and only for v1" {
		t.Fatalf(body)
	}
	if _, body := request(t, ts, "GET", "/v1/test", nil); body != "for all" {
		t.Fatalf(body)
	}
	if _, body := request(t, ts, "GET", "/v2/test", nil); body != "for all and only for v2" {
		t.Fatalf(body)
	}
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
