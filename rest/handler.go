package rest

import (
	"errors"
	"net/http"

	"github.com/ichxxx/shack"
)

// NotFoundHandler returns a handler func to respond to non-existent routes with a REST compliant
// error message.
func NotFoundHandler() shack.Handler {
	return func(ctx *shack.Context) {
		ctx.Response.Status(http.StatusNotFound)
		ctx.Response.Header("Content-Type", "application/json")
		ctx.Response.JSON(Resp(ctx).Error(errors.New("resource not found")))
	}
}

// MethodNotAllowedHandler returns a handler func to respond to routes requested with the wrong verb a
// REST compliant error message.
func MethodNotAllowedHandler() shack.Handler {
	return func(ctx *shack.Context) {
		ctx.Response.Status(http.StatusMethodNotAllowed)
		ctx.Response.Header("Content-Type", "application/json")
		ctx.Response.JSON(Resp(ctx).Error(errors.New("method not allowed")))
	}
}
