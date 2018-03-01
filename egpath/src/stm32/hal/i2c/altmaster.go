package i2c

import (
	"arch/cortexm"

	"stm32/hal/raw/i2c"
)

// AltMasterConn can be used by I2c master. It represents a virtual connection
// to the slave device.
type AltMasterConn struct {
	d     *AltDriver
	addr  uint16
	stopm StopMode
	state byte
}

const (
	nact  = 0
	actrd = 1
	actwr = 2

	manstprd = 1 << 1
	manstpwr = 1 << 2
)

// SetStopMode allows to enable/disable auto-stop mode for read and/or write
// operations. See StopMode for more information.
func (c *AltMasterConn) SetStopMode(stopm StopMode) {
	c.stopm = stopm
}

func (c *AltMasterConn) UnlockDriver() {
	c.d.Unlock()
}

// StopWrite terminates current write transaction and deactivates connection.
func (c *AltMasterConn) StopWrite() {
	if c.state == actwr {
		c.d.P.raw.STOP().Set()
		c.state = nact
		c.d.mutex.Unlock()
	}
}

// Write sends data from buf to slave device. If len(buf) == 0 Write does
// nothing, especially it does not activate inactiv connection nor interrupt
// previous read transaction.
func (c *AltMasterConn) Write(buf []byte) (int, error) {
	if len(buf) == 0 {
		return 0, nil
	}
	var (
		e Error
		n int
	)
	d := c.d
	p := &d.P.raw
	if c.state != actwr {
		if c.state == actrd {
			return 0, ActiveRead
		}
		d.mutex.Lock()
		c.state = actwr
		p.START().Set()
		if e = d.waitEvent(i2c.SB); e != 0 {
			goto err
		}
		p.DR.U16.Store(c.addr) // BUG: 10-bit addr not supported.
		if e = d.waitEvent(i2c.ADDR); e != 0 {
			goto err
		}
		p.SR2.Load()
	}
	for m, b := range buf {
		p.DR.Store(i2c.DR(b))
		if e = d.waitEvent(i2c.BTF); e != 0 {
			n = m
			goto err
		}
	}
	if c.stopm&ASWR != 0 {
		c.StopWrite()
	}
	return len(buf), nil
err:
	p.SR1.Store(0) // Clear error flags.
	if e&Timeout == 0 {
		p.STOP().Set()
	}
	c.state = nact
	return n, e
}

// SetStopRead sets an internal flag which causes that subsequent read finishes
// transaction and deactivates connection. It can be called at any time, but if
// called after first read in current transaction, the subsequent read must read
// at least 2 bytes to properly generate stop condition on I2C bus.
func (c *AltMasterConn) SetStopRead() {
	c.stopm |= stoprd
}

// Read reads data from slave device into buf. If len(buf) == 0 Read does
// nothing, especially it does not: activate inactiv connection, interrupt
// previous write transaction, deactivate connection if SetStopRead was called
// before.
func (c *AltMasterConn) Read(buf []byte) (int, error) {
	if len(buf) == 0 {
		return 0, nil
	}
	if c.stopm&ASRD != 0 {
		c.SetStopRead()
	}
	var (
		e Error
		n int
	)
	d := c.d
	p := &d.P.raw
	stop := c.stopm & stoprd
	if c.state != actrd {
		if c.state == nact {
			d.mutex.Lock()
		}
		c.state = actrd
		p.CR1.SetBits(i2c.ACK | i2c.START)
		if e = d.waitEvent(i2c.SB); e != 0 {
			goto end
		}
		p.DR.U16.Store(c.addr | 1) // BUG: 10-bit addr not supported.
		if e = d.waitEvent(i2c.ADDR); e != 0 {
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
				if e = d.waitEvent(i2c.RXNE); e != 0 {
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
				if e = d.waitEvent(i2c.BTF); e != 0 {
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
			if e = d.waitEvent(i2c.BTF); e != 0 {
				goto end
			}
			buf[n] = byte(p.DR.Load())
		}
		if e = d.waitEvent(i2c.BTF); e != 0 {
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
		if e = d.waitEvent(i2c.RXNE); e != 0 {
			goto end
		}
		buf[n] = byte(p.DR.Load())
		n++
		goto end
	}
	for n = 0; n < len(buf); n++ {
		if e = d.waitEvent(i2c.BTF); e != 0 {
			goto end
		}
		buf[n] = byte(p.DR.Load())
	}
	return n, nil
end:
	c.stopm &^= stoprd
	c.state = nact
	if e != 0 {
		return n, e
	}
	c.d.mutex.Unlock()
	return n, nil
}

func (c *AltMasterConn) WriteByte(b byte) error {
	_, err := c.Write([]byte{b})
	return err
}
