package bits

func Field32(bits, mask uint32) int {
	n := TrailingZeros32(mask)
	return int(bits & mask >> n)
}

func MakeField32(v int, mask uint32) uint32 {
	n := TrailingZeros32(mask)
	return uint32(v<<n) & mask
}

func Field64(bits, mask uint64) int {
	n := TrailingZeros64(mask)
	return int(bits & mask >> n)
}

func MakeField64(v int, mask uint64) uint64 {
	n := TrailingZeros64(mask)
	return uint64(v<<n) & mask

}
