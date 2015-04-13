package strconv

import (
	"bits"
	"math"
)

type diyfp struct {
	f uint64
	e int
}

func normalize(x diyfp) diyfp {
	n := bits.LeadingZeros64(x.f)
	x.f <<= n
	x.e -= int(n)
	return x
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

func normalDiyfp(f64 float64) diyfp {
	bits := math.Float64bits(f64)
	frac := bits & (1<<52 - 1)
	exp := int(bits>>52) & (1<<11 - 1)
	if exp == 0 {
		exp = 1 - (1023 + 52)
	} else {
		exp -= 1023 + 52
		frac += 1 << 52
	}
	return normalize(diyfp{frac, exp})
}

func normalBounds(f64 float64) (lower, upper diyfp) {
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
		// f64 = (1<<52)*2^exp; predecessor(f64) = (1<<53-1)*2^(exp-1).
		lower.f = 4*(1<<52) - 1
		lower.e = exp - 2
	}
	upper.f = 2*frac + 1
	upper.e = exp - 1
	lower = normalize(lower)
	upper = normalize(upper)
	return
}

const (
	exp10start = -348
	exp10step  = 8
)

func cachedFrac(i int) uint64
func cachedExp(i int) int

func cachedPower(exp, alpha, gamma int) (diyfp, int) {
	exp += 64
	exp10 := ((alpha+gamma)/2 - exp + 64 - 1) * 146 / 485
	i := (exp10 - exp10start) / exp10step
	var ce int
	for {
		ce = cachedExp(i)
		if sum := exp + ce; sum < alpha {
			i++
		} else if sum > gamma {
			i--
		} else {
			break
		}
	}
	return diyfp{cachedFrac(i), ce}, exp10start + i*exp10step
}

func Show(f64 float64) (uint64, int, uint64, int) {
	w := normalDiyfp(f64)
	c10, _ := cachedPower(w.e, -59, -32)
	d := multiply(w, c10)
	return w.f, w.e, d.f, d.e
}

func FormatFloat64(buf []byte, f64 float64) int {
	if len(buf) < 25 {
		panicBuffer()
	}
	var n int
	if f64 < 0 {
		buf[0] = '-'
		n++
		f64 = -f64
	}
	// Special cases.
	switch {
	case f64 == 0:
		return n + copy(buf[n:], "0.0")
	case math.IsInf(f64, 0):
		return n + copy(buf[n:], "Inf")
	case math.IsNaN(f64):
		return n + copy(buf[n:], "NaN")
	}
	// Grisu algorithm.
	w := normalDiyfp(f64)
	c10, exp10 := cachedPower(w.e, -59, -32)
	d := multiply(w, c10)
	e := uint(-d.e)
	n += FormatUint32(buf[n:], uint32(d.f>>e), -10)
	buf[n] = '.'
	n++
	mask := uint64(1)<<e - 1
	f := d.f & mask
	for ; n < 20; n++ {
		f *= 10
		buf[n] = '0' + byte(f>>e)
		f &= mask
	}
	buf[n] = 'e'
	n++
	n += FormatInt(buf[n:], -exp10, -10)
	return n
}
