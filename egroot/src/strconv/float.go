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

// FormatBool stores text representation of f in buf using format specified by
// fmt:
//	|fmt| == 'b': -ddddp±ddd, BUG: unsupported
//  |fmt| == 'e': -d.dddde±dd,
// Unused portion of the buffer is filled with spaces.
// If fmt > 0 then formatted value is left-justified and FormatFloat returns
// its length. If base < 0 then formatted value is right-justified and
// FormatFloat returns offset to its first char.

func FormatFloat(buf []byte, f float64, fmt, prec int) int {
	right := fmt < 0
	if right {
		fmt = -fmt
	}
	n := formatFloat(buf, f, fmt, prec)
	if !right {
		bytes.Fill(buf[n:], ' ')
		return n
	}
	m := len(buf) - n
	copy(buf[m:], buf[:n])
	bytes.Fill(buf[:m], ' ')
	return m
}

func formatFloat(buf []byte, f float64, fmt, prec int) int {
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
	switch fmt {
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
		g.Init(f)
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
	panic("strconv: bad fmt")
}
