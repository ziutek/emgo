package strconv

import (
	"bits"
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

func sub(x, y diyfp) diyfp {
	if x.e != y.e || x.f < y.f {
		panic("strconv: minus")
	}
	return diyfp{x.f - y.f, x.e}
}

func mul(x, y diyfp) diyfp {
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

func bounds(w diyfp) (lower, upper diyfp) {
	if w.f != 1<<52 || w.e == 1-(1023+52) {
		lower.f = 2*w.f - 1
		lower.e = w.e - 1
	} else {
		// Lower bound between (1<<53-1)*2^(w.e-1) and (1<<52)*2^w.e
		lower.f = 4*(1<<52) - 1
		lower.e = w.e - 2
	}
	upper.f = 2*w.f + 1
	upper.e = w.e - 1
	return
}

func expUp(x diyfp, eup int) diyfp {
	de := eup - x.e
	x.e += de
	x.f >>= uint(de)
	return x
}
