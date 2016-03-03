package i2c

import (
	"arch/cortexm"

	"stm32/hal/raw/i2c"
)

// MasterConn can be used by I2c master. It represents a virtual connection to
// the slave device.
//
// One peripheral supports multiple concurrent master connections. First read
// or write on inactive connection starts an I2C transaction and the connection
// becomes active until the transaction end. Peripheral supports only one
// active connection at the same time. Starting a subsequent transaction in
// other connection is blocked until the current transaction will end.
//
// Active connection supports both read and write transactions. There is no
// need to terminate write transaction before subsequent read transaction but
// read transaction must be terminated before subsequent write transaction. By
// default auto-stop is enabled for read and write operations.
type MasterConn struct {
	d     *Driver
	addr  uint16
	stop  StopMode
	state byte
}

type StopMode byte

const (
	NOAS StopMode = 0      // Manual mode (use SetStopRead, StopWrite).
	ASRD StopMode = 1 << 1 // Any read is finished by sending stop condition.
	ASWR StopMode = 1 << 2 // Any write is finished by sending stop condition.

	stoprd StopMode = 1 << 0
)

const (
	nact  = 0
	actrd = 1
	actwr = 2

	manstprd = 1 << 1
	manstpwr = 1 << 2
)

// SetStopMode allows to enable/disable auto-stop mode for read and/or write
// operations. If auto-stop mode is enabled for read/write then ay read/write
// operation is finished by sending stop condition on the I2C bus and leaves
// connection inactive. This mode improves ability to sharing I2C bus between
// multiple tasks but at the same time can degrade performance. It is not
// recommended to disable auto-stop mode for read operations.
func (c *MasterConn) SetStopMode(stm StopMode) {
	c.stop = stm
}

// Write sends data from buf to slave device. If len(buf) == 0 Write does
// nothing, especially it does not activate inactiv connection nor interrupt
// previous read transaction.
func (c *MasterConn) Write(buf []byte) (int, error) {
	if len(buf) == 0 {
		return 0, nil
	}
	var (
		e Error
		n int
	)
	d := c.d
	p := &d.Periph.raw
	if c.state != actwr {
		if c.state == actrd {
			return 0, ActiveRead
		}
		d.mutex.Lock()
		c.state = actwr
		p.START().Set()
		if e = i2cWaitEvent(d, i2c.SB); e != 0 {
			goto err
		}
		p.DR.Store(i2c.DR_Bits(c.addr)) // BUG: 10-bit addr not supported.
		if e = i2cWaitEvent(d, i2c.ADDR); e != 0 {
			goto err
		}
		p.SR2.Load()
	}
	for m, b := range buf {
		p.DR.Store(i2c.DR_Bits(b))
		if e = i2cWaitEvent(d, i2c.BTF); e != 0 {
			n = m
			goto err
		}
	}
	if c.stop&ASWR != 0 {
		c.StopWrite()
	}
	return len(buf), nil
err:
	p.SR1.Store(0) // Clear error flags.
	if e&Timeout == 0 {
		d.Periph.raw.STOP().Set()
	}
	c.state = nact
	d.mutex.Unlock()
	return n, e
}

// StopWrite terminates current write transaction and deactivates connection.
func (c *MasterConn) StopWrite() {
	if c.state == actwr {
		c.d.Periph.raw.STOP().Set()
		c.state = nact
		c.d.mutex.Unlock()
	}
}

// Read reads data from slave device into buf. If len(buf) == 0 Read does
// nothing, especially it does not: activate inactiv connection, interrupt
// previous write transaction, deactivate connection if SetStopRead was called
// before.
func (c *MasterConn) Read(buf []byte) (int, error) {
	if len(buf) == 0 {
		return 0, nil
	}
	if c.stop&ASRD != 0 {
		c.SetStopRead()
	}
	var (
		e Error
		n int
	)
	d := c.d
	p := &d.Periph.raw
	stop := c.stop & stoprd
	if c.state != actrd {
		if c.state == nact {
			d.mutex.Lock()
		}
		c.state = actrd
		p.CR1.SetBits(i2c.ACK | i2c.START)
		if e = i2cWaitEvent(d, i2c.SB); e != 0 {
			goto end
		}
		p.DR.Store(i2c.DR_Bits(c.addr | 1)) // BUG: 10-bit addr not supported.
		if e = i2cWaitEvent(d, i2c.ADDR); e != 0 {
			goto end
		}
		if stop != 0 {
			switch len(buf) {
			case 1:
				p.ACK().Clear()
				cortexm.SetPRIMASK()
				p.SR2.Load()
				p.STOP().Set()
				cortexm.ClearPRIMASK()
				if e = i2cWaitEvent(d, i2c.RXNE); e != 0 {
					goto end
				}
				buf[0] = byte(p.DR.Load())
				n = 1
				goto end
			case 2:
				p.POS().Set()
				cortexm.SetPRIMASK()
				p.SR2.Load()
				p.ACK().Clear()
				cortexm.ClearPRIMASK()
				if e = i2cWaitEvent(d, i2c.BTF); e != 0 {
					goto end
				}
				cortexm.SetPRIMASK()
				p.STOP().Set()
				d := p.DR.Load()
				cortexm.ClearPRIMASK()
				p.POS().Clear()
				buf[0] = byte(d)
				buf[1] = byte(p.DR.Load())
				n = 2
				goto end
			}
		}
		p.SR2.Load()
	}
	if stop != 0 {
		m := len(buf) - 3
		if m < 0 {
			e = BelatedStop
			goto end
		}
		for n = 0; n < m; n++ {
			if e = i2cWaitEvent(d, i2c.BTF); e != 0 {
				goto end
			}
			buf[n] = byte(p.DR.Load())
		}
		if e = i2cWaitEvent(d, i2c.BTF); e != 0 {
			goto end
		}
		p.ACK().Clear()
		cortexm.SetPRIMASK()
		b := p.DR.Load()
		p.STOP().Set()
		cortexm.ClearPRIMASK()
		buf[n] = byte(b)
		n++
		buf[n] = byte(p.DR.Load())
		n++
		if e = i2cWaitEvent(d, i2c.RXNE); e != 0 {
			goto end
		}
		buf[n] = byte(p.DR.Load())
		n++
		goto end
	}
	for n = 0; n < len(buf); n++ {
		if e = i2cWaitEvent(d, i2c.BTF); e != 0 {
			goto end
		}
		buf[n] = byte(p.DR.Load())
	}
	return n, nil
end:
	c.stop &^= stoprd
	c.state = nact
	d.mutex.Unlock()
	if e != 0 {
		return n, e
	}
	return n, nil
}

// SetStopRead sets an internal flag which causes that subsequent read finishes
// transaction and deactivates connection. It can be called at any time, but if
// called after first read in current transaction, the subsequent read must read
// at least 2 bytes to properly generate stop condition on I2C bus.
func (c *MasterConn) SetStopRead() {
	c.stop |= stoprd
}

func (c *MasterConn) WriteByte(b byte) error {
	_, err := c.Write([]byte{b})
	return err
}
