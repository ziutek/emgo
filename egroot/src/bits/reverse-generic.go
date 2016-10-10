// +build !cortexm3,!cortexm4,!cortexm4f

package bits

import "unsafe"

func reverse32(u uint32) uint32 {
	var v uint32
	for i := 31; u != 0 && i >= 0; i-- {
		v |= u & 1 << uint(i)
		u >>= 1
	}
	return v
}

func reverse64(u uint64) uint64 {
	var v uint64
	for i := 63; u != 0 && i >= 0; i-- {
		v |= u & 1 << uint(i)
		u >>= 1
	}
	return v
}

func reversePtr(u uintptr) uintptr {
	if unsafe.Sizeof(u) == 64 {
		return uintptr(reverse64(uint64(u)))
	}
	return uintptr(reverse32(uint32(u)))
}
