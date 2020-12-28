package rest

import (
	"github.com/ichxxx/shack"
)

// NotFoundHandler returns a handler func to respond to non existent routes with a REST compliant
// error message
func NotFoundHandler() shack.HandlerFunc {
	return func(ctx *shack.Context) {
		ctx.JSON(Resp().Error("resource not found"))
	}
}

// MethodNotAllowedHandler returns a handler func to respond to routes requested with the wrong verb a
// REST compliant error message
func MethodNotAllowedHandler() shack.HandlerFunc {
	return func(ctx *shack.Context) {
		ctx.JSON(Resp().Error("method not allowed"))
	}
}
