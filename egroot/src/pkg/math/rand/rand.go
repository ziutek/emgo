// Package rand provides functions that can be used to generate pseudorandom
// numbers.
package rand

// XorShift64 is 64-bit xorshift* pseudorandom number generator.
// See http://en.wikipedia.org/wiki/Xorshift for more informations.
type XorShift64 struct {
	x uint64
}

// Seed initializes XorShift64 state. seed must not be zero..
func (g *XorShift64) Seed(seed uint64) {
	g.x = seed
}

// Next seteps generator to next state.
func (g *XorShift64) Next() {
	x := g.x
	x ^= x >> 12
	x ^= x << 25
	x ^= x >> 27
	g.x = x
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
