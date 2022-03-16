package utils

import (
	"unsafe"
)

func UnsafeBytes(s string) []byte {
	tmp := (*[2]uintptr)(unsafe.Pointer(&s))
	x := [3]uintptr{tmp[0], tmp[1], tmp[1]}
	return *(*[]byte)(unsafe.Pointer(&x))
}
