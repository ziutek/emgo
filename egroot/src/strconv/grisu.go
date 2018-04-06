package strconv

const (
	exp10start = -348
	exp10step  = 8
	maxDigits  = 19
)

func cachedPower(exp, alpha, gamma int) (diyfp, int) {
	exp += 64
	exp10 := ((alpha+gamma)/2 - exp + 64 - 1) * 146 / 485
	i := (exp10 - exp10start) / exp10step
	var ce int
	for {
		ce = int(exponents[i])
		if sum := exp + ce; sum < alpha {
			i++
		} else if sum > gamma {
			i--
		} else {
			break
		}
	}
	return diyfp{significands[i], ce}, -(exp10start + i*exp10step)
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
		div := tens[g.i]
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
		div := tens[g.i]
		dig := g.p1 / div
		g.p1 -= dig * div
		return int(dig)
	}
	g.p2 *= 10
	dig := int(g.p2 >> g.e)
	g.p2 &= g.mask
	return dig
}
