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
	state state
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

// Unlock must be used after recovering from error.
func (d *Driver) Unlock() {
	d.mutex.Unlock()
}

func (d *Driver) I2CISR() {
	sr1 := d.Periph.raw.SR1.Load()
	if e := getError(sr1); e != 0 {
		d.handleErrors()
		return
	}
	eventHandlers[d.state](d, sr1)
}

type state byte

const (
	stateIdle state = iota
	stateStart
	stateAddr

	stateWrite
	stateWriteWait

	stateRead
	stateRead1
	stateRead2
	stateReadN
	stateReadWait

	// All state numbers for errors should be >= stateI2CError
	stateError
	stateBelatedStop
	stateBadEvent
)

// eventHandlers is an array containing functions that performs state
// transitions in response to hardware or software events.
//
//emgo:const
var eventHandlers = [...]func(d *Driver, sr1 i2c.SR1_Bits){
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
}

func (d *Driver) idleEH(sr1 i2c.SR1_Bits) {
	p := &d.Periph.raw
	if d.addr < 0 {
		// Slave - not implemented.
	} else {
		// Master
		bits := i2c.START
		if d.addr&1 != 0 {
			bits |= i2c.ACK
		}
		p.CR1.SetBits(bits)
		d.state = stateStart
	}
	// In case of repeated start ensure that BTF was cleared before enable
	// interrupts.
	for sr1&i2c.BTF != 0 {
		sr1 = p.SR1.Load()
		if e := getError(sr1); e != 0 {
			d.handleErrors()
			return
		}
	}
	d.enableIntI2C(i2c.ITEVTEN | i2c.ITERREN)
}

func (d *Driver) startISR(sr1 i2c.SR1_Bits) {
	if sr1&i2c.SB == 0 {
		d.badEvent(sr1)
		return
	}
	d.state = stateAddr
	p := &d.Periph.raw
	p.DR.Store(i2c.DR_Bits(d.addr))
	p.SR1.Load() // Ensure that SB was cleared before return.
}

func (d *Driver) addrISR(sr1 i2c.SR1_Bits) {
	if sr1&i2c.ADDR == 0 {
		d.badEvent(sr1)
		return
	}
	p := &d.Periph.raw
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
		d.enableIntI2C(i2c.ITEVTEN | i2c.ITERREN | i2c.ITBUFEN)
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

func (d *Driver) writeISR(sr1 i2c.SR1_Bits) {
	if sr1&i2c.BTF == 0 {
		d.badEvent(sr1)
		return
	}
	n := d.n
	if n == len(d.buf) {
		d.disableIntI2C(i2c.ITEVTEN | i2c.ITERREN | i2c.ITBUFEN)
		d.state = stateWriteWait
		d.done.Set()
		return
	}
	d.Periph.raw.DR.Store(i2c.DR_Bits(d.buf[n]))
	d.n = n + 1
}

func (d *Driver) writeWaitEH(sr1 i2c.SR1_Bits) {
	d.Periph.raw.DR.Store(i2c.DR_Bits(d.buf[0]))
	d.n = 1
	d.state = stateWrite
	d.enableIntI2C(i2c.ITEVTEN | i2c.ITERREN)
}

func (d *Driver) readISR(sr1 i2c.SR1_Bits) {
	if sr1&i2c.BTF == 0 {
		d.badEvent(sr1)
		return
	}
	d.buf[d.n] = byte(d.Periph.raw.DR.Load())
	d.n++
	if d.n == len(d.buf) {
		d.disableIntI2C(i2c.ITEVTEN | i2c.ITERREN | i2c.ITBUFEN)
		d.state = stateReadWait
		d.done.Set()
	}
}

func (d *Driver) readWaitEH(sr1 i2c.SR1_Bits) {
	if d.stop {
		d.state = stateReadN
	} else {
		d.state = stateRead
	}
	d.enableIntI2C(i2c.ITEVTEN | i2c.ITERREN)
}

func (d *Driver) read1ISR(sr1 i2c.SR1_Bits) {
	if sr1&i2c.RXNE == 0 {
		d.badEvent(sr1)
		return
	}
	d.disableIntI2C(i2c.ITEVTEN | i2c.ITERREN | i2c.ITBUFEN)
	n := d.n
	d.buf[n] = byte(d.Periph.raw.DR.Load())
	d.n = n + 1
	d.state = stateIdle
	d.done.Set()
}

func (d *Driver) read2ISR(sr1 i2c.SR1_Bits) {
	if sr1&i2c.BTF == 0 {
		d.badEvent(sr1)
		return
	}
	p := &d.Periph.raw
	cortexm.SetPRIMASK()
	p.STOP().Set()
	dr := p.DR.Load()
	cortexm.ClearPRIMASK()
	p.POS().Clear()
	d.disableIntI2C(i2c.ITEVTEN | i2c.ITERREN | i2c.ITBUFEN)
	d.buf[0] = byte(dr)
	d.buf[1] = byte(p.DR.Load())
	d.n = 2
	d.state = stateIdle
	d.done.Set()
}

func (d *Driver) readNISR(sr1 i2c.SR1_Bits) {
	if sr1&i2c.BTF == 0 {
		d.badEvent(sr1)
		return
	}
	p := &d.Periph.raw
	n := d.n
	m := len(d.buf) - n
	if m < 3 {
		d.disableIntI2C(i2c.ITEVTEN | i2c.ITERREN | i2c.ITBUFEN)
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
		d.state = stateRead1
		d.enableIntI2C(i2c.ITEVTEN | i2c.ITERREN | i2c.ITBUFEN)
	}
	d.buf[n] = byte(p.DR.Load())
	d.n = n + 1
}

func (d *Driver) enableIntI2C(m i2c.CR2_Bits) {
	d.Periph.raw.CR2.SetBits(m)
}

func (d *Driver) disableIntI2C(m i2c.CR2_Bits) {
	d.Periph.raw.CR2.ClearBits(m)
}


func (d *Driver) handleErrors() {
	d.disableIntI2C(i2c.ITEVTEN | i2c.ITERREN | i2c.ITBUFEN)
	if d.state < stateError {
		d.state = stateError
	}
	d.done.Set()
}

func (d *Driver) badEvent(sr1 i2c.SR1_Bits) {
	d.state = stateBadEvent
	d.handleErrors()
}

func (d *Driver) waitDone(n int) (e Error) {
	deadline := rtos.Nanosec() + 100e6 + 2*9e9*int64(n+1)/int64(d.Speed())
	if d.done.Wait(deadline) {
		d.done.Clear()
	} else {
		e = SoftTimeout
	}
	if d.state >= stateError {
		if d.state == stateBelatedStop {
			e |= BelatedStop
		}
		p := &d.Periph.raw
		e |= getError(p.SR1.Load())
		if e&Timeout == 0 {
			p.STOP().Set()
		}
		p.SR1.Store(0) // Clear error flags.
	}
	return
}
