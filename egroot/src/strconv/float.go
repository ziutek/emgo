package strconv

type diyfp struct {
	f uint64
	e int
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

func cachedPower(exp10 int) diyfp


func FormatFloat64(buf []byte, f float64) {

}
