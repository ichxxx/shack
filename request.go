package shack

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/tidwall/gjson"
)

type Request struct {
	*http.Request
	body  []byte
	query url.Values
}

func (r *Request) Header(key string) string {
	return r.Request.Header.Get(key)
}

func (r *Request) Method() string {
	return r.Request.Method
}

func (r *Request) URI() string {
	return r.Request.RequestURI
}

func (r *Request) RawQuery() string {
	return r.Request.URL.RawQuery
}

func (r *Request) Query(key string, defaultValue ...string) string {
	if r.query == nil {
		r.query = r.Request.URL.Query()
	}
	value := r.query.Get(key)
	if len(value) == 0 && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return value
}

func (r *Request) BindQuery(dst interface{}, tag ...string) error {
	p := reflect.ValueOf(dst)
	if p.Kind() != reflect.Ptr || p.IsNil() {
		return errors.New("dst must be a pointer")
	}

	m := make(map[string]string)
	segments := strings.Split(r.RawQuery(), "&")
	for _, segment := range segments {
		kv := strings.Split(segment, "=")
		if len(kv) > 1 {
			m[kv[0]] = kv[1]
		}
	}

	return mapTo(p.Elem(), m, tag...)
}

// Body returns the request body.
func (r *Request) Body() []byte {
	if r.body == nil {
		b, _ := io.ReadAll(r.Request.Body)
		r.body = b
	}
	return r.body
}

func (r *Request) BindJSON(dst interface{}) error {
	return json.Unmarshal(r.Body(), dst)
}

func (r *Request) JSON(key string) interface{} {
	return gjson.GetBytes(r.Body(), key).Value()
}

func (r *Request) Forms(key string) string {
	if r.PostForm == nil {
		_ = r.ParseMultipartForm(MaxMultipartMemory)
	}
	return r.PostForm.Get(key)
}

func (r *Request) File(key string) []*multipart.FileHeader {
	if r.PostForm == nil {
		_ = r.ParseMultipartForm(MaxMultipartMemory)
	}
	return r.MultipartForm.File[key]
}
