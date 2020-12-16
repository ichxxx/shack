package shack

import (
	"fmt"
	"net/http"
)

type Engine struct {
	router       *Router
	middlewares  []func(http.Handler) http.Handler
}


func New() *Engine{
	return &Engine{}
}


func(e *Engine) Use(middleware func(http.Handler) http.Handler) {
	e.middlewares = append(e.middlewares, middleware)
}


func(e *Engine) ListenAndServe(addr string, router *Router) {
	e.router = router
	err := http.ListenAndServe(addr, e.router)
	if err != nil {
		panic(fmt.Sprint("shack: ", err))
	}
	return
}
