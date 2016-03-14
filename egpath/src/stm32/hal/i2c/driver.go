package i2c

import (
	"rtos"
	"sync"

	"arch/cortexm"

	"stm32/hal/raw/i2c"
)

// Driver implements interrupt driven driver to I2C peripheral. Most of
// data/event processing is performed in interrupt handler mode.
type Driver struct {
	*Periph

	mutex sync.Mutex
	done  rtos.EventFlag
	buf   []byte
	n     int
	addr  int16
	stop  bool
	state byte
}

// NewDriver provides convenient way to create heap allocated Driver struct.
func NewDriver(p *Periph) *Driver {
	d := new(Driver)
	d.Periph = p
	return d
}

// MasterConn returns initialized MasterConn struct that can be used to
// communicate with the slave device. Addr is the I2C address of the slave.
func (d *Driver) MasterConn(addr int16, stopm StopMode) MasterConn {
	return MasterConn{d: d, addr: addr << 1, stopm: stopm}
	// TODO: Add support for 10-bit addr.
}

// NewMasterConn is like MasterConn but returns pointer to heap allocated
// MasterConn struct.
func (d *Driver) NewMasterConn(addr int16, stopm StopMode) *MasterConn {
	mc := new(MasterConn)
	*mc = d.MasterConn(addr, stopm)
	return mc
}

func (d *Driver) EventISR() {
	eventHandlers[d.state](d)
}

func (d *Driver) ErrorISR() {
	d.disableInt()
	d.state = stateI2CError
	d.done.Set()
}

// Unlock must be used after recovering from error.
func (d *Driver) Unlock() {
	d.mutex.Unlock()
}

func (d *Driver) enableInt() {
	d.Periph.raw.CR2.SetBits(i2c.ITEVTEN | i2c.ITERREN)
}

func (d *Driver) enableIntBuf() {
	d.Periph.raw.CR2.SetBits(i2c.ITEVTEN | i2c.ITERREN | i2c.ITBUFEN)
}

func (d *Driver) disableInt() {
	d.Periph.raw.CR2.ClearBits(i2c.ITEVTEN | i2c.ITERREN | i2c.ITBUFEN)
}

const (
	stateIdle = iota
	stateStart
	stateAddr

	stateWrite
	stateWriteWait

	stateRead
	stateRead1
	stateRead2
	stateReadN
	stateReadLast
	stateReadWait

	// All state numbers for errors should be >= stateI2CError
	stateI2CError
	stateBelatedStop
)

// eventHandlers is an array containing functions that performs state
// transitions in response to hardware or software events.
//
//emgo:const
var eventHandlers = [...]func(d *Driver){
	stateIdle:  (*Driver).idleEH,
	stateStart: (*Driver).startISR,
	stateAddr:  (*Driver).addrISR,

	stateWrite:     (*Driver).writeISR,
	stateWriteWait: (*Driver).writeWaitEH,

	stateRead:     (*Driver).readISR,
	stateReadWait: (*Driver).readWaitEH,
	stateRead1:    (*Driver).read1ISR,
	stateRead2:    (*Driver).read2ISR,
	stateReadN:    (*Driver).readNISR,
	stateReadLast: (*Driver).readLastISR,

	stateI2CError:    (*Driver).errorNOP,
	stateBelatedStop: (*Driver).errorNOP,
}

func (d *Driver) idleEH() {
	if d.addr < 0 {
		// Slave - not implemented.
	} else {
		// Master
		bits := i2c.START
		if d.addr&1 != 0 {
			bits |= i2c.ACK
		}
		d.Periph.raw.CR1.SetBits(bits)
		d.state = stateStart
	}
	d.enableInt()
}

func (d *Driver) startISR() {
	p := &d.Periph.raw
	if p.SR1.Load()&i2c.SB == 0 {
		return
	}
	d.state = stateAddr
	p.DR.Store(i2c.DR_Bits(d.addr))
}

