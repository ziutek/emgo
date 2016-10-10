// +build !cortexm3,!cortexm4,!cortexm4f

package bits

import "unsafe"

func trailingZeros32(u uint32) uint {
	n := uint(32)
	x := u << 16
	if x != 0 {
		n -= 16
		u = x
	}
	x = u << 8
	if x != 0 {
		n -= 8
		u = x
	}
	x = u << 4
	if x != 0 {
		n -= 4
		u = x
	}
	x = u << 2
	if x != 0 {
		n -= 2
		u = x
	}
	x = u << 1
	if x != 0 {
		n -= 1
		u = x
	}
	return n - uint(u>>31)
}

func trailingZeros64(u uint64) uint {
	if x := uint32(u); x != 0 {
		return TrailingZeros32(x)
	}
	return 32 + TrailingZeros32(uint32(u>>32))
}

func trailingZerosPtr(u uintptr) uint {
	if unsafe.Sizeof(u) == 64 {
		return TrailingZeros64(uint64(u))
	}
	return TrailingZeros32(uint32(u))
}
