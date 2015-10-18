package fmt

import (
	"io"
	"stack"
	"strings"
	"unsafe"
)

func findVerb(f string) (int, string, byte) {
	start := strings.IndexByte(f, '%')
	if start == -1 {
		return len(f), "", 0
	}
	f = f[start+1:]
	for k := 0; k < len(f); k++ {
		c := f[k]
		if c >= '0' && c <= '9' {
			continue
		}
		switch c {
		case '+', '-', ' ', '*', '#', '.':
			continue
		}
		return start, f[:k], c
	}
	return start, "", 0
}

func Fprintf(w io.Writer, f string, a ...interface{}) (int, error) {
	p, ok := w.(printer)
	if !ok {
		ptr := unsafe.Pointer(stack.Alloc(1, unsafe.Sizeof(*p.writer)))
		p.writer = (*writer)(ptr)
		p.writer.w = w
	}
	var m int
	for {
		start, flags, verb := findVerb(f)
		p.WriteString(f[:start])
		if p.err != nil {
			return p.n, p.err
		}
		if start == len(f) {
			break
		}
		switch verb {
		case '%':
			p.WriteByte('%')
		case 0:
			// Unfinished format.
			p.Ferr(0, "UNFINISHED", nil)
		default:
			if m < len(a) {
				p.parse(flags)
				p.format(verb, a[m])
			}
			m++
		}
		if p.err != nil {
			return p.n, p.err
		}
		if m > len(a) {
			p.Ferr(verb, "MISSING", nil)
			if p.err != nil {
				return p.n, p.err
			}
		}
		f = f[start+2+len(flags):]
	}
	for ; m < len(a); m++ {
		p.Ferr(0, "EXTRA ", a[m])
		if p.err != nil {
			break
		}
	}
	return p.n, p.err
}