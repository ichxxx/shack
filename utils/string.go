package utils

import (
	"unsafe"
)

func UnsafeString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
