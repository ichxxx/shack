package shack

import (
	"fmt"
	"net/http"
)

type Engine struct {
	router *Router
}

/*
func New() *Engine{
	return &Engine{}
}
*/


func ListenAndServe(addr string, router *Router) {
	err := http.ListenAndServe(addr, router)
	if err != nil {
		panic(fmt.Sprint("shack: ", err))
	}
	return
}
