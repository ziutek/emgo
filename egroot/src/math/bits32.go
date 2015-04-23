package math

const (
	uvnan32    = 0x7F880001
	uvinf32    = 0x7F800000
	uvneginf32 = 0xFF800000
	shift32    = 32 - 8 - 1
	bias32     = 127
)

// Inf32 is like Inf but for float32.
func Inf32(sign int) float32 {
	var v uint32
	if sign >= 0 {
		v = uvinf32
	} else {
		v = uvneginf32
	}
	return Float32frombits(v)
}

// NaN32 is like NaN but for float32.
func NaN32() float32 { return Float32frombits(uvnan32) }

// IsNaN32 is like IsNaN but for float32.
func IsNaN32(f float32) (is bool) {
	x := Float32bits(f)
	return x&uvinf32 == uvinf32 && x != uvinf32 && x != uvneginf32
}

// IsInf32 is like IsInf but for float32.
func IsInf32(f float32, sign int) bool {
	x := Float32bits(f)
	return sign >= 0 && x == uvinf32 || sign <= 0 && x == uvneginf32
}

// normalize32 is like normalize but for float32.
func normalize32(x float32) (y float32, exp int) {
	const SmallestNormal = 1.1754944e-38
	if Abs32(x) < SmallestNormal {
		return x * (1 << 23), -23
	}
	return x, 0
}
