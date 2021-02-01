package rest

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/ichxxx/shack"
)


func TestResp(t *testing.T) {
	r := shack.NewRouter()

	r.GET("/resp/1", func(ctx *shack.Context){
		R(ctx).Data("foo", "foo", "bar", 123).OK()
	})
	r.GET("/resp/2", func(ctx *shack.Context){
		DefaultFailCode(2)
		R(ctx).Error(errors.New("fail")).Fail()
	})
	r.GET("/resp/3", func(ctx *shack.Context){
		data := struct {
			Foo string `json:"foo"`
		}{"bar"}
		R(ctx).Data(data).Write()
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	_, data := request(t, ts, "GET", "/resp/1", nil)
	result := map[string]interface{}{"code":0.0, "msg":"success", "data": map[string]interface{}{"foo":"foo", "bar":123.0}}
	if !reflect.DeepEqual(data, result) {
		t.Fatal(data)
	}

	_, data = request(t, ts, "GET", "/resp/2", nil)
	result = map[string]interface{}{"code":2.0, "msg":"fail", "error":"fail"}
	if !reflect.DeepEqual(data, result) {
		t.Fatal(data)
	}

	_, data = request(t, ts, "GET", "/resp/3", nil)
	result = map[string]interface{}{"data":map[string]interface{}{"foo":"bar"}}
	if !reflect.DeepEqual(data, result) {
		t.Fatal(data)
	}
}


func TestRespWithCtx(t *testing.T) {
	r := shack.NewRouter()

	r.GET("/resp/1", func(ctx *shack.Context){
		R(ctx).OK()
	})
	r.GET("/resp/2", func(ctx *shack.Context){
		DefaultFailCode(2)
		DefaultFailMsg("test_2")
		R(ctx).Fail()
	})
	r.GET("/resp/3", func(ctx *shack.Context){
		ctx.Status(3)
		ctx.Error(errors.New("test_3"))
		R(ctx).Fail()
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	_, data := request(t, ts, "GET", "/resp/1", nil)
	result := map[string]interface{}{"code":0.0, "msg":"success"}
	if !reflect.DeepEqual(data, result) {
		t.Fatal(data)
	}

	_, data = request(t, ts, "GET", "/resp/2", nil)
	result = map[string]interface{}{"code":2.0, "msg":"test_2"}
	if !reflect.DeepEqual(data, result) {
		t.Fatal(data)
	}

	_, data = request(t, ts, "GET", "/resp/3", nil)
	result = map[string]interface{}{"code":3.0, "msg":"test_2", "error": "test_3"}
	if !reflect.DeepEqual(data, result) {
		t.Fatal(data)
	}
}


func request(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, map[string]interface{}) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	if err != nil {
		t.Fatal(err)
		return nil, nil
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
		return nil, nil
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
		return nil, nil
	}
	defer resp.Body.Close()

	data := make(map[string]interface{})
	json.Unmarshal(respBody, &data)
	return resp, data
}
