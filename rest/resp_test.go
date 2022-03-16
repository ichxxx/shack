package rest

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/ichxxx/shack"
)

func TestResp(t *testing.T) {
	r := shack.NewRouter()

	r.GET("/resp/1", func(ctx *shack.Context) {
		Resp(ctx).Data("foo", "foo", "bar", 123).OK()
	})
	r.GET("/resp/2", func(ctx *shack.Context) {
		DefaultFailCode(2)
		Resp(ctx).Error(errors.New("fail")).Fail()
	})
	r.GET("/resp/3", func(ctx *shack.Context) {
		data := struct {
			Foo string `json:"foo"`
		}{"bar"}
		Resp(ctx).Data(data).OK()
	})

	go shack.Run(":8080", r)

	_, data := request(t, "127.0.0.1:8080", "GET", "/resp/1", nil)
	result := map[string]interface{}{"status": 0.0, "msg": "success", "data": map[string]interface{}{"foo": "foo", "bar": 123.0}}
	if !reflect.DeepEqual(data, result) {
		t.Fatal(data)
	}

	_, data = request(t, "127.0.0.1:8080", "GET", "/resp/2", nil)
	result = map[string]interface{}{"status": 2.0, "msg": "fail", "error": "fail"}
	if !reflect.DeepEqual(data, result) {
		t.Fatal(data)
	}

	_, data = request(t, "127.0.0.1:8080", "GET", "/resp/3", nil)
	result = map[string]interface{}{"status": 0.0, "msg": "success", "data": map[string]interface{}{"foo": "bar"}}
	if !reflect.DeepEqual(data, result) {
		t.Fatal(data)
	}
}

func request(t *testing.T, url, method, path string, body io.Reader) (*http.Response, map[string]interface{}) {
	req, err := http.NewRequest(method, "http://"+url+path, body)
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
