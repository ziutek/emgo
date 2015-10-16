package strconv

import (
	"io"
)

// WriteBool writes text representation of b to w using format specified by
// base:
//	|fmt| == '1': 0 / 1,
//  |fmt| == 't': false / true.
// Formatted value is extended to |width| characters. If width > 0 then spaces
// are written after value, otherwise spaces (fmt > 0) or zeros (base < 0) are
// written before it.
func WriteBool(w io.Writer, b bool, fmt, width int) (int, error) {
	zeros := fmt < 0
	if zeros {
		fmt = -fmt
	}
	txt := "0false1true"
	switch fmt {
	case '1':
		if b {
			txt = txt[6:7]
		} else {
			txt = txt[:1]
		}
	case 't':
		if b {
			txt = txt[7:]
		} else {
			txt = txt[1:6]
		}
	default:
		panicBase()
	}
	return WriteString(w, txt, width, zeros)
}
