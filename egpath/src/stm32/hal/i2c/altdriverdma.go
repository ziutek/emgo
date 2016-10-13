package i2c

import (
	"rtos"
	"sync"
	"unsafe"

	"stm32/hal/dma"

	"stm32/hal/raw/i2c"
)

// AltDriverDMA uses DMA to implement polling and interrupt driven driver to I2C
// peripheral. Default mode is polling.
type AltDriverDMA struct {
	P     *Periph
	RxDMA *dma.Channel
	TxDMA *dma.Channel

	mutex  sync.Mutex
	evflag rtos.EventFlag
	i2cint bool
	dmaint bool
}

// NewAltDriverDMA provides convenient way to create heap allocated AltDriverDMA
// struct.
func NewAltDriverDMA(p *Periph, rxch, txch *dma.Channel) *AltDriverDMA {
	d := new(AltDriverDMA)
	d.P = p
	d.RxDMA = rxch
	d.TxDMA = txch
	return d
}

// Unlock must be used after recovering from an error.
func (d *AltDriverDMA) Unlock() {
	d.mutex.Unlock()
}

func (d *AltDriverDMA) SetIntMode(i2cen, dmaen bool) {
	d.i2cint = i2cen
	d.dmaint = dmaen
}

func (d *AltDriverDMA) I2CISR() {
	d.P.raw.CR2.ClearBits(i2c.ITBUFEN | i2c.ITEVTEN | i2c.ITERREN)
	d.evflag.Set()
}

func (d *AltDriverDMA) DMAISR(ch *dma.Channel) {
	ch.DisableIRQ(dma.EvAll, dma.ErrAll)
	d.evflag.Set()
}

// MasterConn returns initialized AltMasterConnDMA struct that can be used to
// communicate with the slave device. Addr is the I2C address of the slave.
// See MasterConnDMA.SetStopMode for description of stopm.
func (d *AltDriverDMA) MasterConn(addr int16, stopm StopMode) AltMasterConnDMA {
	return AltMasterConnDMA{d: d, addr: uint16(addr << 1), stop: stopm}
}

// NewMasterConn is like MasterConn but returns pointer to heap allocated
// AltMasterConnDMA struct.
func (d *AltDriverDMA) NewMasterConn(addr int16, stopm StopMode) *AltMasterConnDMA {
	mc := new(AltMasterConnDMA)
	*mc = d.MasterConn(addr, stopm)
	return mc
}

func (d *AltDriverDMA) waitEvent(ev i2c.SR1_Bits) Error {
	p := &d.P.raw
	deadline := rtos.Nanosec() + byteTimeout
	if d.i2cint {
		return i2cWaitIRQ(p, &d.evflag, ev, deadline)

	}
	return i2cPollEvent(p, ev, deadline)
}

func (d *AltDriverDMA) performDMA(ch *dma.Channel, addr *byte, n int) (int, Error) {
	ch.SetAddrM(unsafe.Pointer(addr))
	ch.SetLen(n)
	ch.Clear(dma.EvAll, dma.ErrAll)
	timeout := byteTimeout + 2*9e9*int64(n+1)/int64(d.P.Speed())
	deadline := rtos.Nanosec() + timeout
	var e Error
	if d.dmaint {
		d.evflag.Clear()
		ch.EnableIRQ(dma.Complete, dmaErrMask)
		ch.Enable()
		e = dmaWait(ch, &d.evflag, deadline)
	} else {
		ch.Enable()
		e = dmaPoolTRCE(ch, deadline)
	}
	if e != 0 {
		n -= ch.Len()
	}
	return n, e
}
