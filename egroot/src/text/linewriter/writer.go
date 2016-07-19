// Package linewriter implements a write filter that expect "\n" or "\r" ended
// lines as its input.
package linewriter

import (
	"bufio"
	"io"
)

// Writer allows to write input text line by line to provided output io.Writer.
// If no conversion is enabled it calls output.Write method once per line. In
// conversion mode it converts any newline to provided nl string but in this
// case it can use more than one write for one line. If bufio.Writer is used as
// output it calls its Flush method after each line written.
type Writer struct {
	out   io.Writer
	newnl []byte
}

var (
	CRLF = []byte{'\r', '\n'}
	CR   = CRLF[0:1]
	LF   = CRLF[1:2]
)

func Make(output io.Writer, newnl []byte) Writer {
	return Writer{out: output, newnl: newnl}
}

func New(output io.Writer, newnl []byte) *Writer {
	b := new(Writer)
	*b = Make(output, newnl)
	return b
}

func indexNL(buf []byte) int {
	for n, b := range buf {
		if b == '\n' || b == '\r' {
			return n
		}
	}
	return -1
}

func (w *Writer) Write(buf []byte) (n int, err error) {
	b, _ := w.out.(*bufio.Writer)
	for {
		var m int
		i := indexNL(buf)
		if i == -1 {
			m, err = w.out.Write(buf)
			n += m
			return n, err
		}
		if w.newnl != nil {
			m, err = w.out.Write(buf[:i])
			n += m
			if err == nil {
				m, err = w.out.Write(w.newnl)
				n += m
			}
		} else {
			m, err = w.out.Write(buf[:i+1])
			n += m
		}
		if b != nil && err == nil {
			err = b.Flush()
		}
		if err != nil {
			return n, err
		}
		buf = buf[i+1:]
	}
}
