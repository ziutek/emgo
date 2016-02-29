package i2c

import (
	"rtos"
	"sync"

	"stm32/hal/raw/i2c"
)

// Driver implements polling and interrupt driven driver to I2C peripheral.
// Default mode is polling.
type Driver struct {
	*Periph

	mutex  sync.Mutex
	evflag rtos.EventFlag
	i2cint bool
}

// NewDriver provides convenient way to create heap allocated Driver struct.
func NewDriver(p *Periph) *Driver {
	d := new(Driver)
	d.Periph = p
	return d
}

func (d *Driver) SetIntMode(en bool) {
	d.i2cint = en
}

func (d *Driver) ISR() {
	d.Periph.raw.CR2.ClearBits(i2c.ITBUFEN | i2c.ITEVTEN | i2c.ITERREN)
	d.evflag.Set()
}

// MasterConn returns initialized MasterConn struct that can be used to
// communicate with the slave device. Addr is the I2C address of the slave.
// See MasterConn.SetStopMode for description of stopm.
func (d *Driver) MasterConn(addr int16, stopm StopMode) MasterConn {
	return MasterConn{d: d, addr: uint16(addr << 1), stop: stopm}
	// TODO: Add support for 10-bit addr.
}

// NewMasterConn is like MasterConn but returns pointer to heap allocated
// MasterConn struct.
func (d *Driver) NewMasterConn(addr int16, stopm StopMode) *MasterConn {
	mc := new(MasterConn)
	*mc = d.MasterConn(addr, stopm)
	return mc
}

func (d *Driver) i2cWaitEvent(ev i2c.SR1_Bits) Error {
	p := &d.Periph.raw
	deadline := rtos.Nanosec() + byteTimeout
	if d.i2cint {
		return i2cWaitIRQ(p, &d.evflag, ev, deadline)

	}
	return i2cPollEvent(p, ev, deadline)
}
