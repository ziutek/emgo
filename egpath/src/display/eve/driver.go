package eve

type Driver struct {
	dci DCI
	buf []byte
}

func NewDriver(dci DCI, n int) *Driver {
	d := new(Driver)
	d.dci = dci
	d.buf = make([]byte, 0, n)
	return d
}

func (d *Driver) Err() error {
	return d.dci.Err()
}

type HostCmd byte

// Cmd invokes host command. Arg is command argument. It must be zero for
// commands that do not require arguments.
func (d *Driver) Cmd(cmd HostCmd, arg byte) {
	d.buf = d.buf[:3]
	d.buf[0] = byte(cmd)
	d.buf[1] = arg
	d.buf[2] = 0
	d.dci.Write(d.buf)
	d.dci.End()
}

func checkAddr(addr int) {
	if uint(addr) >= 1<<22 {
		panic("eve: bad addr")
	}
}

func (d *Driver) StartWrite(addr int) *WriteCloser {
	checkAddr(addr)
	d.buf = d.buf[:3]
	d.buf[0] = 1<<7 | byte(addr>>16)
	d.buf[1] = byte(addr >> 8)
	d.buf[2] = byte(addr)
	return (*WriteCloser)(d)
}
