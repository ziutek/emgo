// +build cortexm3 cortexm4 cortexm4f cortexm7f cortexm7d

package bits

func trailingZeros32(u uint32) uint {
	return leadingZeros32(reverse32(u))
}

func trailingZeros64(u uint64) uint {
	return leadingZeros64(reverse64(u))
}

func trailingZerosPtr(u uintptr) uint {
	return leadingZerosPtr(reversePtr(u))
}
