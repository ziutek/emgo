package strconv

import (
	"io"
)

// WriteBool writes text representation of b to w using format specified by fmt:
//	false / true for 't', 'v',
//	0 / 1 for 'd', 'x', 'X'.
// If width > 0 then written value is right-justified, otherwise it is
// left-justified.
func WriteBool(w io.Writer, b bool, fmt byte, width int, pad rune) (int, error) {
	txt := "0false1true"
	switch fmt {
	case 'd', 'x', 'X':
		if b {
			txt = txt[6:7]
		} else {
			txt = txt[:1]
		}
	case 't', 'v':
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
