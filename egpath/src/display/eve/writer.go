package eve

// Writer allows to write data to thr EVE.
type Writer Driver

func (w *Writer) flush() {
	w.n += len(w.buf)
	w.dci.Write32(w.buf)
	w.buf = w.buf[:0]
}

func (w *Writer) WriteWord32(u uint32) {
	if len(w.buf) == cap(w.buf) {
		w.flush()
	}
	n := len(w.buf)
	w.buf = w.buf[:n+1]
	w.buf[n] = u
}
