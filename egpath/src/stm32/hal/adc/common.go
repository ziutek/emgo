package adc

import (
	"unsafe"

	"stm32/hal/dma"
)

func panicCN() {
	panic("adc: bad channel number")
}

func panicSeq() {
	panic("adc: sequence too long")
}

func enableDMA(ch *dma.Channel, circ dma.Mode, half dma.Event,
	paddr, maddr unsafe.Pointer, wordSize uintptr, n int) {

	ch.Setup(dma.PTM | dma.IncM | dma.FT1 | circ)
	ch.SetWordSize(wordSize, wordSize)
	ch.SetAddrP(paddr)
	ch.SetAddrM(maddr)
	ch.SetLen(n)
	ch.Clear(dma.EvAll, dma.ErrAll)
	ch.EnableIRQ(dma.Complete|half, dma.ErrAll&^dma.ErrFIFO)
	ch.Enable()
}
