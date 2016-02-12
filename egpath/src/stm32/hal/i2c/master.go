package i2c

import (
	"stm32/hal/raw/i2c"
)

// MasterConn can be used by I2c master. It represent virtual connection to
// slave device. One peripheral supports multiple concurrent master connections.
// First Read or Write on inactive connection starts I2C transaction and
// connection becomes active until transaction end. Peripheral supports only one
// active connection at the same time. Starting subsequent transaction (in
// other connection) is blocked until current transaction will end.
type MasterConn struct {
	d    *Driver
	addr uint16
}

const started = 0x8000

func (c *MasterConn) Write(buf []byte) (int, error) {
	var (
		e Error
		n int
	)
	d := c.d
	p := d.Periph.raw
	if c.addr&started == 0 {
		d.mutex.Lock()
		c.addr |= started
		p.START().Set()
		if e = d.waitEvent(i2c.SB); e != 0 {
			goto end
		}
		p.DR.Store(i2c.DR_Bits(c.addr)) // BUG: 10-bit addr not supported.
		if e = d.waitEvent(i2c.ADDR); e != 0 {
			goto end
		}
		p.SR2.Load()
	}
	for _, b := range buf {
		p.DR.Store(i2c.DR_Bits(b))
		if e = d.waitEvent(i2c.BTF); e != 0 {
			goto end
		}
		n++
	}
end:
	if e != 0 {
		p.SR1.Store(0) // Clear error flags.
		if e&Timeout == 0 {
			c.Stop()
		}
		return n, e
	}
	return n, nil
}

// Stop sets connection to inactive state (terminates current transaction).
func (c *MasterConn) Stop() {
	c.d.Periph.raw.STOP().Set()
	c.addr &^= started
	c.d.mutex.Unlock()
}
