package eve

// Writer allows to write data to the EVE memory.
type Writer struct {
	d *Driver
}

func (w Writer) Write32(s ...uint32) {
	d := w.d
	for _, u := range s {
		if len(d.buf) == cap(d.buf) {
			d.flush()
		}
		n := len(d.buf)
		d.buf = d.buf[:n+1]
		d.buf[n] = u
	}
}
