// +build f10x_ld f10x_ld_vl f10x_md f10x_md_vl f10x_hd f10x_hd_vl f10x_xl f10x_cl l1xx_md l1xx_mdp l1xx_hd l1xx_xl

package dma

import (
	"unsafe"

	"stm32/hal/raw/dma"
	"stm32/hal/raw/rcc"
)

type dmaregs struct {
	raw dma.DMA_Periph
	chs [7]struct {
		raw dma.DMA_Channel_Periph
		_   uint32
	}
}

func enableClock(p *DMA, _ bool) {
	bit := bit(p, &rcc.RCC.AHBENR.U32, rcc.DMA1ENn)
	bit.Set()
	bit.Load() // RCC delay (workaround for silicon bugs).
}

func disableClock(p *DMA) {
	bit(p, &rcc.RCC.AHBENR.U32, rcc.DMA1ENn).Clear()
}

func reset(p *DMA) {}

type channel struct {
	raw *dma.DMA_Channel_Periph
}

func getChannel(p *DMA, n, _ int) (ch Channel) {
	n--
	dman := dmanum(p)
	if dman == 0 && uint(n) > 6 || dman != 0 && uint(n) > 4 {
		panic(badStream)
	}
	ch.raw = &p.chs[n].raw
	return
}

// snum returns stream number - 1, eg. 0 for fisrt stream.
func snum(ch Channel) uintptr {
	off := uintptr(unsafe.Pointer(ch.raw)) & 0x3ff
	step := unsafe.Sizeof(dmaregs{}.chs[0])
	return (off - unsafe.Sizeof(dma.DMA_Periph{})) / step
}

func sdma(ch Channel) *dma.DMA_Periph {
	addr := uintptr(unsafe.Pointer(ch.raw)) &^ 0x3ff
	return (*dma.DMA_Periph)(unsafe.Pointer(addr))
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
)

const (
	tce = 1 << 1
	hce = 1 << 2
	err = 1 << 3
)

func events(ch Channel) Events {
	isr := sdma(ch).ISR.U32.Load()
	return Events(isr >> (snum(ch) * 4) & 0xf)
}

func clearEvents(ch Channel, e Events) {
	mask := uint32(e&0xf) << (snum(ch) * 4)
	sdma(ch).IFCR.U32.Store(mask)
}

func enable(ch Channel) {
	ch.raw.EN().Set()
}

func disable(ch Channel) {
	ch.raw.EN().Clear()
}

func intEnabled(ch Channel) Events {
	return Events(ch.raw.CCR.U32.Load() & 0xe)
}

func enableInt(ch Channel, e Events) {
	ch.raw.CCR.U32.SetBits(uint32(e) & 0xe)
}

func disableInt(ch Channel, e Events) {
	ch.raw.CCR.U32.ClearBits(uint32(e) & 0xe)
}

const modeMask = 0x70f0

func mode(ch Channel) Mode {
	return Mode(ch.raw.CCR.Bits(modeMask))
}

func setup(ch Channel, m Mode) {
	ch.raw.CCR.U32.StoreBits(0x70f0, uint32(m))
}

func wordSize(ch Channel) (p, m uintptr) {
	ccr := uintptr(ch.raw.CCR.Load())
	p = 1 << (ccr >> 8 & 3)
	m = 1 << (ccr >> 10 & 3)
	return
}

func setWordSize(ch Channel, p, m uintptr) {
	ccr := p&6<<7 | m&6<<9
	ch.raw.CCR.U32.StoreBits(0xf00, uint32(ccr))
}

func length(ch Channel) int {
	return int(ch.raw.NDT().Load())
}

func setLen(ch Channel, n int) {
	ch.raw.NDT().UM32.Store(uint32(n))
}

func setAddrP(ch Channel, a unsafe.Pointer) {
	ch.raw.CPAR.U32.Store(uint32(uintptr(a)))
}

func setAddrM(ch Channel, a unsafe.Pointer) {
	ch.raw.CMAR.U32.Store(uint32(uintptr(a)))
}
