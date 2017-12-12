package eve

// Driver uses DCI to communicate with EVE graphics controller. Commands/data
// are received/sent using DCI read/write transactions. R,W, DL, GE methods
// starts new transaction and leaves it in open state. Any subsequent
// transaction implicitly closes current open transaction.
type Driver struct {
	dci           DCI
	buf           []byte
	n             int
	mmap          *mmap
	width, height uint16
	intflags      byte
}

// NewDriver returns new driver to the EVE graphics controller accessed via dci.
// N sets the capacity of the internal buffer (bytes, must be >= 4).
func NewDriver(dci DCI, n int) *Driver {
	d := new(Driver)
	d.dci = dci
	d.buf = make([]byte, 0, n)
	d.n = -3
	return d
}

func (d *Driver) Width() int {
	return int(d.width)
}

func (d *Driver) Height() int {
	return int(d.height)
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
// case of commands that do not require parameters.
func (d *Driver) HostCmd(cmd HostCmd, param byte) {
	d.end()
	d.dci.Write([]byte{byte(cmd), param, 0})
	d.dci.End()
}

// WriteByte writes byte to the EVE memory at address addr.
func (d *Driver) WriteByte(addr int, val byte) {
	d.end()
	d.dci.Write([]byte{
		1<<7 | byte(addr>>16), byte(addr >> 8), byte(addr),
		val,
	})
	d.dci.End()
}

// WriteUint16 writes 16-bit word to the EVE memory at address addr.
func (d *Driver) WriteUint16(addr int, val uint16) {
	d.end()
	d.dci.Write([]byte{
		1<<7 | byte(addr>>16), byte(addr >> 8), byte(addr),
		byte(val), byte(val >> 8),
	})
	d.dci.End()
}

// WriteUint32 writes 32-bit word to the EVE memory at address addr.
func (d *Driver) WriteUint32(addr int, val uint32) {
	d.end()
	d.dci.Write([]byte{
		1<<7 | byte(addr>>16), byte(addr >> 8), byte(addr),
		byte(val), byte(val >> 8), byte(val >> 16), byte(val >> 24),
	})
	d.dci.End()
}

// WriteInt writes signed 32-bit word to the EVE memory at address addr.
func (d *Driver) WriteInt(addr int, val int) {
	d.WriteUint32(addr, uint32(val))
}

// W starts a write transaction to the EVE memory at address addr. It returns
// Writer that proviedes set of methods for buffered writes. Any other Drivers's
// method flushes internal buffer and finishes the write transaction started by
// W. After that, the returned Writer is invalid and should not be used.
func (d *Driver) W(addr int) Writer {
	checkAddr(addr)
	d.end()
	d.buf = d.buf[:3]
	d.buf[0] = 1<<7 | byte(addr>>16)
	d.buf[1] = byte(addr >> 8)
	d.buf[2] = byte(addr)
	return Writer{d}
}

// R starts a read transaction from the EVE memory at address addr. It
// returns Reader that provides set of reading methods. Any other Driver's
// method finish the read transaction started by R. After that, the returned
// Reader is invalid and should not be used.
func (d *Driver) R(addr int) Reader {
	checkAddr(addr)
	d.end()
	d.dci.Write([]byte{byte(addr >> 16), byte(addr >> 8), byte(addr)})
	d.n = 0
	d.dci.Read([]byte{0}) // Read dummy byte (input mode required by QSPI ).
	return Reader{d}
}

// ReadByte reads byte from EVE memory at address addr.
func (d *Driver) ReadByte(addr int) byte {
	r := d.R(addr)
	val := r.ReadByte()
	r.Close()
	return val
}

// ReadUint16 reads 16-bit word from EVE memory at address addr.
func (d *Driver) ReadUint16(addr int) uint16 {
	r := d.R(addr)
	val := r.ReadUint16()
	r.Close()
	return val
}

// ReadUint32 reads 32-bit word from EVE memory at address addr.
func (d *Driver) ReadUint32(addr int) uint32 {
	r := d.R(addr)
	val := r.ReadUint32()
	r.Close()
	return val
}

// ReadUint32 reads signed 32-bit word from EVE memory at address addr.
func (d *Driver) ReadInt(addr int) int {
	return int(int32(d.ReadUint32(addr)))
}

// Err returns and clears the internal error status.
func (d *Driver) Err(clear bool) error {
	d.end()
	return d.dci.Err(clear)
}

// DL wraps W to return Display List writer. See W for more information.
func (d *Driver) DL(addr int) DL {
	return DL{d.W(addr)}
}

// GE wraps DL to retun Graphics Engine command writer. See DL for more
// information.
func (d *Driver) GE(addr int) GE {
	return GE{d.DL(addr)}
}

// IRQ returns channel that can be used to wait for IRQ.
func (d *Driver) IRQ() <-chan struct{} {
	return d.dci.IRQ()
}
