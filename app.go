package shack

import (
	"context"
	"net/http"
	"sync"

	jsoniter "github.com/json-iterator/go"
)

type (
	// Map is a shortcut for map[string]interface{}
	Map     map[string]interface{}
	Handler func(*Context)
)

var (
	json                     = jsoniter.ConfigCompatibleWithStandardLibrary
	MaxMultipartMemory int64 = 8 << 20
	runningApps              = make(map[string]*http.Server)
	appMutex                 = sync.Mutex{}
)

type Option struct {
	ShutdownFunc func()
}

var defaultOpt = Option{
	ShutdownFunc: func() {},
}

func Run(addr string, router *Router, opts ...Option) error {
	if len(opts) == 0 {
		opts = append(opts, defaultOpt)
	}
	for _, opt := range opts {
		defer opt.ShutdownFunc()
	}

	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}
	addRunningApp(addr, server)
	return server.ListenAndServe()
}

func addRunningApp(addr string, app *http.Server) {
	appMutex.Lock()
	runningApps[addr] = app
	appMutex.Unlock()
}

func Shutdown(addrs ...string) {
	appMutex.Lock()
	defer appMutex.Unlock()

	wg := sync.WaitGroup{}
	if len(addrs) > 0 {
		wg.Add(len(addrs))
		for _, addr := range addrs {
			_addr := addr
			go func() {
				_ = shutdownRunningApp(runningApps[_addr])
				wg.Done()
			}()
			delete(runningApps, addr)
		}

	} else {
		wg.Add(len(runningApps))
		for _, app := range runningApps {
			_app := app
			go func() {
				_ = shutdownRunningApp(_app)
				wg.Done()
			}()
		}
		runningApps = make(map[string]*http.Server)
	}

	wg.Wait()
}

func shutdownRunningApp(app *http.Server) error {
	return app.Shutdown(context.Background())
}
