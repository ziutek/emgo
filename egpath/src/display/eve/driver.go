package eve

type Driver struct {
	dci DCI
	buf []uint32
}

// n must be >= 1
func NewDriver(dci DCI, n int) *Driver {
	d := new(Driver)
	d.dci = dci
	d.buf = make([]uint32, 0, n)
	return d
}

func (d *Driver) ClearErr() error {
	return d.dci.Err(true)
}

// End should be used to ensure that previous command finished.
func (d *Driver) End() error {
	if len(d.buf) > 0 {
		d.dci.Write32(d.buf)
		d.buf = d.buf[:0]
	}
	d.dci.End()
	return d.dci.Err(false)
}

type HostCmd byte

// Cmd invokes host command. Arg is command argument. It must be zero for
// commands that do not require arguments.
func (d *Driver) Cmd(cmd HostCmd, arg byte) {
	d.End()
	d.buf = d.buf[:1]
	d.buf[0] = uint32(cmd)<<16 | uint32(arg)<<8
	d.dci.Write32(d.buf)
}

func checkAddr(addr int) {
	if uint(addr) >= 1<<22 {
		panic("eve: bad addr")
	}
}

func (d *Driver) Writer(addr int) *Writer {
	checkAddr(addr)
	d.End()
	d.buf = d.buf[:1]
	d.buf[0] = 1<<7 | uint32(addr)
	return (*Writer)(d)
}

func (d *Driver) Reader(addr int) *Reader {
	checkAddr(addr)
	d.End()
	d.dci.Write32([]uint32{uint32(addr)})
	d.dci.Read([]byte{0}) // Read dummy byte (switch QSPI to input mode).
	return (*Reader)(d)
}
