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
	cmdStart      int16
	flags         byte
	waitSwap      bool
}

// NewDriver returns new driver to the EVE graphics controller accessed via dci.
// N sets the capacity of the internal buffer (bytes, must be >= 4).
func NewDriver(dci DCI, n int) *Driver {
	d := new(Driver)
	d.dci = dci
	d.buf = make([]byte, 0, n)
	d.n = -3
	d.cmdStart = -1
	return d
}

// Width returns screen width.
func (d *Driver) Width() int {
	return int(d.width)
}

// Height returns screen height.
func (d *Driver) Height() int {
	return int(d.height)
}

// RAMDL returns the start address of RAM_DL.
func (d *Driver) RAMDL() int {
	return d.mmap.ramdl
}

// RAMCMD returns the start address of RAM_CMD.
func (d *Driver) RAMCMD() int {
	return d.mmap.ramcmd
}

// IRQ returns channel that can be used to wait for IRQ.
func (d *Driver) IRQ() <-chan struct{} {
	return d.dci.IRQ()
}

func (d *Driver) writeByte(addr int, val byte) {
	d.dci.Write([]byte{
		1<<7 | byte(addr>>16), byte(addr >> 8), byte(addr),
		val,
	})
	d.dci.End()
}

func (d *Driver) writeUint16(addr int, val uint16) {
	d.dci.Write([]byte{
		1<<7 | byte(addr>>16), byte(addr >> 8), byte(addr),
		byte(val), byte(val >> 8),
	})
	d.dci.End()
}

func (d *Driver) writeUint32(addr int, val uint32) {
	d.dci.Write([]byte{
		1<<7 | byte(addr>>16), byte(addr >> 8), byte(addr),
		byte(val), byte(val >> 8), byte(val >> 16), byte(val >> 24),
	})
	d.dci.End()
}

func (d *Driver) readByte(addr int) byte {
	buf := [3]byte{byte(addr >> 16), byte(addr >> 8), byte(addr)}
	d.dci.Write(buf[:3])
	d.dci.Read(buf[:2])
	d.dci.End()
	return buf[1]
}

func (d *Driver) readUint16(addr int) uint16 {
	buf := [3]byte{byte(addr >> 16), byte(addr >> 8), byte(addr)}
	d.dci.Write(buf[:3])
	d.dci.Read(buf[:3])
	d.dci.End()
	return uint16(buf[1]) | uint16(buf[2])<<8
}

func (d *Driver) readUint32(addr int) uint32 {
	buf := [5]byte{byte(addr >> 16), byte(addr >> 8), byte(addr)}
	d.dci.Write(buf[:3])
	d.dci.Read(buf[:5])
	d.dci.End()
	return uint32(buf[1]) | uint32(buf[2])<<8 | uint32(buf[3])<<16 |
		uint32(buf[4])<<24
}

func (d *Driver) intFlags() byte {
	d.flags |= d.readByte(d.mmap.regintflags)
	return d.flags
}

func (d *Driver) clearIntFlags(mask byte) {
	d.flags |= d.readByte(d.mmap.regintflags)
	d.flags &^= mask
}

func (d *Driver) intMask() byte {
	return d.readByte(d.mmap.regintflags + ointmask)
}

func (d *Driver) setIntMask(mask byte) {
	d.writeByte(d.mmap.regintflags+ointmask, mask)
}

func (d *Driver) wait(flags byte) {
	if d.flags&flags != 0 || d.intFlags()&flags != 0 {
		return
	}
	mask := d.intMask()
	d.setIntMask(flags)
	for {
		<-d.IRQ()
		if d.intFlags()&flags != 0 {
			break
		}
	}
	d.setIntMask(mask)
}

func (d *Driver) flush() {
	if d.waitSwap {
		d.waitSwap = false
		d.wait(INT_SWAP)
	}
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
	if d.cmdStart >= 0 {
		cmdEnd := (uint32(d.cmdStart) + uint32(n+3)&^3) & 4095
		d.cmdStart = -1
		d.writeUint32(eve1_regcmdwrite, cmdEnd)
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

func checkAddr(addr int) {
	if uint(addr) >= 1<<22 {
		panic("eve: bad addr")
	}
}

// WriteByte writes byte to the EVE memory at address addr.
func (d *Driver) WriteByte(addr int, val byte) {
	d.end()
	checkAddr(addr)
	d.writeByte(addr, val)
}

// WriteUint16 writes 16-bit word to the EVE memory at address addr.
func (d *Driver) WriteUint16(addr int, val uint16) {
	d.end()
	checkAddr(addr)
	d.writeUint16(addr, val)
}

// WriteUint32 writes 32-bit word to the EVE memory at address addr.
func (d *Driver) WriteUint32(addr int, val uint32) {
	d.end()
	checkAddr(addr)
	d.writeUint32(addr, val)
}

// WriteInt writes signed 32-bit word to the EVE memory at address addr.
func (d *Driver) WriteInt(addr int, val int) {
	d.WriteUint32(addr, uint32(val))
}

// ReadByte reads byte from EVE memory at address addr.
func (d *Driver) ReadByte(addr int) byte {
	d.end()
	checkAddr(addr)
	return d.readByte(addr)
}

// ReadUint16 reads 16-bit word from EVE memory at address addr.
func (d *Driver) ReadUint16(addr int) uint16 {
	d.end()
	checkAddr(addr)
	return d.readUint16(addr)
}

// ReadUint32 reads 32-bit word from EVE memory at address addr.
func (d *Driver) ReadUint32(addr int) uint32 {
	d.end()
	checkAddr(addr)
	return d.readUint32(addr)
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

// IntFlags reads interrupt flags from EVE and accumulates them (using logical
// OR) in internal variable. It returns accumulated flags.
func (d *Driver) IntFlags() byte {
	d.end()
	return d.intFlags()
}

// ClearIntFlags clears interrupt flags specified by mask.
func (d *Driver) ClearIntFlags(mask byte) {
	d.end()
	d.clearIntFlags(mask)
}

// IntMask returns current interrupt mask.
func (d *Driver) IntMask() byte {
	d.end()
	return d.intMask()
}

// SetIntMask sets interrupt mask.
func (d *Driver) SetIntMask(mask byte) {
	d.end()
	d.setIntMask(mask)
}

// Wait waits for any interrupt in flags.
func (d *Driver) Wait(flags byte) {
	d.end()
	d.wait(flags)
}

// SetBacklight sets backlight PWM duty cycle. Pwmduty range is from 0 to 128.
func (d *Driver) SetBacklight(pwmduty int) {
	d.end()
	d.writeByte(d.mmap.regintflags+opwmduty, byte(pwmduty))
}

// SwapDL clears INT_SWAP and schedules the display lists swap, to be performed
// after rendering the current frame.
func (d *Driver) SwapDL() {
	d.end()
	d.clearIntFlags(INT_SWAP)
	d.writeByte(d.mmap.regdlswap, DLSWAP_FRAME)
}
