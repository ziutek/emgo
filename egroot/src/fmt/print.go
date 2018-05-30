package fmt

import (
	"io"
)

func Fprint(w io.Writer, a ...interface{}) (int, error) {
	neww := writer{w: w}
	p, ok := w.(printer)
	if !ok {
		p.writer = &neww
	}
	p.parse("")
	for _, v := range a {
		p.format('v', v)
		if p.err != nil {
			break
		}
	}
	return p.n, p.err
}

//emgo:noinline
func Fprintln(w io.Writer, a ...interface{}) (int, error) {
	neww := writer{w: w}
	p, ok := w.(printer)
	if !ok {
		p.writer = &neww
	}
	p.parse("")
	for i, v := range a {
		if i > 0 {
			p.WriteByte(' ')
			if p.err != nil {
				return p.n, p.err
			}
		}
		p.format('v', v)
		if p.err != nil {
			return p.n, p.err
		}
	}
	p.WriteByte('\n')
	return p.n, p.err
}

var DefaultWriter io.Writer

func Print(a ...interface{}) (int, error) {
	return Fprint(DefaultWriter, a...)
}

func Println(a ...interface{}) (int, error) {
	return Fprintln(DefaultWriter, a...)
}

func Printf(f string, a ...interface{}) (int, error) {
	return Fprintf(DefaultWriter, f, a...)
}
