package rand

type XorShift64 struct {
	x uint64
}

// Seed initializes XorShift. seed should not be zero.
func (g *XorShift64) Seed(seed uint64) {
	g.x = seed
}

// Next seteps generator to next state.
func (g XorShift64) Next() {
	g.x ^= g.x >> 12
	g.x ^= g.x << 25
	g.x ^= g.x >> 27
}

// Uint64 converts current generator state to unsigned 64-bit integer.
func (g XorShift64) Uint64() uint64 {
	return g.x * 2685821657736338717
}

// Uint64 converts current generator state to signed 64-bit integer.
func (g XorShift64) Int64() int64 {
	return int64(g.Uint64())
}

// Uint32 converts current generator state to unsigned 32-bit integer.
func (g XorShift64) Uint32() uint32 {
	return uint32(g.Uint64())
}

// Int32 converts current generator state to signed 32-bit integer.
func (g XorShift64) Int32() int32 {
	return int32(g.Uint64())
}
