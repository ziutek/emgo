package bits

func LeadingZeros32(u uint32) uint {
	return leadingZeros32(u)
}

func LeadingZeros64(u uint64) uint {
	return leadingZeros64(u)
}

func LeadingZerosPtr(u uintptr) uint {
	return leadingZerosPtr(u)
}
