package i2c

import (
	"sync"
)

// Driver implements interrupt driven driver for I2C peripheral.
type Driver struct {
	Periph

	mt   sync.Mutex
	wg   sync.WaitGroup
	data []byte
	n    uint32
}

func (d *Driver) MasterConn(addr int16) MasterConn {
	return MasterConn{d: d, addr: uint16(addr<<1) &^ started}
	// TODO: Add support for 10-bit addr.
}

func (d *Driver) ISR() {
	
	sr1 := d.sr1 | d.Periph.raw.SR1.Load()
}


func (d *Driver) enai() {
}
