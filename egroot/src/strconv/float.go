package strconv

import (
	"bits"
	"math"
)

type diyfp struct {
	f uint64
	e int
}

func (x *diyfp) normalize() {
	n := bits.LeadingZeros64(x.f)
	x.f <<= n
	x.e -= int(n)
}

func minus(x, y diyfp) diyfp {
	if x.e != y.e || x.f < y.f {
		panic("strconv: minus")
	}
	return diyfp{x.f - y.f, x.e}
}

func multiply(x, y diyfp) diyfp {
	const M32 = uint64(0xFFFFFFFF)
	a := x.f >> 32
	b := x.f & M32
	c := y.f >> 32
	d := y.f & M32
	ac := a * c
	bc := b * c
	ad := a * d
	bd := b * d
	tmp := bd>>32 + ad&M32 + bc&M32 + 1<<31
	return diyfp{
		ac + ad>>32 + bc>>32 + tmp>>32,
		x.e + y.e + 64,
	}
}

func bounds(f64 float64) (lower, upper diyfp) {
	bits := math.Float64bits(f64)
	frac := bits & (1<<52 - 1)
	exp := int(bits>>52) & (1<<11 - 1)
	if exp == 0 {
		exp = 1 - (1023 + 52)
	} else {
		exp -= 1023 + 52
		frac += 1 << 52
	}
	if frac != 1<<52 || exp == 1-(1023+52) {
		lower.f = 2*frac - 1
		lower.e = exp - 1
	} else {
		// Lower border between (1<<53-1)*2^(exp-1) and (1<<52)*2^exp.
		lower.f = 4*(1<<52) - 1
		lower.e = exp - 2
	}
	upper.f = 2*frac + 1
	upper.e = exp - 1
	lower.normalize()
	upper.normalize()
	return
}

func cachedPower(exp10 int) diyfp

func exp10comp(exp, alpha, gamma int) int {
	exp10 := (alpha - exp + 64 - 1) * 485 / 146
1
}

func Show(f float64) (uint64, int, uint64, int) {
	l, u := bounds(f)
	return l.f, l.e, u.f, u.e
}

func FormatFloat64(buf []byte, f float64) {
	if f < 0 {
		f = -f
		n := copy(buf, "-")
		buf = buf[n:]
	}
	// Special cases.
	switch {
	case f == 0:
		copy(buf, "0.0")
	case math.IsInf(f, 0):
		copy(buf, "Inf")
	case math.IsNaN(f):
		copy(buf, "NaN")
	}
}
