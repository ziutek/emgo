package i2c

import (
	"rtos"
	"sync"

	"arch/cortexm/nvic"

	"stm32/hal/raw/i2c"
)

type Error byte

const (
	BusErr      Error = 1 << 0
	ArbLost     Error = 1 << 1
	AckFail     Error = 1 << 2
	Overrun     Error = 1 << 3
	PECErr      Error = 1 << 4
	SoftTimeout Error = 1 << 5
	Timeout     Error = 1 << 6
	SMBAlert    Error = 1 << 7
)

func (e Error) Error() string {
	return "I2C error"
}

// Driver implements interrupt driven driver for I2C peripheral. By default
// driver works in polling mode.
type Driver struct {
	Periph

	mutex  sync.Mutex
	evflag rtos.EventFlag
	evirq  nvic.IRQ
	errirq nvic.IRQ
}

func (d *Driver) SetIntMode(evirq, errirq nvic.IRQ) {
	d.evirq = evirq
	d.errirq = errirq
	d.Periph.raw.CR2.SetBits(i2c.ITBUFEN | i2c.ITEVTEN | i2c.ITERREN)
}

func (d *Driver) ISR() {
	d.evirq.Disable()
	d.errirq.Disable()
	d.evflag.Set()
}

func (d *Driver) MasterConn(addr int16) MasterConn {
	return MasterConn{d: d, addr: uint16(addr<<1) &^ started}
	// TODO: Add support for 10-bit addr.
}
