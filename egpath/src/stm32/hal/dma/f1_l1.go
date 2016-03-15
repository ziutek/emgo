// +build f10x_ld f10x_ld_vl f10x_md f10x_md_vl f10x_hd f10x_hd_vl f10x_xl f10x_cl l1xx_md l1xx_mdp l1xx_hd l1xx_xl

package dma

import (
	"unsafe"

	"stm32/hal/raw/dma"
	"stm32/hal/raw/rcc"
)

type channel struct {
	raw dma.DMA_Channel_Periph
	_   uint32
}

type dmaperiph struct {
	raw dma.DMA_Periph
	chs [7]channel
}

func (p *DMA) enableClock(_ bool) {
	bit := bit(p, &rcc.RCC.AHBENR.U32, rcc.DMA1ENn)
	bit.Set()
	bit.Load() // RCC delay (workaround for silicon bugs).
}

func (p *DMA) disableClock() {
	bit(p, &rcc.RCC.AHBENR.U32, rcc.DMA1ENn).Clear()
}

func (p *DMA) reset() {}

func (p *DMA) getChannel(n, _ int) *Channel {
	n--
	dman := dmanum(p)
	if dman == 0 && uint(n) > 6 || dman != 0 && uint(n) > 4 {
		panic(badStream)
	}
	return (*Channel)(&p.chs[n])
}

func sdma(ch *Channel) *dma.DMA_Periph {
	addr := uintptr(unsafe.Pointer(ch)) &^ 0x3ff
	return (*dma.DMA_Periph)(unsafe.Pointer(addr))
}

// snum returns stream number - 1, eg. 0 for fisrt stream.
func snum(ch *Channel) uintptr {
	off := uintptr(unsafe.Pointer(ch)) & 0x3ff
	step := unsafe.Sizeof(channel{})
	return (off - unsafe.Sizeof(dma.DMA_Periph{})) / step
}

const (
	trce = 1 << 1
	htce = 1 << 2

	trerr = 1 << 3
	fferr = 0
	dmerr = 0
)

func (ch *Channel) events() Events {
	isr := sdma(ch).ISR.U32.Load()
	return Events(isr >> (snum(ch) * 4) & 0xf)
}

func (ch *Channel) clearEvents(e Events) {
	mask := uint32(e&0xf) << (snum(ch) * 4)
	sdma(ch).IFCR.U32.Store(mask)
}

func (ch *Channel) enable() {
	ch.raw.EN().Set()
}

func (ch *Channel) disable() {
	ch.raw.EN().Clear()
}

func (ch *Channel) intEnabled() Events {
	return Events(ch.raw.CCR.U32.Load() & 0xe)
}

func (ch *Channel) enableInt(e Events) {
	ch.raw.CCR.U32.SetBits(uint32(e) & 0xe)
}

func (ch *Channel) disableInt(e Events) {
	ch.raw.CCR.U32.ClearBits(uint32(e) & 0xe)
}

const (
	mtp = 1 << dma.DIRn
	mtm = 1 << dma.MEM2MEMn

	circ = 1 << dma.CIRCn
	incP = 1 << dma.PINCn
	incM = 1 << dma.MINCn

	prioM = 1 << dma.PLn
	prioH = 2 << dma.PLn
	prioV = 3 << dma.PLn

	fifo_1_4 = 0
	fifo_2_4 = 0
	fifo_3_4 = 0
	fifo_4_4 = 0
)

func (ch *Channel) setup(m Mode) {
	mask := dma.DIR | dma.MEM2MEM | dma.CIRC | dma.PINC | dma.MINC | dma.PL
	ch.raw.CCR.StoreBits(mask, dma.CCR_Bits(m))
}

func (ch *Channel) wordSize() (p, m uintptr) {
	ccr := uintptr(ch.raw.CCR.Load())
	p = 1 << (ccr >> 8 & 3)
	m = 1 << (ccr >> 10 & 3)
	return
}

func (ch *Channel) setWordSize(p, m uintptr) {
	ccr := p&6<<7 | m&6<<9
	ch.raw.CCR.U32.StoreBits(0xf00, uint32(ccr))
}

func (ch *Channel) len() int {
	return int(ch.raw.NDT().Load())
}

func (ch *Channel) setLen(n int) {
	ch.raw.NDT().UM32.Store(uint32(n))
}

func (ch *Channel) setAddrP(a unsafe.Pointer) {
	ch.raw.CPAR.U32.Store(uint32(uintptr(a)))
}

func (ch *Channel) setAddrM(a unsafe.Pointer) {
	ch.raw.CMAR.U32.Store(uint32(uintptr(a)))
}
