package eve

// Writer allows to write data to the EVE memory.
type Writer struct {
	d *Driver
}

func (w Writer) start(addr int) {
	d := w.d
	d.buf = d.buf[:3]
	d.buf[0] = 1<<7 | byte(addr>>16)
	d.buf[1] = byte(addr >> 8)
	d.buf[2] = byte(addr)
	d.state |= stateOpen
}

// W opens a write transaction to the EVE memory at address addr. It returns
// Writer that proviedes set of methods for buffered writes. If special addr -1
// is used W waits for INT_SWAP and starts write to RAM_DL.
func (d *Driver) W(addr int) Writer {
	d.end()
	if addr == -1 {
		d.wait(INT_SWAP)
		addr = d.mmap.ramdl
	} else {
		checkAddr(addr)
	}
	d.addr = addr
	d.state = d.state&^3 | stateWrite
	w := Writer{d}
	w.start(addr)
	return w
}

func (w Writer) restart(n int) {
	d := w.d
	if d.state&stateOpen == 0 {
		w.start(d.Addr())
	}
	d.addr += n
}

// Flush writes all data from internal buffer, closes current transaction and
// returns the current write address.
func (w Writer) Flush() int {
	w.d.end()
	return w.d.Addr()
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
	if len(s) == 0 {
		return
	}
	w.restart(len(s))
	if len(s) == 1 {
		w.wr8(s[0])
		return
	}
	d := w.d
	if len(s) >= cap(d.buf) {
		if len(d.buf) > 0 {
			d.flush()
		}
		d.dci.Write(s) // Write long data directly.
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
	n := d.addr & 3
	if n == 0 {
		return
	}
	n = 4 - n
	d.addr += n
	n += len(d.buf)
	if n > cap(d.buf) {
		n -= len(d.buf)
		d.flush()
	}
	d.buf = d.buf[:n]
}

// Align32 writes random bytes to align the current write address to 32 bit.
func (w Writer) Align32() {
	w.restart(0)
	w.align32()
}

func (w Writer) Write32(s ...uint32) {
	w.restart(4 * len(s))
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
	w.restart(len(s))
	w.ws(s)
	return len(s), nil
	// BUG?: WriteString always succeeds. Rationale: there is no case for
	// infinite write transaction so Driver.Err can be called after all writes.
}
