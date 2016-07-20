package io

type stringWriter interface {
	WriteString(string) (int, error)
}

func WriteString(w Writer, s string) (int, error) {
	if sw, ok := w.(stringWriter); ok {
		return sw.WriteString(s)
	}
	return w.Write([]byte(s))
}

type discard byte

func (_ discard) Write(b []byte) (int, error)       { return len(b), nil }
func (_ discard) WriteString(s string) (int, error) { return len(s), nil }

// Discard implements Writer and StringWriter and simply
// discards all data written to it.
const Discard discard = 0

// ReadFull reads exactly len(buf) bytes from r into buf. It returns the
// number of bytes copied and an error if fewer bytes were read. The error
// is EOF only if no bytes were read. If an EOF happens after reading some
// but not all the bytes, ReadFull returns ErrUnexpectedEOF.
func ReadFull(r Reader, buf []byte) (n int, err error) {
	for n < len(buf) {
		m, err := r.Read(buf[n:])
		n += m
		if err != nil {
			if err == EOF && n < len(buf) {
				err = ErrUnexpectedEOF
			}
			return n, err
		}
	}
	return n, nil
}
