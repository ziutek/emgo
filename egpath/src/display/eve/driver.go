package eve

type Driver struct {
	dci DCI
	buf []uint32
	n   int
}

// NewDriver returns new driver to the EVE graphics controller accessed via dci.
// N sets capacity of the internal buffer (in 32-bit words, must be >= 1).
func NewDriver(dci DCI, n int) *Driver {
	d := new(Driver)
	d.dci = dci
	d.buf = make([]uint32, 0, n)
	return d
}

func checkAddr(addr int) {
	if uint(addr) >= 1<<22 {
		panic("eve: bad addr")
	}
}

func (d *Driver) flush() {
	d.n += len(d.buf)
	d.dci.Write32(d.buf)
	d.buf = d.buf[:0]
}

// End should be used to ensure that the previous Reader/Writer finished. It
// returns the number of bytes written to the EVE memory. End is implicitly
// called at the beginning of Cmd, Reader, Writer and Err methods.
func (d *Driver) End() int {
	if len(d.buf) > 0 {
		d.flush()
	}
	n := d.n
	if n > 0 {
		n = (n - 1) * 4
		d.n = 0
		d.dci.End()
	}
	return n
}

type HostCmd byte

// Cmd invokes host command. Param is a command parameter. It must be zero in
// case of commands that do not require parameters. Cmd is not buffered by the
// Driver: does not require calling End after it.
func (d *Driver) Cmd(cmd HostCmd, param byte) {
	d.End()
	d.dci.Write32([]uint32{uint32(cmd)<<16 | uint32(param)<<8})
	d.dci.End()
}

// Writer starts writing to the EVE memory at the address addr.
func (d *Driver) Writer(addr int) Writer {
	checkAddr(addr)
	d.End()
	d.buf = d.buf[:1]
	d.buf[0] = 1<<23 | uint32(addr)
	return Writer{d}
}

// Reader starts reading from the EVE memory at the address addr.
func (d *Driver) Reader(addr int) Reader {
	checkAddr(addr)
	d.End()
	d.dci.Write32([]uint32{uint32(addr)})
	d.dci.Read([]byte{0}) // Read dummy byte (switch QSPI to input mode).
	d.n = 1
	return Reader{d}
}

// Err returns and clears the internal error status.
func (d *Driver) Err(clear bool) error {
	d.End()
	return d.dci.Err(clear)
}

func (d *Driver) DL(addr int) DL {
	return DL{d.Writer(addr)}
}

func (d *Driver) GE(addr int) GE {
	return GE{d.DL(addr)}
}
