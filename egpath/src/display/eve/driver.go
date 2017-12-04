package eve

// Driver uses DCI to communicate with EVE graphics controller. Commands/data
// are received/sent using DCI read/write transactions. Start* methods starts
// new transaction and leaves it in open state. Subsequent call of any Driver's
// method implicitly closes previously opened transaction..
type Driver struct {
	dci DCI
	buf []byte
	n   int
}

// NewDriver returns new driver to the EVE graphics controller accessed via dci.
// N sets capacity of the internal buffer (in 32-bit words, must be >= 1).
func NewDriver(dci DCI, n int) *Driver {
	d := new(Driver)
	d.dci = dci
	d.buf = make([]byte, 0, n*4)
	d.n = -3
	return d
}

func checkAddr(addr int) {
	if uint(addr) >= 1<<22 {
		panic("eve: bad addr")
	}
}

func (d *Driver) flush() {
	d.dci.Write(d.buf)
	d.n += len(d.buf)
	d.buf = d.buf[:0]
}

func (d *Driver) end() int {
	if len(d.buf) > 0 {
		d.flush()
	}
	n := d.n
	if n >= 0 {
		d.n = -3
		d.dci.End()
	}
	return n
}

type HostCmd byte

// Cmd invokes host command. Param is a command parameter. It must be zero in
// case of commands that do not require parameters. Cmd is not buffered by the
// Driver: does not require calling end after it.
func (d *Driver) Cmd(cmd HostCmd, param byte) {
	d.end()
	d.dci.Write([]byte{byte(cmd), param, 0})
	d.dci.End()
}

// Writer starts writing to the EVE memory at the address addr.
func (d *Driver) StartW(addr int) Writer {
	checkAddr(addr)
	d.end()
	d.buf = d.buf[:3]
	d.buf[0] = 1<<7 | byte(addr>>16)
	d.buf[1] = byte(addr >> 8)
	d.buf[2] = byte(addr)
	return Writer{d}
}

// Reader starts reading from the EVE memory at the address addr.
func (d *Driver) StartR(addr int) Reader {
	checkAddr(addr)
	d.end()
	d.dci.Write([]byte{byte(addr >> 16), byte(addr >> 8), byte(addr)})
	d.n = 0
	d.dci.Read([]byte{0}) // Read dummy byte (input mode required by QSPI ).
	return Reader{d}
}

// Err returns and clears the internal error status.
func (d *Driver) Err(clear bool) error {
	d.end()
	return d.dci.Err(clear)
}

func (d *Driver) StartDL(addr int) DL {
	return DL{d.StartW(addr)}
}

func (d *Driver) StartGE(addr int) GE {
	return GE{d.StartDL(addr)}
}
