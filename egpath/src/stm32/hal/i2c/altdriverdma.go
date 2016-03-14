package i2c

import (
	"rtos"
	"sync"
	"unsafe"

	"stm32/hal/dma"

	"stm32/hal/raw/i2c"
)

// DriverDMA uses DMA to implement polling and interrupt driven driver to I2C
// peripheral. Default mode is polling.
type DriverDMA struct {
	*Periph
	RxDMA dma.Channel
	TxDMA dma.Channel

	mutex  sync.Mutex
	evflag rtos.EventFlag
	i2cint bool
	dmaint bool
}

// NewDriverDMA provides convenient way to create heap allocated DriverDMA
// struct.
func NewDriverDMA(p *Periph, rxch, txch dma.Channel) *DriverDMA {
	d := new(DriverDMA)
	d.Periph = p
	d.RxDMA = rxch
	d.TxDMA = txch
	return d
}

func (d *DriverDMA) SetIntMode(i2cen, dmaen bool) {
	d.i2cint = i2cen
	d.dmaint = dmaen
}

func (d *DriverDMA) I2CISR() {
	d.Periph.raw.CR2.ClearBits(i2c.ITBUFEN | i2c.ITEVTEN | i2c.ITERREN)
	d.evflag.Set()
}

func (d *DriverDMA) DMAISR(ch dma.Channel) {
	ch.DisableInt(dma.EV | dma.ERR)
	d.evflag.Set()
}

// MasterConn returns initialized MasterConnDMA struct that can be used to
// communicate with the slave device. Addr is the I2C address of the slave.
// See MasterConnDMA.SetStopMode for description of stopm.
func (d *DriverDMA) MasterConn(addr int16, stopm StopMode) MasterConnDMA {
	return MasterConnDMA{d: d, addr: uint16(addr << 1), stop: stopm}
}

// NewMasterConn is like MasterConn but returns pointer to heap allocated
// MasterConnDMA struct.
func (d *DriverDMA) NewMasterConn(addr int16, stopm StopMode) *MasterConnDMA {
	mc := new(MasterConnDMA)
	*mc = d.MasterConn(addr, stopm)
	return mc
}

func (d *DriverDMA) waitEvent(ev i2c.SR1_Bits) Error {
	p := &d.Periph.raw
	deadline := rtos.Nanosec() + byteTimeout
	if d.i2cint {
		return i2cWaitIRQ(p, &d.evflag, ev, deadline)

	}
	return i2cPollEvent(p, ev, deadline)
}

func (d *DriverDMA) startDMA(ch dma.Channel, addr *byte, n int) (int, Error) {
	ch.SetAddrM(unsafe.Pointer(addr))
	ch.SetLen(n)
	ch.Enable()
	// Set timeout to 2 * calculated transfer time.
	speed := d.Speed()
	deadline := rtos.Nanosec() + (2*9e9*int64(n)+int64(speed))/int64(speed)
	var e Error
	if d.dmaint {
		e = dmaWaitTRCE(ch, &d.evflag, deadline)
	} else {
		e = dmaPoolTRCE(ch, deadline)
	}
	ch.Disable()
	ch.ClearEvents(dma.EV | dma.ERR)
	return n - ch.Len(), e
}
