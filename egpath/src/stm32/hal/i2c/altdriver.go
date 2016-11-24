package i2c

import (
	"rtos"
	"sync"

	"stm32/hal/raw/i2c"
)

// AltDriver implements polling and interrupt driven driver to I2C peripheral.
// Default mode is polling. Interrupt mode avoids active polling of event flags
// at cost of system call to sleep until interupt occurs. ISR is very short: it
// only informs thread part of AltDriver that some event occured. The all
// remaining porcessing of events/data for both modes is handled by the same
// code that is executed in thread mode.
type AltDriver struct {
	P *Periph

	mutex  sync.Mutex
	evflag rtos.EventFlag
	i2cint bool
}

// NewAltDriver provides convenient way to create heap allocated AltDriver
// struct.
func NewAltDriver(p *Periph) *AltDriver {
	d := new(AltDriver)
	d.P = p
	return d
}

// Unlock must be used after recovering from error.
func (d *AltDriver) Unlock() {
	d.mutex.Unlock()
}

// SetIntMode enables/disables interrupt mode.
func (d *AltDriver) SetIntMode(en bool) {
	d.i2cint = en
}

func (d *AltDriver) ISR() {
	d.P.raw.CR2.ClearBits(i2c.ITBUFEN | i2c.ITEVTEN | i2c.ITERREN)
	d.evflag.Signal(1)
}

// MasterConn returns initialized AltMasterConn struct that can be used to
// communicate with the slave device. Addr is the I2C address of the slave.
func (d *AltDriver) MasterConn(addr int16, stopm StopMode) AltMasterConn {
	return AltMasterConn{d: d, addr: uint16(addr << 1), stopm: stopm}
	// TODO: Add support for 10-bit addr.
}

// NewMasterConn is like MasterConn but returns pointer to heap allocated
// AltMasterConn struct.
func (d *AltDriver) NewMasterConn(addr int16, stopm StopMode) *AltMasterConn {
	mc := new(AltMasterConn)
	*mc = d.MasterConn(addr, stopm)
	return mc
}

func (d *AltDriver) waitEvent(ev i2c.SR1_Bits) Error {
	p := &d.P.raw
	deadline := rtos.Nanosec() + byteTimeout
	if d.i2cint {
		return i2cWaitIRQ(p, &d.evflag, ev, deadline)
	}
	return i2cPollEvent(p, ev, deadline)
}
