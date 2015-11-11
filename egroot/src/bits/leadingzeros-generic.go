// +build !cortexm3,!cortexm4,!cortexm4f

package bits

import "unsafe"

// TODO: Implement these functions more efficiently.

func leadingZeros32(u uint32) uint {
	n := uint(32)
	for u != 0 {
		u >>= 1
		n--
	}
	return n
}

func leadingZeros64(u uint64) uint {
	n := uint(64)
	for u != 0 {
		u >>= 1
		n--
	}
	return n
}

func leadingZerosPtr(u uintptr) uint {
	n := uint(unsafe.Sizeof(u) * 8)
	for u != 0 {
		u >>= 1
		n--
	}
	return n
}
