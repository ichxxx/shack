package shack

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"
)


func forAll() HandlerFunc {
	return func(c *Context) {
		c.String("for all")
	}
}


func onlyForV1() HandlerFunc {
	return func(c *Context) {
		c.String(" and ")
		c.String("only for v1")
	}
}



func onlyForV2() HandlerFunc {
	return func(c *Context) {
		c.String(" and ")
		c.String("only for v2")
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

		r.Group("/v3", func(r *Router) {
			r.Add(func(r *Router) {
				r.GET("/test", func(ctx *Context) {
					ctx.String(" and ")
					ctx.String("this is v3")
				})
			})
		})
	}
	r := NewRouter()
	r.Group("/", r2).Use(forAll())

	go Run(":8080", r)

	if _, body := request(t, "127.0.0.1:8080", "GET", "/", nil); body != "for all" {
		t.Fatalf(body)
	}
	if _, body := request(t, "127.0.0.1:8080", "GET", "/v1", nil); body != "for all and only for v1" {
		t.Fatalf(body)
	}
	if _, body := request(t, "127.0.0.1:8080", "GET", "/v1/test", nil); body != "for all" {
		t.Fatalf(body)
	}
	if _, body := request(t, "127.0.0.1:8080", "GET", "/v2/test", nil); body != "for all and only for v2" {
		t.Fatalf(body)
	}
	if _, body := request(t, "127.0.0.1:8080", "GET", "/v3/test", nil); body != "for all and this is v3" {
		t.Fatalf(body)
	}
}


func request(t *testing.T, url, method, path string, body io.Reader) (*http.Response, string) {
	now := time.Now()
	defer log.Printf("%s in %v", path, time.Since(now).Microseconds())

	req, err := http.NewRequest(method, "http://"+url+path, body)
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
