// +build !cortexm3,!cortexm4,!cortexm4f,!cortexm7f,!cortexm7d

package bits

import "unsafe"

func leadingZeros32(u uint32) uint {
	var n uint = 32
	x := u >> 16
	if x != 0 {
		n -= 16
		u = x
	}
	x = u >> 8
	if x != 0 {
		n -= 8
		u = x
	}
	x = u >> 4
	if x != 0 {
		n -= 4
		u = x
	}
	x = u >> 2
	if x != 0 {
		n -= 2
		u = x
	}
	x = u >> 1
	if x != 0 {
		n -= 1
		u = x
	}
	return n - uint(u)
}

func leadingZeros64(u uint64) uint {
	if x := uint32(u >> 32); x != 0 {
		return leadingZeros32(x)
	}
	return 32 + leadingZeros32(uint32(u))
}

func leadingZerosPtr(u uintptr) uint {
	if unsafe.Sizeof(u) == 64 {
		return leadingZeros64(uint64(u))
	}
	return leadingZeros32(uint32(u))
}
