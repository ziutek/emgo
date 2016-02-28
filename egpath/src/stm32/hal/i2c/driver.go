package i2c

import (
	"sync"

	"arch/cortexm/nvic"

	"stm32/hal/raw/i2c"
)

// Driver implements polling and interrupt driven driver to I2C peripheral.
// Default mode is polling.
type Driver struct {
	Periph

	mutex sync.Mutex
	irqs  irqs
}

// NewDriver provides convenient way to create heap allocated Driver struct.
func NewDriver(p Periph) *Driver {
	d := new(Driver)
	d.Periph = p
	return d
}

func (d *Driver) SetIntMode(evirq, errirq nvic.IRQ) {
	d.irqs.i2cev = evirq
	d.irqs.i2cerr = errirq
	d.Periph.raw.CR2.SetBits(i2c.ITBUFEN | i2c.ITEVTEN | i2c.ITERREN)
}

func (d *Driver) SetPollMode() {
	d.irqs.i2cev = 0
	d.irqs.i2cerr = 0
	d.Periph.raw.CR2.ClearBits(i2c.ITBUFEN | i2c.ITEVTEN | i2c.ITERREN)
}

func (d *Driver) ISR() {
	irqs := &d.irqs
	irqs.i2cev.Disable()
	irqs.i2cerr.Disable()
	irqs.evflag.Set()
}

// MasterConn returns initialized MasterConn struct that can be used to
// communicate with the slave device. Addr is the I2C address of the slave.
// See MasterConn.SetStopMode for description of stm.
func (d *Driver) MasterConn(addr int16, stm StopMode) MasterConn {
	return MasterConn{d: d, addr: uint16(addr << 1), stop: stm}
	// TODO: Add support for 10-bit addr.
}

// NewMasterConn is like MasterConn but returns pointer to heap allocated
// MasterConn struct.
func (d *Driver) NewMasterConn(addr int16, stm StopMode) *MasterConn {
	mc := new(MasterConn)
	*mc = d.MasterConn(addr, stm)
	return mc
}

func (d *Driver) i2cWaitEvent(ev i2c.SR1_Bits) Error {
	return i2cWaitEvent(d.Periph.raw, &d.irqs, ev)
}
