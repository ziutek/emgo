package i2c

import (
	"rtos"
	"sync"

	"arch/cortexm/irq"
)

// Driver implements interrupt driven driver for I2C peripheral.
type Driver struct {
	Periph
	EventIRQ
	ErrorIRQ

	mt sync.Mutex
	ev rtos.Event
}

func (d *Driver) MasterConn(addr int16) MasterConn {
	return MasterConn{d: d, addr: uint16(addr<<1) &^ started}
	// TODO: Add support for 10-bit addr.
}

func (d *Driver) ISR() {
	d.ev.Send()
}
