package bits

func TrailingZeros32(u uint32) uint {
	return trailingZeros32(u)
}

func TrailingZeros64(u uint64) uint {
	return trailingZeros64(u)
}

func TrailingZerosPtr(u uintptr) uint {
	return trailingZerosPtr(u)
}
