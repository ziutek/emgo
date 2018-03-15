package strconv

import (
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
	return WriteString(w, txt, width, ' ')
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

func round(buf []byte, g *grisu) int {
	n := len(buf)
sw:
	switch d := g.NextDigit(); {
	case d < 5:
		return 0
	case d > 5:
		break
	default:
		for i := maxDigits - n - 1; i > 0; i-- {
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

const (
	flagF = 1 << iota
	flagE
	flagG
	flagNeg
	flagLeft
)

func padBefore(w io.Writer, num, flags int, pad rune) (int, error) {
	var n int
	if flags&flagNeg != 0 {
		num--
	}
	if num > 0 && pad != '0' {
		m, err := writeRuneN(w, pad, num)
		n += m
		if err != nil {
			return n, err
		}
	}
	if flags&flagNeg != 0 {
		m, err := w.Write([]byte{'-'})
		n += m
		if err != nil {
			return n, err
		}
	}
	if num > 0 && pad == '0' {
		m, err := writeRuneN(w, pad, num)
		n += m
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

func writeZero(w io.Writer, width, prec, flags, e int, pad rune) (int, error) {
	var padn int
	if flags&flagLeft == 0 {
		padn = width - 1
		if prec > 0 {
			padn -= prec + 1
		}
		if flags&flagE != 0 {
			padn -= 4
		}
	}
	n, err := padBefore(w, padn, flags, pad)
	if err != nil {
		return n, err
	}
	m, err := w.Write([]byte{'0'})
	n += m
	if err != nil {
		return n, err
	}
	if prec > 0 {
		m, err = w.Write([]byte{'.'})
		n += m
		if err != nil {
			return n, err
		}
		m, err = writeRuneN(w, '0', prec)
		n += m
		if err != nil {
			return n, err
		}
	}
	if flags&flagE != 0 {
		m, err = w.Write([]byte{byte(e), '+', '0', '0'})
		n += m
		if err != nil {
			return n, err
		}
	}
	if padn := width - n; padn > 0 {
		if pad == '0' {
			pad = ' '
		}
		m, err = writeRuneN(w, pad, padn)
		n += m
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

// WriteFloat writes text representation of f using format specified by fmt:
//	'b': -ddddp±ddd,
//	'e': -d.dddde±dd, prec sets the number of digits after decimal,
//	'f': -ddddd.dddd, prec sets the number of digits after decimal,
//  'g': shortest from 'e' or 'f', prec sets the number of significant digits.
// For description of width and pad see WriteInt.
func WriteFloat(w io.Writer, f float64, fmt, prec, bitsize, width int, pad rune) (int, error) {
	if n, err := specials(w, f, width); n > 0 {
		return n, err
	}
	var flags int
	if math.Signbit(f) {
		flags |= flagNeg
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
	if fmt == 'b' {
		var buf [1 + 16 + 1 + 1 + 4]byte
		n := formatExp(buf[:], exp) - 1
		buf[n] = 'p'
		n = formatUint64(buf[:n], frac, 10)
		if flags&flagNeg != 0 {
			n--
			buf[n] = '-'
		}
		return writePadded(w, buf[n:], width, pad)
	}
	switch fmt {
	case 'f', 'F':
		flags |= flagF
	case 'e', 'E':
		flags |= flagE
	case 'g', 'G':
		flags |= flagG
		fmt -= 2
		prec-- // Convert prec to number of digits after decimal point.
	default:
		panic("strconv: bad fmt")
	}
	if width < 0 {
		flags |= flagLeft
		width = -width
	}
	if prec < 0 {
		prec = 0
	}
	if frac == 0 {
		return writeZero(w, width, prec, flags, fmt, pad)
	}
	var g grisu
	dig, exp := g.Init(frac, exp)
	var sigd int
	if flags&flagF != 0 {
		sigd = exp + 1 + prec
		if sigd < 0 {
			return writeZero(w, width, prec, flags, 0, pad)
		}
	} else {
		sigd = prec + 1
	}
	if sigd > maxDigits {
		sigd = maxDigits
	}
	var arr [1 + maxDigits]byte
	buf := arr[1:]
	buf[0] = byte('0' + dig)
	for i := 1; i < sigd; i++ {
		buf[i] = byte('0' + g.NextDigit())
	}
	if sigd < maxDigits {
		exp += round(buf[:sigd], &g)
	}
	if flags&flagG != 0 && exp <= prec && exp >= -4 {
		flags |= flagF
		prec -= exp
	}
	var padn int
	if flags&flagLeft == 0 {
		length := prec
		if prec > 0 {
			length++ // Dot
		}
		if flags&flagF != 0 && exp > 0 {
			length += exp + 1
		} else {
			length++ // Digit before dot.
		}
		if flags&flagF == 0 {
			length += 4
			if exp >= 100 || exp <= -100 {
				length++
			}
		}
		padn = width - length
	}
	n, err := padBefore(w, padn, flags, pad)
	if err != nil {
		return n, err
	}
	var m int
	if flags&flagF != 0 {
		todot := exp + 1
		if todot <= 0 {
			m, err = w.Write([]byte{'0', '.'})
			n += m
			if err != nil {
				return n, err
			}
			m, err = writeRuneN(w, '0', -todot)
			prec -= m
			n += m
			if err != nil {
				return n, err
			}
			m, err = w.Write(buf[:sigd])
			prec -= m
			n += m
			if err != nil {
				return n, err
			}
		} else {
			if padn = todot - sigd; padn <= 0 {
				m, err = w.Write(buf[:todot])
				n += m
				if err != nil {
					return n, err
				}
			} else {
				m, err = w.Write(buf[:sigd])
				n += m
				if err != nil {
					return n, err
				}
				m, err = writeRuneN(w, '0', padn)
				n += m
				if err != nil {
					return n, err
				}
			}
			if prec > 0 {
				m, err = w.Write([]byte{'.'})
				n += m
				if err != nil {
					return n, err
				}
				if padn < 0 {
					m, err = w.Write(buf[todot:sigd])
					prec -= m
					n += m
					if err != nil {
						return n, err
					}
				}
			}
		}
	} else {
		arr[0] = arr[1]
		if prec > 0 {
			prec -= sigd - 1
			arr[1] = '.'
			sigd++
		}
		m, err = w.Write(arr[:sigd])
		n += m
		if err != nil {
			return n, err
		}
	}
	if prec > 0 {
		m, err = writeRuneN(w, '0', prec)
		n += m
		if err != nil {
			return n, err
		}
	}
	if flags&flagF == 0 {
		m = formatExp(buf, exp) - 1
		buf[m] = byte(fmt)
		m, err = w.Write(buf[m:])
		n += m
		if err != nil {
			return n, err
		}
	}
	if padn = width - n; padn > 0 {
		if pad == '0' {
			pad = ' '
		}
		m, err = writeRuneN(w, pad, padn)
		n += m
	}
	return n, err
}
