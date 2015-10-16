package strconv

const (
	exp10start = -348
	exp10step  = 8
	maxDigits  = 19
)

func cachedFrac(i int) uint64
func cachedExp(i int) int
func cachedTens(i int) uint32

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
	return diyfp{cachedFrac(i), ce}, -(exp10start + i*exp10step)
}

// grisu implements Grisu algorithm published by Florian Loitsch in paper
// "Printing floating-point numbers quickly and accurately with integers".
type grisu struct {
	p2   uint64
	mask uint64
	i    int
	e    uint
	p1   uint32
}

// Init initializes g and returns first digit and its exp10.
func (g *grisu) Init(frac uint64, exp int) (int, int) {
	w := normalize(diyfp{frac, exp})
	d, exp10 := cachedPower(w.e, -59, -32)
	d = mul(w, d)
	g.i = 9
	g.e = uint(-d.e)
	g.p1 = uint32(d.f >> g.e)
	g.mask = uint64(1)<<g.e - 1
	g.p2 = d.f & g.mask
	for {
		div := cachedTens(g.i)
		dig := g.p1 / div
		if dig != 0 {
			g.p1 -= dig * div
			return int(dig), exp10 + g.i
		}
		g.i--
	}
}

func (g *grisu) NextDigit() int {
	if g.i > 0 {
		g.i--
		div := cachedTens(g.i)
		dig := g.p1 / div
		g.p1 -= dig * div
		return int(dig)
	}
	g.p2 *= 10
	dig := int(g.p2 >> g.e)
	g.p2 &= g.mask
	return dig
}

/*func (g *grisu) Exp10() int {
	return g.exp10 + g.i
}*/

// grisu2 needs len(buf) >= 20
func grisu2(buf []byte, frac uint64, exp int) (n, exp10 int) {
	low, hig := bounds(diyfp{frac, exp})
	hig = normalize(hig)
	low = expUp(normalize(low), hig.e)
	var c10 diyfp
	c10, exp10 = cachedPower(hig.e, -59, -32)
	low = mul(low, c10)
	hig = mul(hig, c10)
	low.f++
	hig.f--
	delta := sub(hig, low)
	// Digits generation.
	e := uint(-hig.e)
	mask := uint64(1)<<e - 1
	p1 := uint32(hig.f >> e)
	p2 := hig.f & mask
	kappa := 10
	for kappa > 0 {
		kappa--
		div := cachedTens(kappa)
		if d := p1 / div; d != 0 || n != 0 {
			buf[n] = byte('0' + d)
			p1 -= d * div
			n++
		}
		if uint64(p1)<<e+p2 <= delta.f {
			goto end
		}
	}
	for {
		kappa--
		p2 *= 10
		if d := int(p2 >> e); d != 0 || n != 0 {
			buf[n] = byte('0' + d)
			n++
		}
		p2 &= mask
		delta.f *= 10
		if p2 <= delta.f {
			break
		}
	}
end:
	// TODO: rounding
	exp10 += kappa
	return
}
