package shack

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strings"
	"sync"
	"unsafe"
)

type Context struct {
	HttpStatusCode int
	StatusCode     *int
	Writer         http.ResponseWriter
	Request        *http.Request
	Params         map[string]string
	Bucket         map[string]interface{}
	SyncBucket     *sync.Map
	Err            error
	errOnce        *sync.Once
	bodyBuf        []byte
	handlers       []HandlerFunc
	index          int8
}

const abortIndex int8 = math.MaxInt8 / 2


func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer : w,
		Request: r,
		index  : -1,
	}
}


// String writes string to http.ResponseWriter.
func(c *Context) String(s ...string) *Context {
	c.Header("Content-Type", "text/plain")
	_, err := c.Writer.Write(bytesFromString(strings.Join(s, "")))
	if err != nil {
		log.Printf("shack: ResponseWriter write error, %s", err.Error())
	}
	return c
}


// JSON marshals and writes data to http.ResponseWriter.
// Support raw json (string or []byte), struct and map.
func(c *Context) JSON(data interface{}) *Context {
	c.Header("Content-Type", "application/json")

	if b, ok := getBytes(data); ok {
		_, err := c.Writer.Write(b)
		if err != nil {
			log.Printf("shack: ResponseWriter write error, %s", err.Error())
		}
		return c
	}

	b, err := json.Marshal(data)
	if err != nil {
		log.Printf("shack: marshal json error, %s", err.Error())
	}
	_, err = c.Writer.Write(b)
	if err != nil {
		log.Printf("shack: ResponseWriter write error, %s", err.Error())
	}
	return c
}


// Data writes data to http.ResponseWriter.
func (c *Context) Data(data []byte) *Context {
	_, err := c.Writer.Write(data)
	if err != nil {
		log.Printf("shack: ResponseWriter write error, %s", err.Error())
	}
	return c
}


// Status sets the status of response.
func(c *Context) Status(code int) *Context {
	c.StatusCode = &code
	return c
}


// HttpStatus sets the http status of response.
func(c *Context) HttpStatus(code int) *Context {
	if code < 100 || code > 500 {
		return c
	}
	c.HttpStatusCode = code
	c.Writer.WriteHeader(code)
	return c
}


// Header sets the header of response.
func(c *Context) Header(key string, value string) *Context {
	c.Writer.Header().Set(key, value)
	return c
}


// Param returns the value of the url param.
func(c *Context) Param(key string) string {
	return c.Params[key]
}


// Body returns the request body.
func(c *Context) Body() []byte {
	if len(c.bodyBuf) > 0 {
		return c.bodyBuf
	}

	buf, err := ioutil.ReadAll(c.Request.Body)
	if err == nil {
		c.bodyBuf = make([]byte, len(buf))
		copy(c.bodyBuf, buf)
	}
	return buf
}


// BodyFlow returns a workflow of the request body.
func(c *Context) BodyFlow() bodyFlow {
	return newBodyFlow(c.Body())
}


// Form returns the first value for the named component of the POST, PATCH, or PUT request body.
func(c *Context) Form(key string) string {
	return c.Request.PostFormValue(key)
}


// FormFlow returns a workflow of the first value for the named component of the POST, PATCH, or PUT request body.
func(c *Context) FormFlow(key string) valueFlow {
	return newValueFlow(c.Form(key))
}


// Forms returns all the values for the named component of the POST, PATCH, or PUT request body.
func(c *Context) Forms() map[string][]string {
	err := c.Request.ParseMultipartForm(1024*1024*1024) // 10Mb
	if err != nil {
		return nil
	}
	return c.Request.MultipartForm.Value
}


// FormsFlow returns a workflow of all the values for the named component of the POST, PATCH, or PUT request body.
func(c *Context) FormsFlow() formFlow {
	return newFormFlow(c.Forms())
}


// Query returns a workflow of the keyed url query value.
func(c *Context) Query(key string, defaultValue ...string) string {
	value := c.Request.URL.Query().Get(key)
	if value == "" && len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return value
}


// QueryFlow returns a workflow of the keyed url query value.
func(c *Context) QueryFlow(key string, defaultValue ...string) valueFlow {
	return newValueFlow(c.Query(key, defaultValue...))
}


// RawQuery returns the url query values, without '?'.
func(c *Context) RawQuery() string {
	return c.Request.URL.RawQuery
}


// RawQueryFlow returns a workflow of the url query values, without '?'.
func(c *Context) RawQueryFlow() rawFlow {
	return newRawFlow(c.RawQuery())
}


// Error sets the first non-nil error of the context.
func(c *Context) Error(err error) {
	if err != nil {
		if c.errOnce == nil {
			c.errOnce = &sync.Once{}
		}
		c.errOnce.Do(func() {
			c.Err = err
		})
	}
}


// SetSync stores a key/value pair in the context bucket synchronicity.
func(c *Context) SetSync(key string, value interface{}) {
	if c.SyncBucket == nil {
		c.SyncBucket = &sync.Map{}
	}
	c.SyncBucket.Store(key, value)
}


// GetSync returns the value for the given key in the context bucket synchronicity.
func(c *Context) GetSync(key string) (value interface{}, ok bool) {
	if c.SyncBucket == nil {
		return
	}
	return c.SyncBucket.Load(key)
}


// DeleteSync removes the value for the given key in the context bucket synchronicity.
func(c *Context) DeleteSync(key string) {
	if c.SyncBucket == nil {
		return
	}
	c.SyncBucket.Delete(key)
	return
}


// Set stores a key/value pair in the context bucket.
func(c *Context) Set(key string, value interface{}) {
	if c.Bucket == nil {
		c.Bucket = make(map[string]interface{})
	}
	c.Bucket[key] = value
}


// Get returns the value for the given key in the context bucket.
func(c *Context) Get(key string) (value interface{}, ok bool) {
	if c.Bucket == nil {
		return
	}
	value, ok = c.Bucket[key]
	return
}


// Delete removes the value for the given key in the context bucket.
func(c *Context) Delete(key string) {
	if c.Bucket == nil {
		return
	}
	delete(c.Bucket, key)
}


// Abort prevents pending handlers from being called.
func(c *Context) Abort() {
	c.index = abortIndex
}


// Next should be used only inside middleware.
// It executes the pending handlers in the chain inside the calling handler.
func(c *Context) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}


func getBytes(v interface{}) (b []byte, ok bool) {
	switch d := v.(type) {
	case []byte:
		return d, true
	case string:
		return bytesFromString(d), true
	}
	return nil, false
}


func bytesFromString(s string) []byte {
	tmp := (*[2]uintptr)(unsafe.Pointer(&s))
	x := [3]uintptr{tmp[0], tmp[1], tmp[1]}
	return *(*[]byte)(unsafe.Pointer(&x))
}
