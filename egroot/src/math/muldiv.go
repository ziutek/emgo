package math

// MulDiv returns uint64(uint128(x) * uint128(m) / uint128(d)).
func MulDiv(x, m, d uint64) uint64 {
	divx := x / d
	modx := x - divx*d
	divm := m / d
	modm := m - divm*d
	return divx*m + modx*divm + modx*modm/d
}

func MulDivUp(x, m, d uint64) uint64 {
	o := d - 1
	divx := (x + o) / d
	modx := x - divx*d
	divm := (m + o) / d
	modm := m - divm*d
	return divx*m + modx*divm + (modx*modm+o)/d
}
