package i2c

import (
	"unsafe"

	"arch/cortexm"

	"stm32/hal/dma"

	"stm32/hal/raw/i2c"
)

type AltMasterConnDMA struct {
	d     *AltDriverDMA
	addr  uint16
	stop  StopMode
	state byte
}

func (c *AltMasterConnDMA) UnlockDriver() {
	c.d.Unlock()
}

// StopWrite terminates current write transaction and deactivates connection.
func (c *AltMasterConnDMA) StopWrite() {
	if c.state == actwr {
		p := &c.d.Periph.raw
		p.DMAEN().Clear()
		p.STOP().Set()
		c.state = nact
		c.d.mutex.Unlock()
	}
}

// Write sends data from buf to slave device. If len(buf) == 0 Write does
// nothing, especially it does not activate inactiv connection nor interrupt
// previous read transaction.
func (c *AltMasterConnDMA) Write(buf []byte) (int, error) {
	if len(buf) == 0 {
		return 0, nil
	}
	var (
		e Error
		n int
	)
	d := c.d
	p := &d.Periph.raw
	txd := d.TxDMA
	if c.state != actwr {
		if c.state == actrd {
			return 0, ActiveRead
		}
		d.mutex.Lock()
		c.state = actwr
		txd.Setup(dma.MTP | dma.IncM | dma.FIFO_4_4)
		txd.SetWordSize(1, 1)
		txd.SetAddrP(unsafe.Pointer(p.DR.U16.Addr()))
		p.DMAEN().Set()
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
	for {
		m := len(buf) - n
		if m == 0 {
			break
		}
		if m == 1 {
			p.DR.Store(i2c.DR_Bits(buf[n]))
		} else {
			if m > 0xffff {
				m = 0xffff
			}
			if len(buf)-(n+m) == 1 {
				m-- // Avoid last transfer size 1.
			}
			m, e = d.startDMA(txd, &buf[n], m)
		}
		n += m
		if e != 0 {
			e |= getError(p.SR1.Load())
			goto err
		}
		if e = d.waitEvent(i2c.BTF); e != 0 {
			goto err
		}
	}
	if c.stop&ASWR != 0 {
		c.StopWrite()
	}
	return n, nil
err:
	p.SR1.Store(0) // Clear error flags.
	if e&Timeout == 0 {
		d.Periph.raw.STOP().Set()
	}
	c.state = nact
	return n, e
}

// SetStopRead sets an internal flag which causes that subsequent read finishes
// transaction and deactivates connection. It can be called at any time, but if
// called after first read in current transaction, the subsequent read must read
// at least 2 bytes to properly generate stop condition on I2C bus.
func (c *AltMasterConnDMA) SetStopRead() {
	c.stop |= stoprd
}

// Read reads data from slave device into buf. If len(buf) == 0 Read does
// nothing, especially it does not: activate inactiv connection, interrupt
// previous write transaction, deactivate connection if SetStopRead was called
// before.
func (c *AltMasterConnDMA) Read(buf []byte) (int, error) {
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
	rxd := d.RxDMA
	stop := c.stop & stoprd
	if c.state != actrd {
		if c.state == nact {
			d.mutex.Lock()
		}
		c.state = actrd
		p.CR1.SetBits(i2c.ACK | i2c.START)
		if e = d.waitEvent(i2c.SB); e != 0 {
			goto end
		}
		p.DR.Store(i2c.DR_Bits(c.addr | 1)) // BUG: 10-bit addr not supported.
		if e = d.waitEvent(i2c.ADDR); e != 0 {
			goto end
		}
		if stop != 0 && len(buf) == 1 {
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
		}
		rxd.Setup(dma.PTM | dma.IncM | dma.FIFO_1_4)
		rxd.SetWordSize(1, 1)
		rxd.SetAddrP(unsafe.Pointer(p.DR.U16.Addr()))
		p.DMAEN().Set()
		p.SR2.Load()
	}
	if e = d.waitEvent(i2c.BTF); e != 0 {
		goto end
	}
	if len(buf) == 1 {
		if stop != 0 {
			e = BelatedStop
			goto end
		}
		buf[0] = byte(p.DR.Load())
		return 1, nil
	}
	for {
		m := len(buf) - n
		if m == 0 {
			break
		}
		if m > 0xffff {
			m = 0xffff
		}
		if len(buf)-(n+m) == 1 {
			m-- // Avoid last transfer size 1.
		}
		if stop != 0 && n+m == len(buf) {
			p.LAST().Set()
		}
		m, e = d.startDMA(rxd, &buf[n], m)
		n += m
		if e != 0 {
			e |= getError(p.SR1.Load())
			goto end
		}
	}
	if stop == 0 {
		return n, nil
	}
end:
	p.CR2.ClearBits(i2c.DMAEN | i2c.LAST)
	c.stop &^= stoprd
	c.state = nact
	if e != 0 {
		return n, e
	}
	c.d.mutex.Unlock()
	return n, nil
}
