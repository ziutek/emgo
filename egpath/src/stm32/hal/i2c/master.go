package i2c

import (
	"arch/cortexm"
	"stm32/hal/raw/i2c"
)

// MasterConn can be used by I2c master. It represents a virtual connection to
// the slave device. One peripheral supports multiple concurrent master
// connections. First read or write on inactive connection starts an I2C
// transaction and the connection becomes active until the transaction end.
// Peripheral supports only one active connection at the same time. Starting a
// subsequent transaction in other connection is blocked until the current
// transaction will end. Active connection supports both read and write
// transactions, which can be interleaved. Change of direction interrupts
// previous transaction but keeps the connection active.
type MasterConn struct {
	d      *Driver
	addr   uint16
	state  byte
	stoprd bool
}

const (
	nact  = 0
	actrd = 1
	actwr = 2
)

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
	p := d.Periph.raw
	if c.state != actwr {
		if c.state == nact {
			d.mutex.Lock()
		}
		c.state = actwr
		p.START().Set()
		if e = d.waitEvent(i2c.SB); e != 0 {
			goto err
		}
		p.DR.Store(i2c.DR_Bits(c.addr)) // BUG: 10-bit addr not supported.
		if e = d.waitEvent(i2c.ADDR); e != 0 {
			goto err
		}
		p.SR2.Load()
	}
	for _, b := range buf {
		p.DR.Store(i2c.DR_Bits(b))
		if e = d.waitEvent(i2c.BTF); e != 0 {
			goto err
		}
		n++
	}
	return n, nil
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
	var err error
	d := c.d
	p := d.Periph.raw
	if c.state != actrd {
		if c.state == nact {
			d.mutex.Lock()
		}
		c.state = actrd
		p.CR1.SetBits(i2c.ACK | i2c.START)
		if e := d.waitEvent(i2c.SB); e != 0 {
			err = e
			goto end
		}
		p.DR.Store(i2c.DR_Bits(c.addr | 1)) // BUG: 10-bit addr not supported.
		if e := d.waitEvent(i2c.ADDR); e != 0 {
			err = e
			goto end
		}
		if c.stoprd {
			switch len(buf) {
			case 1:
				p.ACK().Clear()
				cortexm.SetPRIMASK()
				p.SR2.Load()
				p.STOP().Set()
				cortexm.ClearPRIMASK()
				if e := d.waitEvent(i2c.RXNE); e != 0 {
					err = e
					goto end
				}
				buf[0] = byte(p.DR.Load())
				goto end
			case 2:
				p.POS().Set()
				cortexm.SetPRIMASK()
				p.SR2.Load()
				p.ACK().Clear()
				cortexm.ClearPRIMASK()
				if e := d.waitEvent(i2c.BTF); e != 0 {
					err = e
					goto end
				}
				cortexm.SetPRIMASK()
				p.STOP().Set()
				d := p.DR.Load()
				cortexm.ClearPRIMASK()
				p.POS().Clear()
				buf[0] = byte(d)
				buf[1] = byte(p.DR.Load())
				goto end
			}
		}
		p.SR2.Load()
	}
	if c.stoprd {
		n := len(buf) - 2
		if n < 0 {
			err = BelatedStop
			goto end
		}
		for i := 0; i < n; i++ {
			if e := d.waitEvent(i2c.BTF); e != 0 {
				err = e
				goto end
			}
			buf[i] = byte(p.DR.Load())
		}
		p.ACK().Clear()
		cortexm.SetPRIMASK()
		d := p.DR.Load()
		p.STOP().Set()
		cortexm.ClearPRIMASK()
		buf[n] = byte(d)
		buf[n+1] = byte(p.DR.Load())
		goto end
	}
	for i := range buf {
		if e := d.waitEvent(i2c.BTF); e != 0 {
			err = e
			goto end
		}
		buf[i] = byte(p.DR.Load())
	}
	return len(buf), nil
end:
	c.stoprd = false
	c.state = nact
	d.mutex.Unlock()
	return len(buf), err
}

// SetStopRead sets an internal flag which causes that subsequent read finishes
// transaction and deactivates connection. It can be called at any time, but if
// called after first read in current transaction, the subsequent read must read
// at least 2 bytes to properly generate stop condition on I2C bus.
func (c *MasterConn) SetStopRead() {
	c.stoprd = true
}
