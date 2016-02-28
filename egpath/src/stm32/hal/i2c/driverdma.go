package i2c

import (
	"rtos"
	"sync"
	"unsafe"

	"arch/cortexm/nvic"

	"stm32/hal/dma"

	"stm32/hal/raw/i2c"
)

// DriverDMA uses DMA to implement polling and interrupt driven driver to I2C
// peripheral. Default mode is polling.
type DriverDMA struct {
	Periph
	TxDMA dma.Channel
	RxDMA dma.Channel

	mutex sync.Mutex
	irqs  irqs
}

// NewDriverDMA provides convenient way to create heap allocated DriverDMA
// struct.
func NewDriverDMA(p Periph, txch, rxch dma.Channel) *DriverDMA {
	d := new(DriverDMA)
	d.Periph = p
	d.TxDMA = txch
	d.RxDMA = rxch
	return d
}

func (d *DriverDMA) SetIntMode(evirq, errirq nvic.IRQ) {
	d.irqs.i2cev = evirq
	d.irqs.i2cerr = errirq
	d.Periph.raw.CR2.SetBits(i2c.ITEVTEN | i2c.ITERREN)
}

func (d *DriverDMA) SetPollMode() {
	d.irqs.i2cev = 0
	d.irqs.i2cerr = 0
	d.Periph.raw.CR2.ClearBits(i2c.ITEVTEN | i2c.ITERREN)
}

func (d *DriverDMA) ISR() {
	irqs := &d.irqs
	irqs.i2cev.Disable()
	irqs.i2cerr.Disable()
	irqs.evflag.Set()
}

// MasterConn returns initialized MasterConnDMA struct that can be used to
// communicate with the slave device. Addr is the I2C address of the slave.
// See MasterConnDMA.SetStopMode for description of stm.
func (d *DriverDMA) MasterConn(addr int16, stm StopMode) MasterConnDMA {
	return MasterConnDMA{d: d, addr: uint16(addr << 1), stop: stm}
}

func (d *DriverDMA) i2cWaitEvent(ev i2c.SR1_Bits) Error {
	return i2cWaitEvent(d.Periph.raw, &d.irqs, ev)
}

func startDMA(ch dma.Channel, addr *byte, n, speed int) (int, Error) {
	ch.SetAddrM(unsafe.Pointer(addr))
	ch.SetLen(n)
	ch.Enable()
	// Set timeout to 2 * calculated transfer time.
	deadline := rtos.Nanosec() + (2*9e9*int64(n)+int64(speed))/int64(speed)
	e := dmaPoolTCE(ch, deadline)
	ch.Disable()
	ch.ClearEvents(dma.TCE | dma.HCE | dma.ERR)
	return n - ch.Len(), e
}
