package shack

import (
	"net/http"

	"github.com/ichxxx/shack/utils"
	"github.com/valyala/bytebufferpool"
)

var responseBodyPool bytebufferpool.Pool

type Response struct {
	http.ResponseWriter
	StatusCode int
	body       *bytebufferpool.ByteBuffer
	hasFlush   bool
}

func (r *Response) Header(key, value string) {
	r.ResponseWriter.Header().Set(key, value)
}

func (r *Response) Write(data []byte) error {
	_, err := r.bodyBuffer().Write(data)
	return err
}

func (r *Response) String(s string) error {
	r.Header("Content-Type", "text/plain")
	_, err := r.bodyBuffer().Write(utils.UnsafeBytes(s))
	return err
}

func (r *Response) JSON(data interface{}) error {
	r.Header("Content-Type", "application/json")

	bytes, err := getBytes(data)
	if err != nil {
		return err
	}
	return r.Write(bytes)
}

func (r *Response) Flush() error {
	if r.hasFlush {
		return nil
	}
	r.hasFlush = true
	if r.StatusCode != 0 {
		r.ResponseWriter.WriteHeader(r.StatusCode)
	}
	_, err := r.ResponseWriter.Write(r.body.Bytes())
	return err
}

func getBytes(v interface{}) ([]byte, error) {
	switch d := v.(type) {
	case []byte:
		return d, nil
	case string:
		return utils.UnsafeBytes(d), nil
	default:
		return json.Marshal(v)
	}
}

// Status sets the http status of response.
func (r *Response) Status(code int) {
	r.StatusCode = code
}

func (r *Response) bodyBuffer() *bytebufferpool.ByteBuffer {
	if r.body == nil {
		r.body = responseBodyPool.Get()
	}
	return r.body
}
