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

// Uint64 returns pseudorandom uint64 number.
func (g XorShift64) Uint64() uint64 {
	x := g.x
	x ^= x >> 12
	x ^= x << 25
	x ^= x >> 27
	g.x = x
	return x * 2685821657736338717
}

// Int64 returns pseudorandom int64 number.
func (g XorShift64) Int64() int64 {
	return int64(g.Uint64())
}

// Uint32 returns pseudorandom uint32 number.
func (g XorShift64) Uint32() uint32 {
	return uint32(g.Uint64())
}

// Int32 returns pseudorandom int32 number.
func (g XorShift64) Int32() int32 {
	return int32(g.Uint64())
}