package bits

// ReverseBytes16 returns u with reversed byte order.
func ReverseBytes16(u uint16) uint16 {
	return reverseBytes16(u)
}

// ReverseBytes32 returns u with reversed byte order.
func ReverseBytes32(u uint32) uint32 {
	return reverseBytes32(u)
}

// ReverseBytes64 returns u with reversed byte order.
func ReverseBytes64(u uint64) uint64 {
	return reverseBytes64(u)
}