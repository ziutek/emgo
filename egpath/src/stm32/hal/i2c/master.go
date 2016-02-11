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

func (c *MasterConn) waitSR1(bit i2c.SR1_Bits) Error {
	timeout := rtos.Nanosec() + 100*1e6 // 100 ms timeout
	for {
		sr1 := c.d.raw.SR1.Load()
		if e := Error(sr1>>8) &^ (1 << 5); e != 0 {
			return e
		}
		if sr1&bit != 0 {
			return 0
		}
		if rtos.Nanosec() >= timeout {
			return SoftTimeout
		}
	}
}

/*func (c *MasterConn) Write(buf []byte) (int, error) {

}*/

const (
	started = 0x8000
)

func (c *MasterConn) Write(buf []byte) (int, error) {
	var (
		e Error
		n int
	)
	raw := c.d.Periph.raw
	if c.addr&started == 0 {
		c.d.mt.Lock()
		c.addr |= started
		raw.START().Set()
		if e = c.waitSR1(i2c.SB); e != 0 {
			goto end
		}
		raw.DR.Store(i2c.DR_Bits(c.addr)) // BUG: 10-bit addr not supported.
		if e = c.waitSR1(i2c.ADDR); e != 0 {
			goto end
		}
		raw.SR2.Load()
	}
	for _, b := range buf {
		raw.DR.Store(i2c.DR_Bits(b))
		if e = c.waitSR1(i2c.BTF); e != 0 {
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

// SetStop is equivalent to call SetStopMode(StopOnce, 0).
func (c *MasterConn) Stop() {
	c.d.Periph.raw.STOP().Set()
	c.addr &^= started
	c.d.mt.Unlock()
}
