// +build !cortexm3,!cortexm4,!cortexm4f

package bits

import "unsafe"

// TODO: Implement generic leadingZeros more efficiently.

func leadingZeros32(u uint32) uint {
	var i uint
	for u>>i != 0 {
		i++
	}
	return 32 - i
}

func leadingZeros64(u uint64) uint {
	var i uint
	for u>>i != 0 {
		i++
	}
	return 64 - i
}

func leadingZerosPtr(u uintptr) uint {
	var i uint
	for u>>i != 0 {
		i++
	}
	return uint(unsafe.Sizeof(u)*8) - i
}
