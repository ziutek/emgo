package strings

import (
	"io"
)

type Reader struct {
	s string
}

func MakeReader(s string) Reader {
	return Reader{s}
}

func (r *Reader) Read(buf []byte) (n int, err error) {
	n = copy(buf, r.s)
	r.s = r.s[n:]
	if len(r.s) == 0 {
		err = io.EOF
	}
	return
}
