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

func specials(buf []byte, f float64) (int, int) {
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

func round(buf []byte, n int, g *grisu) int {
	if n >= 19 {
		return 0
	}
sw:
	switch d := g.NextDigit(); {
	case d < 5:
		return 0
	case d > 5:
		break
	default:
		for i := 19 - 1 - n; i > 0; i-- {
			if g.NextDigit() > 0 {
				break sw
			}
		}
		// For border cases prefer to round to even number.
		if buf[n-1]&1 == 0 {
			return 0
		}
	}
	// Round up.
	for i := n - 1; i >= 0; i-- {
		if buf[i] < '9' {
			buf[i]++
			return 0
		}
		buf[i] = '0'
	}
	copy(buf[1:], buf[:n-1])
	buf[0] = '1'
	return 1
}

// FormatFloat stores text representation of f in buf using format specified by
// fmt:
//	|fmt| == 'b': -ddddp±ddd, BUG: unsupported
//  |fmt| == 'e': -d.dddde±dd,
// Unused portion of the buffer is filled with spaces.
// If fmt > 0 then formatted value is left-justified and FormatFloat returns
// its length. If base < 0 then formatted value is right-justified and
// FormatFloat returns offset to its first char.

func FormatFloat(buf []byte, f float64, fmt, prec, bitsize int) int {
	right := fmt < 0
	if right {
		fmt = -fmt
	}
	n := formatFloat(buf, f, fmt, prec, bitsize)
	if !right {
		bytes.Fill(buf[n:], ' ')
		return n
	}
	m := len(buf) - n
	copy(buf[m:], buf[:n])
	bytes.Fill(buf[:m], ' ')
	return m
}

func frexp64(f float64) (uint64, int) {
	bits := math.Float64bits(f)
	frac := bits & (1<<52 - 1)
	exp := int(bits>>52) & (1<<11 - 1)
	if exp == 0 {
		exp = 1 - (1023 + 52)
	} else {
		exp -= 1023 + 52
		frac += 1 << 52
	}
	return frac, exp
}

func frexp32(f float32) (uint64, int) {
	bits := math.Float32bits(f)
	frac := bits & (1<<23 - 1)
	exp := int(bits>>23) & (1<<8 - 1)
	if exp == 0 {
		exp = 1 - (127 + 23)
	} else {
		exp -= 127 + 23
		frac += 1 << 23
	}
	return uint64(frac), exp
}

func formatExp(buf []byte, exp int) int {
	n := 0
	if exp >= 0 {
		buf[n] = '+'
	} else {
		buf[n] = '-'
		exp = -exp
	}
	n++
	if exp < 10 {
		buf[n] = '0'
		n++
	}
	return n + FormatUint(buf[n:], uint(exp), 10)
}

func formatFloat(buf []byte, f float64, fmt, prec, bitsize int) int {
	n, spec := specials(buf, f)
	if spec > 0 {
		return n
	}
	if n > 0 {
		f = -f
	}
	if prec < 0 {
		panic("strconv: prec<0")
	}
	var (
		frac uint64
		exp  int
	)
	switch bitsize {
	case 32:
		frac, exp = frexp32(float32(f))
	case 64:
		frac, exp = frexp64(f)
	default:
		panic("strconv: illegal FormatFloat bitsize")
	}
	switch fmt {
	case 'b':
		n += FormatUint64(buf[n:], frac, 10)
		if len(buf) < n+1+1+2 {
			panicBuffer()
		}
		buf[n] = 'p'
		n++
		return n + formatExp(buf[n:], exp)
	case 'e', 'E':
		prec++ // Add first digit.
		l := prec + 5
		if prec > 1 {
			l++ // Dot.
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
		var g grisu
		g.Init(frac, exp)
		for i := 0; i < prec; i++ {
			nb[i] = byte('0' + g.NextDigit())
		}
		exp10 := g.Exp10() + prec - 1
		exp10 += round(nb, prec, &g)
		n += prec
		if prec > 1 {
			copy(nb[2:], nb[1:prec])
			nb[1] = '.'
			n++
		}
		buf[n] = byte(fmt)
		n++
		return n + formatExp(buf[n:], exp10)
	}
	panic("strconv: bad fmt")
}
