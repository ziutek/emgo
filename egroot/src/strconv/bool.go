package strconv

import (
	"io"
)

// WriteBool writes text representation of b to w using format specified by
// base:
//	|base| == 1: 0 / 1,
//  |base| == 2: false / true.
// Formatted value is extended to |width| characters. If width > 0 then spaces
// are written after value, otherwise spaces (base > 0) or zeros (base < 0) are
// written before it.
func WriteBool(w io.Writer, b bool, base, width int) (int, error) {
	zeros := base < 0
	if zeros {
		base = -base
	}
	txt := "0false1true"
	switch base {
	case 1:
		if b {
			txt = txt[6:7]
		} else {
			txt = txt[:1]
		}
	case 2:
		if b {
			txt = txt[7:]
		} else {
			txt = txt[1:6]
		}
	default:
		panicBase()
	}
	return writeStringPadded(w, txt, width, zeros)
}
