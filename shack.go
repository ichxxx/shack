package shack

import (
	"fmt"
	"net/http"
)


func ListenAndServe(addr string, router *Router) {
	err := http.ListenAndServe(addr, router)
	if err != nil {
		panic(fmt.Sprint("shack: ", err))
	}
	return
}


func Logger(name string) *logger {
	Log.name = name
	return Log
}
