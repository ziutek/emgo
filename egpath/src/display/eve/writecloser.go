package eve

type WriteCloser Driver

func (w *WriteCloser) flush() {
	w.dci.Write(w.buf)
	w.buf = w.buf[:0]
}

func (w *WriteCloser) Write(s []byte) (int, error) {
	for p := 0; p < len(s); {
		if len(w.buf) == cap(w.buf) {
			w.flush()
		}
		m := len(w.buf)
		n := copy(w.buf[m:], s[p:])
		w.buf = w.buf[:m+n]
		p += n
	}
	if err := w.dci.Err(); err != nil {
		return 0, err
	}
	return len(s), nil
}

func (w *WriteCloser) WriteString(s string) (int, error) {
	for p := 0; p < len(s); {
		if len(w.buf) == cap(w.buf) {
			w.flush()
		}
		m := len(w.buf)
		n := copy(w.buf[m:], s[p:])
		w.buf = w.buf[:m+n]
		p += n
	}
	if err := w.dci.Err(); err != nil {
		return 0, err
	}
	return len(s), nil
}

func (w *WriteCloser) Close() error {
	if len(w.buf) > 0 {
		w.flush()
	}
	w.dci.End()
	return w.dci.Err()
}
