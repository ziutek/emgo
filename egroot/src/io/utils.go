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

// Copy copies from src to dst until either EOF is reached on src or an error
// occurs. It returns the number of bytes copied and the first error encountered
// while copying, if any.
func Copy(dst Writer, src Reader) (written int64, err error) {
	return ioCopy(dst, src, nil)
}

// CopyBuffer works like copy but use provided buffer if neither dst nor src
// implement ReaderFrom/WriteTo interfaces.
func CopyBuffer(dst Writer, src Reader, buf []byte) (written int64, err error) {
	if buf != nil && len(buf) == 0 {
		panic("io.CopyBuffer: empty buffer")
	}
	return ioCopy(dst, src, buf)
}

func ioCopy(dst Writer, src Reader, buf []byte) (written int64, err error) {
	if wt, ok := src.(WriterTo); ok {
		return wt.WriteTo(dst)
	}
	if rf, ok := dst.(ReaderFrom); ok {
		return rf.ReadFrom(src)
	}
	if buf != nil {
		return copyBuffer(dst, src, buf)
	}
	var lbuf [256]byte
	return copyBuffer(dst, src, lbuf[:])
}

func copyBuffer(dst Writer, src Reader, buf []byte) (written int64, err error) {
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != EOF {
				err = er
			}
			break
		}
	}
	return written, err
}
