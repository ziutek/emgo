package i2c

type MasterConnDMA struct {
	d      *DriverDMA
	addr   int16
	stopm  StopMode
	locked bool
}

// SetStopMode allows to enable/disable auto-stop mode for read and/or write
// operations. See StopMode for more information.
func (c *MasterConnDMA) SetStopMode(stopm StopMode) {
	c.stopm = stopm
}

func (c *MasterConnDMA) lockDriver() {
	if !c.locked {
		c.d.mutex.Lock()
		c.locked = true
	}
}

func (c *MasterConnDMA) unlockDriver() {
	c.locked = false
	c.d.mutex.Unlock()
}

func (c *MasterConnDMA) UnlockDriver() {
	c.d.Unlock()
}

// StopWrite terminates current write transaction and deactivates connection.
func (c *MasterConnDMA) StopWrite() {
	d := c.d
	if d.state == dstateIdle {
		return
	}
	p := &d.Periph.raw
	p.STOP().Set()
	d.state = dstateIdle
	c.unlockDriver()
}

// Write sends data from buf to slave device. If len(buf) == 0 Write does
// nothing, especially it does not activate inactiv connection nor interrupt
// previous read transaction.
func (c *MasterConnDMA) Write(buf []byte) (int, error) {
	if len(buf) == 0 {
		return 0, nil
	}
	c.lockDriver()
	d := c.d
	d.buf = buf
	d.addr = c.addr
	d.stop = false
	d.I2CISR()
	if e := d.waitDone(d.TxDMA); e != 0 {
		c.locked = false // d.Unlock must be used to unlock the driver.
		d.state = dstateIdle
		return d.n, e
	}
	if c.stopm&ASWR != 0 {
		c.StopWrite()
	}
	return len(buf), nil
}

// SetStopRead sets an internal flag which causes that subsequent read finishes
// transaction and deactivates connection. It can be called at any time, but if
// called after first read in current transaction, the subsequent read must read
// at least 2 bytes to properly generate stop condition on I2C bus.
func (c *MasterConnDMA) SetStopRead() {
	c.stopm |= stoprd
}

// Read reads data from slave device into buf. If len(buf) == 0 Read does
// nothing, especially it does not: activate inactiv connection, interrupt
// previous write transaction, deactivate connection if SetStopRead was called
// before.
func (c *MasterConnDMA) Read(buf []byte) (int, error) {
	if len(buf) == 0 {
		return 0, nil
	}
	if c.stopm&ASRD != 0 {
		c.SetStopRead()
	}
	c.lockDriver()
	d := c.d
	d.buf = buf
	d.n = 0
	d.addr = c.addr | 1
	d.stop = c.stopm&stoprd != 0
	if d.state == dstateWriteWait {
		d.state = dstateIdle
	}
	d.I2CISR()
	if e := d.waitDone(d.RxDMA); e != 0 {
		c.locked = false // d.Unlock must be used to unlock the driver.
		d.state = dstateIdle
		return d.n, e
	}
	if d.stop {
		c.unlockDriver()
		c.stopm &^= stoprd
	}
	return len(buf), nil
}

func (c *MasterConnDMA) WriteByte(b byte) error {
	_, err := c.Write([]byte{b})
	return err
}
