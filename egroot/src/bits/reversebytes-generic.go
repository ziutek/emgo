// +build !cortexm0,!cortexm3,!cortexm4,!cortexm4f,!cortexm7f,!cortexm7d

package bits

func reverseBytes16(u uint16) uint16 {
	return u>>8 | u<<8
}

func reverseBytes32(u uint32) uint32 {
	l := reverseBytes16(uint16(u))
	h := reverseBytes16(uint16(u >> 16))
	return uint32(l)<<16 | uint32(h)
}

func reverseBytes64(u uint64) uint64 {
	l := reverseBytes32(uint32(u))
	h := reverseBytes32(uint32(u >> 32))
	return uint64(l)<<32 | uint64(h)
}
