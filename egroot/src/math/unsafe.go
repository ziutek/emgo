package math

import "unsafe"

func Float32bits(f float32) uint32 { return *(*uint32)(unsafe.Pointer(&f)) }
func Float64bits(f float64) uint64 { return *(*uint64)(unsafe.Pointer(&f)) }
