package bits

// Reverse32 returns u with reversesd bit order.
func Reverse32(u uint32) uint32 {
	return reverse32(u)
}

// Reverse64 returns u with reversesd bit order.
func Reverse64(u uint64) uint64 {
	return reverse64(u)
}
