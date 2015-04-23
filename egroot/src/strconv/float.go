package strconv

import (
	"bytes"
	"math"
)

const (
	exp10start = -348
	exp10step  = 8
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

func specials64(buf []byte, f float64) (int, int) {
	if len(buf) < 4 {
		panicBuffer()
	}
	var n int
	if f < 0 {
		buf[0] = '-'
		n++
		f = -f
	}
	switch {
	case f == 0:
		return n, 0
	case math.IsInf(f, 0):
		return n + copy(buf[n:], "Inf"), 1
	case math.IsNaN(f):
		return n + copy(buf[n:], "NaN"), 1
	}
	return n, -1
}

func round(buf []byte, prec int) int {
sw:
	switch b := buf[prec]; {
	case b < '5':
		return 0
	case b > '5':
		break
	default:
		// Prefer round to even number.
		for _, b := range buf[prec+1:] {
			if b > '0' {
				break sw
			}
		}
		// For border cases prefer round to even number.
		if buf[prec-1]&1 == 0 {
			return 0
		}
	}
	// Round up.
	for i := prec - 1; i >= 0; i-- {
		if buf[i] < '9' {
			buf[i]++
			return 0
		}
		buf[i] = '0'
	}
	copy(buf[1:], buf[:prec-1])
	buf[0] = '1'
	return 1
}

func FormatFloat64(buf []byte, f float64, fmt, prec int) int {
	right := fmt < 0
	if right {
		fmt = -fmt
	}
	n := formatFloat64(buf, f, fmt, prec)
	if !right {
		bytes.Fill(buf[n:], ' ')
		return n
	}
	m := len(buf) - n
	copy(buf[m:], buf[:n])
	bytes.Fill(buf[:m], ' ')
	return m
}

func formatFloat64(buf []byte, f float64, fmt, prec int) int {
	n, spec := specials64(buf, f)
	if spec > 0 {
		return n
	}
	if n > 0 {
		f = -f
	}
	if prec < 0 {
		panic("strconv: prec<0")
	}
	switch fmt {
	case 'e', 'E':
		prec++
		l := prec + 5
		if prec > 1 {
			l++
		}
		nb := buf[n:]
		if len(nb) < l {
			panicBuffer()
		} else {
			nb = nb[:l]
		}
		if spec == 0 {
			if prec == 1 {
				return n + copy(nb, "0e+00")
			}
			nb = nb[:prec+1+4]
			bytes.Fill(nb, '0')
			nb[1] = '.'
			nb[prec+1] = byte(fmt)
			nb[prec+2] = '+'
			return n + len(nb)
		}
		exp10 := grisu(nb, f) + len(nb) - 1
		exp10 += round(nb, prec)
		n += prec
		if prec > 1 {
			copy(nb[2:], nb[1:prec])
			nb[1] = '.'
			n++
		}
		buf[n] = byte(fmt)
		n++
		if exp10 >= 0 {
			buf[n] = '+'
		} else {
			buf[n] = '-'
			exp10 = -exp10
		}
		n++
		if exp10 < 10 {
			buf[n] = '0'
			n++
		}
		return n + FormatUint(buf[n:], uint(exp10), 10)
	}
	panic("strconv: bad format")
}

// grisu accepts len(buf) >= 1
func grisu(buf []byte, f64 float64) int {
	w := normalize(makediyfp(f64))
	c10, exp10 := cachedPower(w.e, -59, -32)
	d := mul(w, c10)
	e := uint(-d.e)
	p1 := uint32(d.f >> e)
	n := 0
	for i := 9; i >= 0; i-- {
		div := cachedTens(i)
		if d := p1 / div; d != 0 || n != 0 {
			buf[n] = byte('0' + d)
			p1 -= d * div
			if n++; n >= len(buf) {
				return exp10 + i
			}
		}
	}
	mask := uint64(1)<<e - 1
	p2 := d.f & mask
	buf = buf[n:]
	for n = range buf {
		p2 *= 10
		buf[n] = '0' + byte(p2>>e)
		p2 &= mask
	}
	return exp10 - len(buf)
}

// grisu2 needs len(buf) == 18 ?????
func grisu2(buf []byte, f64 float64) (n, exp10 int) {
	low, hig := bounds(makediyfp(f64))
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
