package strconv

import (
	"io"
)

// WriteBool writes text representation of b to w using format specified by
// base:
//	fmt = '1': 0 / 1,
//	fmt = 't': false / true.
// If width > 0 then written value is right-justified, otherwise it is
// left-justified.
func WriteBool(w io.Writer, b bool, fmt, width int, pad rune) (int, error) {
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
	return WriteString(w, txt, width, pad)
}
