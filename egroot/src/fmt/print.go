package fmt

import (
	"io"
	"stack"
	"unsafe"
)

func Fprint(w io.Writer, a ...interface{}) (int, error) {
	p, ok := w.(printer)
	if !ok {
		ptr := unsafe.Pointer(stack.Alloc(1, unsafe.Sizeof(*p.writer)))
		p.writer = (*writer)(ptr)
		p.writer.w = w
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

func Fprintln(w io.Writer, a ...interface{}) (int, error) {
	p, ok := w.(printer)
	if !ok {
		ptr := unsafe.Pointer(stack.Alloc(1, unsafe.Sizeof(*p.writer)))
		p.writer = (*writer)(ptr)
		p.writer.w = w
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
