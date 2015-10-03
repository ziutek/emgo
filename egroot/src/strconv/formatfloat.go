package strconv

import (
	"bytes"
	"io"
	"math"
)

func specials(w io.Writer, f float64, width int) (int, error) {
	txt := "Nan-Inf+Inf"
	switch {
	case math.IsNaN(f):
		txt = txt[:3]
	case math.IsInf(f, 0):
		if f < 0 {
			txt = txt[3:7]
		} else {
			txt = txt[7:]
		}
	default:
		return 0, nil
	}
	return writeStringPadded(w, txt, width, false)
}

func formatExp(buf []byte, exp int) int {
	neg := exp < 0
	if neg {
		exp = -exp
	}
	n := formatUint32(buf, uint32(exp), 10) - 1
	if exp < 10 {
		buf[n] = '0'
		n--
	}
	if neg {
		buf[n] = '-'
	} else {
		buf[n] = '+'
	}
	return n
}

const maxprec = 19

func round(buf []byte, n int, g *grisu) int {
	if n >= maxprec {
		return 0
	}
sw:
	switch d := g.NextDigit(); {
	case d < 5:
		return 0
	case d > 5:
		break
	default:
		for i := maxprec - n - 1; i > 0; i-- {
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

// WriteFloat writes text representation of f using format specified by fmt:
//	|fmt| == 'b': -ddddp±ddd,
//  |fmt| == 'e': -d.dddde±dd,
func WriteFloat(w io.Writer, f float64, fmt, width, prec, bitsize int) (int, error) {
	if n, err := specials(w, f, width); n > 0 {
		return n, err
	}
	neg := math.Signbit(f)
	if neg {
		f = -f
	}
	var (
		frac uint64
		exp  int
	)
	switch bitsize {
	case 32:
		fr, ex := math.Frexpi32(float32(f))
		frac = uint64(fr)
		exp = ex
	case 64:
		fr, ex := math.Frexpi(f)
		frac = uint64(fr)
		exp = ex
	default:
		panic("strconv: illegal bitsize")
	}
	zeros := fmt < 0
	if zeros {
		fmt = -fmt
	}
	if fmt == 'b' {
		var buf [1 + 16 + 1 + 1 + 4]byte
		n := formatExp(buf[:], exp) - 1
		buf[n] = 'p'
		n = formatUint64(buf[:n], frac, 10)
		if neg {
			n--
			buf[n] = '-'
		}
		return writePadded(w, buf[n:], width, zeros)
	}
	var ft int
	switch fmt {
	case 'f', 'F':
		ft = 'f'
		fmt--
		prec++ // Convert prec to total number of digits.
	case 'e', 'E':
		ft = 'e'
		prec++
	case 'g', 'G':
		ft = 'g'
		fmt -= 2
	default:
		panic("strconv: bad fmt")
	}
	if prec < 1 {
		prec = 1
	} else if prec > maxprec {
		// BUG: Allow prec > maxprec by simply write zeros after decimal point.
		prec = maxprec
	}
	var (
		buf [1 + maxprec + 1 + 1 + 4]byte // len(buf) == maxprec + 7
		n   int
	)
	if neg {
		buf[0] = '-'
		n++
	}
	if frac == 0 {
		buf[n] = '0'
		n++
		if prec > 1 {
			buf[n] = '.'
			bytes.Fill(buf[n+1:n+prec], '0')
			n += prec
		}
		if ft == 'e' {
			n += copy(buf[n:], "e+00")
		}
	} else {
		var g grisu
		g.Init(frac, exp)
		n++ // Save place for dot.
		for i := 0; i < prec; i++ {
			buf[n+i] = byte('0' + g.NextDigit())
		}
		exp10 := g.Exp10() + prec - 1
		exp10 += round(buf[n:], prec, &g)
		if ft == 'g' && exp10 < prec && exp10 >= -4 {
			if exp10 >= 0 {
				dot := n + exp10
				copy(buf[n-1:], buf[n:dot+1])
				if exp10 == prec-1 {
					n += exp10
				} else {
					buf[dot] = '.'
					n += prec + 1
				}
			} else {
				firstd := n - exp10
				copy(buf[firstd:], buf[n:n+prec])
				buf[n-1] = '0'
				buf[n] = '.'
				bytes.Fill(buf[n+1:firstd], '0')
				n = firstd + prec
			}
		} else {
			buf[n-1] = buf[n]
			if prec > 1 {
				buf[n] = '.'
			} else {
				n-- // No dot.
			}
			n += prec
			buf[n] = byte(fmt)
			n++
			m := formatExp(buf[n:], exp10)
			n += copy(buf[n:], buf[n+m:])
		}
	}
	return writePadded(w, buf[:n], width, zeros)
}
