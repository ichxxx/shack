package shack

import (
	"net/http"
)

const (
	_GET     = http.MethodGet
	_POST    = http.MethodPost
	_DELETE  = http.MethodDelete
	_PUT     = http.MethodPut
	_PATCH   = http.MethodPatch
	_OPTIONS = http.MethodOptions
	_HEAD    = http.MethodHead
	_ALL     = "ALL"
)
