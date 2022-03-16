package rest

import (
	"github.com/ichxxx/shack"
	"github.com/ichxxx/shack/middleware"
)


func Default(r *shack.Router) {
	r.Use(middleware.Recovery())
	r.Use(middleware.AccessLog())
	r.NotFound(NotFoundHandler())
	r.MethodNotAllowed(MethodNotAllowedHandler())
}