func (d *Driver) addrISR() {
	p := &d.Periph.raw
	if p.SR1.Load()&i2c.ADDR == 0 {
		return
	}
	if d.addr&1 == 0 {
		// Write.
		d.n = 1
		d.state = stateWrite
		p.SR2.Load()
		p.DR.Store(i2c.DR_Bits(d.buf[0]))
		return
	}
	// Read.
	if !d.stop {
		d.state = stateRead
		p.SR2.Load()
		return
	}
	switch len(d.buf) {
	case 1:
		d.state = stateRead1
		p.ACK().Clear()
		cortexm.SetPRIMASK()
		p.SR2.Load()
		p.STOP().Set()
		cortexm.ClearPRIMASK()
		d.enableIntBuf()
	case 2:
		d.state = stateRead2
		p.POS().Set()
		cortexm.SetPRIMASK()
		p.SR2.Load()
		p.ACK().Clear()
		cortexm.ClearPRIMASK()
	default:
		d.state = stateReadN
		p.SR2.Load()
	}
}

func (d *Driver) writeISR() {
	p := &d.Periph.raw
	if p.SR1.Load()&i2c.BTF == 0 {
		return
	}
	n := d.n
	if n == len(d.buf) {
		d.disableInt()
		d.state = stateWriteWait
		d.done.Set()
		return
	}
	p.DR.Store(i2c.DR_Bits(d.buf[n]))
	d.n = n + 1
}

func (d *Driver) writeWaitEH() {
	d.Periph.raw.DR.Store(i2c.DR_Bits(d.buf[0]))
	d.n = 1
	d.state = stateWrite
	d.enableInt()
}

func (d *Driver) readISR() {
	p := &d.Periph.raw
	if p.SR1.Load()&i2c.BTF == 0 {
		return
	}
	d.buf[d.n] = byte(p.DR.Load())
	d.n++
	if d.n < len(d.buf) {
		return
	}
	d.disableInt()
	d.state = stateReadWait
	d.done.Set()
}

func (d *Driver) readWaitEH() {
	d.n = 0
	if d.stop {
		d.state = stateReadN
	} else {
		d.state = stateRead
	}
	d.enableInt()
}

func (d *Driver) read1ISR() {
	p := &d.Periph.raw
	if p.SR1.Load()&i2c.RXNE == 0 {
		return
	}
	d.disableInt()
	d.buf[0] = byte(p.DR.Load())
	d.n = 1
	d.state = stateIdle
	d.done.Set()
}

func (d *Driver) read2ISR() {
	p := &d.Periph.raw
	if p.SR1.Load()&i2c.BTF == 0 {
		return
	}
	cortexm.SetPRIMASK()
	p.STOP().Set()
	dr := p.DR.Load()
	cortexm.ClearPRIMASK()
	p.POS().Clear()
	d.disableInt()
	d.buf[0] = byte(dr)
	d.buf[1] = byte(p.DR.Load())
	d.n = 2
	d.state = stateIdle
	d.done.Set()
}

func (d *Driver) readNISR() {
	p := &d.Periph.raw
	if p.SR1.Load()&i2c.BTF == 0 {
		return
	}
	n := d.n
	m := len(d.buf) - n
	if m < 3 {
		d.disableInt()
		d.state = stateBelatedStop
		d.done.Set()
		return
	}
	if m == 3 {
		p.ACK().Clear()
		cortexm.SetPRIMASK()
		dr := p.DR.Load()
		p.STOP().Set()
		cortexm.ClearPRIMASK()
		d.buf[n] = byte(dr)
		n++
		d.state = stateReadLast
		d.enableIntBuf()
	}
	d.buf[n] = byte(p.DR.Load())
	d.n = n + 1
}

func (d *Driver) readLastISR() {
	p := &d.Periph.raw
	if p.SR1.Load()&i2c.RXNE == 0 {
		return
	}
	d.disableInt()
	n := d.n
	d.buf[n] = byte(p.DR.Load())
	d.n = n + 1
	d.state = stateIdle
	d.done.Set()
}

func (d *Driver) errorNOP() {
	// Ignore all events in error state.
}

func (d *Driver) waitDone(n int) (e Error) {
	deadline := rtos.Nanosec() + 2*9e9*int64(n+1)/int64(d.Speed())
	if d.done.Wait(deadline) {
		d.done.Clear()
	} else {
		e = SoftTimeout
	}
	if d.state >= stateI2CError {
		if d.state == stateBelatedStop {
			e |= BelatedStop
		}
		p := &d.Periph.raw
		e |= getError(p.SR1.Load())
		if e&Timeout == 0 {
			p.STOP().Set()
		}
		d.state = stateIdle
		p.SR1.Store(0) // Clear error flags.
	}
	return
}
