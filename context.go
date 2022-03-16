package shack

import (
	"math"
	"net/http"
	"sync"
)

const abortIndex int8 = math.MaxInt8 / 2

var (
	ctxPool = &sync.Pool{New: func() interface{} { return new(Context) }}
)

type Context struct {
	index       int8
	Request     Request
	Response    Response
	PathParams  map[string]string
	handlers    []Handler
	Err         error
	errOnce     *sync.Once
	Bucket      map[string]interface{}
	bucketMutex *sync.RWMutex
}

func getContext(request *http.Request, response http.ResponseWriter) *Context {
	ctx := ctxPool.Get().(*Context)
	ctx.reset()
	ctx.Request = Request{Request: request}
	ctx.Response = Response{ResponseWriter: response}
	ctx.index = -1
	ctx.errOnce = &sync.Once{}
	ctx.bucketMutex = &sync.RWMutex{}
	return ctx
}

func releaseContext(ctx *Context) {
	ctxPool.Put(ctx)
}

func (c *Context) reset() {
	c.PathParams = nil
	c.handlers = nil
	c.Err = nil
	c.errOnce = nil
	c.Bucket = nil
	c.bucketMutex = nil
}

// Set stores a key/value pair in the context bucket.
func (c *Context) Set(key string, value interface{}) {
	c.bucketMutex.Lock()
	defer c.bucketMutex.Unlock()

	if c.Bucket == nil {
		c.Bucket = make(map[string]interface{})
	}
	c.Bucket[key] = value
}

// Get returns the value for the given key in the context bucket.
func (c *Context) Get(key string) (value interface{}, ok bool) {
	c.bucketMutex.RLock()
	defer c.bucketMutex.RUnlock()

	if c.Bucket == nil {
		return
	}
	value, ok = c.Bucket[key]
	return
}

// Error sets the first non-nil error of the context.
func (c *Context) Error(err error) {
	if err != nil {
		c.errOnce.Do(func() {
			c.Err = err
		})
	}
}

// Next should be used only inside middleware.
// It executes the pending handlers in the chain inside the calling handler.
func (c *Context) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
	_ = c.Response.Flush()
}

// Abort prevents pending handlers from being called.
func (c *Context) Abort() {
	c.index = abortIndex
}
