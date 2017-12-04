package eve

// Writer allows to write data to the EVE memory.
type Writer struct {
	d *Driver
}

func (w Writer) align32() {
	d := w.d
	m := (len(d.buf) + d.n) & 3
	if m == 0 {
		return
	}
	m = len(d.buf) + 4 - m
	if m > cap(d.buf) {
		m -= len(d.buf)
		d.flush()
	}
	d.buf = d.buf[:m]
}

func (w Writer) wr32(u uint32) {
	d := w.d
	if len(d.buf)+4 > cap(d.buf) {
		d.flush()
	}
	n := len(d.buf)
	d.buf = d.buf[:n+4]
	d.buf[n] = byte(u)
	d.buf[n+1] = byte(u >> 8)
	d.buf[n+2] = byte(u >> 16)
	d.buf[n+3] = byte(u >> 24)
}

func (w Writer) aw32(u uint32) {
	w.align32()
	w.wr32(u)
}

func (w Writer) Write32(s ...uint32) {
	w.align32()
	for _, u := range s {
		w.wr32(u)
	}
}

func (w Writer) Write(s []byte) (int, error) {
	d := w.d
	if len(s) >= cap(d.buf) {
		if len(d.buf) > 0 {
			d.flush()
		}
		d.dci.Write(s)
		d.n += len(s)
	} else {
		n := copy(d.buf[:cap(d.buf)], s)
		d.buf = d.buf[:len(d.buf)+n]
		if n < len(s) {
			d.flush()
			d.buf = d.buf[:len(s)-n]
			copy(d.buf, s[n:])
		}
	}
	// Write always succeeds. There is no case for infinite write transaction
	// so Driver.Err can be called after all writes.
	return len(s), nil
}

func (w Writer) WriteString(s string) (int, error) {
	d := w.d
	n := len(s)
	for len(s) != 0 {
		if len(d.buf) == cap(d.buf) {
			d.flush()
		}
		m := len(d.buf)
		c := copy(d.buf[m:cap(d.buf)], s)
		d.buf = d.buf[:m+c]
		s = s[c:]
	}
	return n, nil
}

// Close closes the write transaction and returns number of bytes written.
func (w Writer) Close() int {
	return w.d.end()
}
