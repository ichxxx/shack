package rest

import (
	"errors"
	"net/http"

	"github.com/ichxxx/shack"
)

// NotFoundHandler returns a handler func to respond to non existent routes with a REST compliant
// error message.
func NotFoundHandler() shack.HandlerFunc {
	return func(ctx *shack.Context) {
		ctx.HttpStatus(http.StatusNotFound)
		ctx.Header("Content-Type", "application/json")
		ctx.JSON(Resp(ctx).Error(errors.New("resource not found")))
	}
}

// MethodNotAllowedHandler returns a handler func to respond to routes requested with the wrong verb a
// REST compliant error message.
func MethodNotAllowedHandler() shack.HandlerFunc {
	return func(ctx *shack.Context) {
		ctx.HttpStatus(http.StatusMethodNotAllowed)
		ctx.Header("Content-Type", "application/json")
		ctx.JSON(Resp(ctx).Error(errors.New("method not allowed")))
	}
}
