package eve

// Writer allows to write data to the EVE memory.
type Writer struct {
	d *Driver
}

func (d *Driver) writer(addr int) Writer {
	d.buf = d.buf[:3]
	d.buf[0] = 1<<7 | byte(addr>>16)
	d.buf[1] = byte(addr >> 8)
	d.buf[2] = byte(addr)
	return Writer{d}
}

// W starts a write transaction to the EVE memory at address addr. It returns
// Writer that proviedes set of methods for buffered writes. If special addr -1
// is used W writes to RAM_DL and waits for INT_SWAP before sending first data
// from internal buffer.
func (d *Driver) W(addr int) Writer {
	d.end()
	d.cmdStart = -1
	if addr == -1 {
		addr = d.mmap.ramdl
		d.waitSwap = true
	}
	checkAddr(addr)
	return d.writer(addr)
}

func (w Writer) wr8(b byte) {
	d := w.d
	if len(d.buf) == cap(d.buf) {
		d.flush()
	}
	n := len(d.buf)
	d.buf = d.buf[:n+1]
	d.buf[n] = b
}

func (w Writer) Write8(s ...byte) {
	if len(s) == 1 {
		w.wr8(s[0])
		return
	}
	if len(s) == 0 {
		return
	}
	d := w.d
	if len(s) >= cap(d.buf) {
		if len(d.buf) > 0 {
			d.flush()
		}
		d.dci.Write(s)
		d.n += len(s)
		return
	}
	n := len(d.buf)
	m := copy(d.buf[n:cap(d.buf)], s)
	d.buf = d.buf[:n+m]
	if m < len(s) {
		s = s[m:]
		d.flush()
		d.buf = d.buf[:len(s)]
		copy(d.buf, s)
	}
}

func (w Writer) Write(s []byte) (int, error) {
	w.Write8(s...)
	return len(s), nil
	// BUG?: Write always succeeds. Rationale: there is no case for infinite
	// write transaction so Driver.Err can be called after all writes.
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

func (w Writer) ws(s string) {
	d := w.d
	for len(s) != 0 {
		if len(d.buf) == cap(d.buf) {
			d.flush()
		}
		n := len(d.buf)
		c := copy(d.buf[n:cap(d.buf)], s)
		d.buf = d.buf[:n+c]
		s = s[c:]
	}
}

func (w Writer) WriteString(s string) (int, error) {
	w.ws(s)
	return len(s), nil
	// BUG?: WriteString always succeeds. Rationale: there is no case for
	// infinite write transaction so Driver.Err can be called after all writes.
}

// Close flushesh internal buffer, closes the write transaction and returns
// number of bytes written. After Close w is invalid and should not be used.
// Close is called implicitly when another read/write/command transaction is
// started.
func (w Writer) Close() int {
	return w.d.end()
}
