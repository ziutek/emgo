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

type stregs struct {
	raw dma.DMA_Channel_Periph
}

func getStream(p *DMA, n int) *Stream {
	n--
	dman := dmanum(p)
	if dman == 0 && uint(n) > 6 || dman != 0 && uint(n) > 4 {
		panic(badStream)
	}
	return (*Stream)(unsafe.Pointer(&p.chs[n].raw))
}

// snum returns channel number - 1, eg. 0 for fisrt channel.
func snum(s *Stream) uintptr {
	off := uintptr(unsafe.Pointer(s)) & 0x3ff
	step := unsafe.Sizeof(dmaregs{}.chs[0])
	return (off - unsafe.Sizeof(dma.DMA_Periph{})) / step
}

func sdma(s *Stream) *dma.DMA_Periph {
	addr := uintptr(unsafe.Pointer(s)) &^ 0x3ff
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

func events(s *Stream) Events {
	isr := sdma(s).ISR.U32.Load()
	return Events(isr >> (snum(s) * 4) & 0xf)
}

func clearEvents(s *Stream, e Events) {
	mask := uint32(e&0xf) << (snum(s) * 4)
	sdma(s).IFCR.U32.Store(mask)
}

func enable(s *Stream) {
	s.raw.EN().Set()
}

func disable(s *Stream) {
	s.raw.EN().Clear()
}

func intEnabled(s *Stream) Events {
	return Events(s.raw.CCR.U32.Load() & 0xe)
}

func enableInt(s *Stream, e Events) {
	s.raw.CCR.U32.SetBits(uint32(e) & 0xe)
}

func disableInt(s *Stream, e Events) {
	s.raw.CCR.U32.ClearBits(uint32(e) & 0xe)
}

const modeMask = 0x70f0

func mode(s *Stream) Mode {
	return Mode(s.raw.CCR.Bits(modeMask))
}

func setup(s *Stream, m Mode, _ Channel) {
	s.raw.CCR.U32.StoreBits(0x70f0, uint32(m))
}

func wordSize(s *Stream) (p, m uintptr) {
	ccr := uintptr(s.raw.CCR.Load())
	p = 1 << (ccr >> 8 & 3)
	m = 1 << (ccr >> 10 & 3)
	return
}

func setWordSize(s *Stream, p, m uintptr) {
	ccr := p&6<<7 | m&6<<9
	s.raw.CCR.U32.StoreBits(0xf00, uint32(ccr))
}

func num(s *Stream) int {
	return int(s.raw.NDT().Load())
}

func setNum(s *Stream, n int) {
	s.raw.NDT().UM32.Store(uint32(n))
}

func setAddrP(s *Stream, a unsafe.Pointer) {
	s.raw.CPAR.U32.Store(uint32(uintptr(a)))
}

func setAddrM(s *Stream, a unsafe.Pointer) {
	s.raw.CMAR.U32.Store(uint32(uintptr(a)))
}
