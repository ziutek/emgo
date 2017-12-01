package eve

// Writer allows to write data to the EVE memory.
type Writer struct {
	d *Driver
}

func (w Writer) wr32(u uint32) {
	d := w.d
	if len(d.buf) == cap(d.buf) {
		d.flush()
	}
	n := len(d.buf)
	d.buf = d.buf[:n+1]
	d.buf[n] = u
}

func (w Writer) Write32(s ...uint32) {
	if d := w.d; len(s) > cap(d.buf) {
		d.flush()
		d.dci.Write32(s)
		d.n += len(s)
		return
	}
	for _, u := range s {
		w.wr32(u)
	}
}
