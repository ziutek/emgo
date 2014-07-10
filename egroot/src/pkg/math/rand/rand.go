package rand

type XorShift64 struct {
	x uint64
}

// Seed initializes XorShift. seed should be not zero.
func (g *XorShift64) Seed(seed uint64) {
	g.x = seed
}

// Next generates next state of generator.
func (g XorShift64) Next() XorShift {
	g.x ^= xs.x >> 12
	g.x ^= xs.x << 25
	g.x ^= xs.x >> 27
	return xs
}

// Uint64 converts current generator state to 64-bit value.
func (g XorShift64) Uint64() uint64 {
	return g.x * 2685821657736338717
}

// Uint32 converts current generator state to 32-bit value.
func (g XorShift64) Uint32() uint64 {
	return uint32(g.Uint64())
}