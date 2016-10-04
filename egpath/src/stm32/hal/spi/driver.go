package spi

import (
	"stm32/hal/dma"
)

type Driver struct {
	*Periph
	RxDMA *dma.Channel
	TxDMA *dma.Channel
}

// NewDriver provides convenient way to create heap allocated Driver struct.
func NewDriver(p *Periph, rxdma, txdma *dma.Channel) *Driver {
	d := new(Driver)
	d.Periph = p
	d.RxDMA = rxdma
	d.TxDMA = txdma
	return d
}
