package shack

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"net/http"
	"strings"
	"sync"
)

type Context struct {
	HttpStatusCode int
	StatusCode     *int
	Writer         http.ResponseWriter
	Request        *http.Request
	Path           string
	Method         string
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
		Path   : r.URL.Path,
		Method : r.Method,
		index  : -1,
	}
}


// String writes string to http.ResponseWriter.
func(c *Context) String(s ...string) *Context {
	c.Header("Content-Type", "text/plain")
	c.Writer.Write([]byte(strings.Join(s, "")))
	return c
}


// JSON marshals and writes data to http.ResponseWriter.
func(c *Context) JSON(data interface{}) *Context {
	c.Header("Content-Type", "application/json")
	b, _ := json.Marshal(data)
	c.Writer.Write(b)
	return c
}


// Data writes data to http.ResponseWriter.
func (c *Context) Data(data []byte) *Context {
	c.Writer.Write(data)
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


// Body returns a workflow of the request body.
func(c *Context) Body() *bodyFlow {
	if len(c.bodyBuf) > 0 {
		return newBodyFlow(c.bodyBuf)
	}

	buf, err := ioutil.ReadAll(c.Request.Body)
	if err == nil {
		c.bodyBuf = make([]byte, len(buf))
		copy(c.bodyBuf, buf)
	}
	return newBodyFlow(buf)
}


// Form returns a workflow of the first value for the named component of the POST, PATCH, or PUT request body.
func(c *Context) Form(key string) *valueFlow {
	return newValueFlow(c.Request.PostFormValue(key))
}


// Forms returns a workflow of all the values for the named component of the POST, PATCH, or PUT request body.
func(c *Context) Forms() *formFlow {
	err := c.Request.ParseMultipartForm(1024*1024*1024) // 10Mb
	if err != nil {
		return newFormFlow(nil)
	}
	return newFormFlow(c.Request.MultipartForm.Value)
}


// Query returns a workflow of the keyed url query value.
func(c *Context) Query(key string, defaultValue ...string) *valueFlow {
	value := c.Request.URL.Query().Get(key)
	if value == "" && len(defaultValue) > 0 {
		return newValueFlow(defaultValue[0])
	}

	return newValueFlow(value)
}


// RawQuery returns a workflow of the url query values, without '?'.
func(c *Context) RawQuery() *rawFlow {
	f := newRawFlow(c.Request.URL.RawQuery)
	return f
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
