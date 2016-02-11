package i2c

import (
	"rtos"

	"stm32/hal/raw/i2c"
)

// MasterConn represents virtual connection from master to slave device. One
// peripheral supports multiple concurrent master connections. First Read or
// Write on inactive connection starts I2C transaction and connection becomes
// active until transaction end. Peripheral supports only one active connection.
// Starting subsequent transaction (in other connection) is blocked until
// current transaction end.
type MasterConn struct {
	d    *Driver
	addr uint16
}

const timeout = 100 * 1e6 // 100 ms

/*func (c *MasterConn) poolEvent(event i2c.SR1_Bits) Error {
	deadline := rtos.Nanosec() + timeout
	for {
		sr1 := c.d.raw.SR1.Load()
		if e := Error(sr1>>8) &^ SoftTimeout; e != 0 {
			return e
		}
		if sr1&event != 0 {
			return 0
		}
		if rtos.Nanosec() >= deadline {
			return SoftTimeout
		}
	}
}*/

func (c *MasterConn) waitIRQ(event i2c.SR1_Bits) Error {
	deadline := rtos.Nanosec() + timeout
	for {
		rtos.IRQ(c.d.evirq).Enable()
		rtos.IRQ(c.d.errirq).Enable()
		if !c.d.evflag.Wait(deadline) {
			return SoftTimeout
		}
		c.d.evflag.Clear()
		sr1 := c.d.Periph.raw.SR1.Load()
		if e := Error(sr1>>8) &^ SoftTimeout; e != 0 {
			return e
		}
		if sr1&event != 0 {
			return 0
		}
	}
}

const started = 0x8000

func (c *MasterConn) Write(buf []byte) (int, error) {
	var (
		e Error
		n int
	)
	raw := c.d.Periph.raw
	if c.addr&started == 0 {
		c.d.mutex.Lock()
		c.addr |= started
		raw.START().Set()
		if e = c.waitIRQ(i2c.SB); e != 0 {
			goto end
		}
		raw.DR.Store(i2c.DR_Bits(c.addr)) // BUG: 10-bit addr not supported.
		if e = c.waitIRQ(i2c.ADDR); e != 0 {
			goto end
		}
		raw.SR2.Load()
	}
	for _, b := range buf {
		raw.DR.Store(i2c.DR_Bits(b))
		if e = c.waitIRQ(i2c.BTF); e != 0 {
			goto end
		}
		n++
	}
end:
	if e != 0 {
		raw.SR1.Store(0) // Clear error flags.
		if e&Timeout == 0 {
			c.Stop()
		}
		return n, e
	}
	return n, nil
}

func (c *MasterConn) Stop() {
	c.d.Periph.raw.STOP().Set()
	c.addr &^= started
	c.d.mutex.Unlock()
}
