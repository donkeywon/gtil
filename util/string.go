package util

import "unsafe"

func Bytes2String(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}

func String2Bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}
