package shack

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"golang.org/x/sync/errgroup"
)

var (
	Json = jsoniter.ConfigCompatibleWithStandardLibrary
	server = make(map[string]*fasthttp.Server)
	serverMutex = sync.RWMutex{}
)

type M map[string]interface{}

func Run(addr string, router *Router) error {
	s := &fasthttp.Server{
		Handler: router.ServeHTTP,
	}
	serverMutex.Lock()
	server[addr] = s
	serverMutex.Unlock()

	go func() {
		<- listenShutdown()
		shutdown(s)
	}()

	return s.ListenAndServe(addr)
}

func Stop(addr... string) error {
	errGroup := errgroup.Group{}
	serverMutex.RLock()
	if len(addr) > 0 {
		for _, a := range addr {
			_addr := a
			errGroup.Go(func() error {
				return server[_addr].Shutdown()
			})
		}
	} else {
		for _, s := range server {
			_server := s
			errGroup.Go(func() error {
				return _server.Shutdown()
			})
		}
	}
	serverMutex.RUnlock()
	return errGroup.Wait()
}

func listenShutdown() <- chan os.Signal {
	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGQUIT)
	return shutdownSignal
}

func shutdown(server *fasthttp.Server) {
	err := server.Shutdown()
	if err != nil {
		log.Println("shack: server shutdown failed,", err.Error())
	}
	log.Println("shack: shutdown")
}

// Logger returns a logger by specify a name
func Logger(name string) *logger {
	Log.name = name
	return Log
}


type _Router interface {
	ServeHTTP(w http.ResponseWriter, req *http.Request)

	GET(pattern string, handler HandlerFunc) *trie
	POST(pattern string, handler HandlerFunc) *trie
	DELETE(pattern string, handler HandlerFunc) *trie
	PUT(pattern string, handler HandlerFunc) *trie
	PATCH(pattern string, handler HandlerFunc) *trie
	OPTIONS(pattern string, handler HandlerFunc) *trie
	HEAD(pattern string, handler HandlerFunc) *trie
	Handle(pattern string, fn HandlerFunc, methods ...string)

	Use(middlewares ...HandlerFunc)
	Mount(pattern string, router *Router)
	Group(pattern string, fn func(r *Router)) *Router
	Add(fn func(r *Router)) *Router

	NotFound(handler HandlerFunc)
	MethodNotAllowed(handler HandlerFunc)
}


type _Logger interface {
	Enable()

	Level(level int8) *logger
	Encoding(encoding string) *logger
	Output(paths string) *logger
	Dev() *logger

	Debug(msg string, keyAndValues ...interface{})
	Info(msg string, keyAndValues ...interface{})
	Warn(msg string, keyAndValues ...interface{})
	Error(msg string, keyAndValues ...interface{})
	Panic(msg string, keyAndValues ...interface{})
	Fatal(msg string, keyAndValues ...interface{})
}


type _Context interface {
	Status(code int) *Context
	HttpStatus(code int) *Context
	Header(key string, value string) *Context
	String(s ...string) *Context
	JSON(data interface{}) *Context
	Data(data []byte) *Context

	Param(key string) string
	Body() *bodyFlow
	Form(key string) *valueFlow
	Forms() *formFlow
	Query(key string, defaultValue ...string) *valueFlow
	RawQuery() *rawFlow

	Set(key string, value interface{})
	Get(key string) (value interface{}, ok bool)
	Delete(key string)

	SetSync(key string, value interface{})
	GetSync(key string) (value interface{}, ok bool)
	DeleteSync(key string)

	Error(err error)
	Abort()
	Next()
}


type _Flow interface {
	Value()
	Int() int
	Int8() int8
	Int64() int64
	Float64() float64
	Bool() bool
	BindJson(dst interface{}) error
	Bind(dst interface{}, tag ...string) error
}


type _Config interface {
	File(file string) *configManager
	Add(config config, section string)
	Load()
}
