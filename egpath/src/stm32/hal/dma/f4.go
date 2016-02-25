// +build f40_41xxx f411xe

package dma

import (
	"bits"
	"mmio"
	"unsafe"

	"stm32/hal/raw/dma"
	"stm32/hal/raw/rcc"
)

type dmaregs struct {
	raw dma.DMA_Periph
	sts [8]dma.DMA_Stream_Periph
}

func enableClock(p *DMA, lp bool) {
	enbit := bit(p, &rcc.RCC.AHB1ENR.U32, rcc.DMA1ENn)
	enbit.Set()
	bit(p, &rcc.RCC.AHB1LPENR.U32, rcc.DMA1LPENn).Store(bits.One(lp))
	enbit.Load() // RCC delay (workaround for silicon bugs).
}

func disableClock(p *DMA) {
	bit(p, &rcc.RCC.AHB1ENR.U32, rcc.DMA1ENn).Clear()
}

func reset(p *DMA) {
	bit := bit(p, &rcc.RCC.AHB1RSTR.U32, rcc.DMA1RSTn)
	bit.Set()
	bit.Clear()
}

type stregs struct {
	raw dma.DMA_Stream_Periph
}

func getStream(p *DMA, n int) *Stream {
	if uint(n) > 7 {
		panic(badStream)
	}
	return (*Stream)(unsafe.Pointer(&p.sts[n]))
}

// snum returns stream number.
func snum(s *Stream) uintptr {
	off := uintptr(unsafe.Pointer(s)) & 0x3ff
	step := unsafe.Sizeof(dmaregs{}.sts[0])
	return (off - unsafe.Sizeof(dma.DMA_Periph{})) / step
}

func sdma(s *Stream) *dma.DMA_Periph {
	addr := uintptr(unsafe.Pointer(s)) &^ 0x3ff
	return (*dma.DMA_Periph)(unsafe.Pointer(addr))
}

const (
	FFERR Events = 1 << 0 // FIFO error.
	DMERR Events = 1 << 2 // Direct mode error.
	TRERR Events = 1 << 3 // Transfer error.

	err = FFERR | DMERR | TRERR
	hce = 1 << 4
	tce = 1 << 5
)

func events(s *Stream) Events {
	d := sdma(s)
	n := snum(s)
	var isr *mmio.U32
	if n < 4 {
		isr = &d.LISR.U32
	} else {
		isr = &d.HISR.U32
		n -= 4
	}
	n *= 6
	if n > 6 {
		n += 4
	}
	return Events(isr.Load() >> n & 0x3d)
}

func clearEvents(s *Stream, e Events) {
	d := sdma(s)
	n := snum(s)
	var ifcr *mmio.U32
	if n < 4 {
		ifcr = &d.LIFCR.U32
	} else {
		ifcr = &d.HIFCR.U32
		n -= 4
	}
	n *= 6
	if n > 6 {
		n += 4
	}
	ifcr.Store(uint32(e) & 0x3d << n)
}

func enable(s *Stream) {
	s.raw.EN().Set()
}

func disable(s *Stream) {
	s.raw.EN().Clear()
}

func intEnabled(s *Stream) Events {
	return Events(s.raw.CR.Load() & 0x1e << 1)
}

func enableInt(s *Stream, e Events) {
	s.raw.CR.U32.SetBits(uint32(e) >> 1 & 0x1e)
}

func disableInt(s *Stream, e Events) {
	s.raw.CR.U32.ClearBits(uint32(e) >> 1 & 0x1e)
}

const (
	mtp = 1 << dma.DIRn
	mtm = 2 << dma.DIRn

	circ = 1 << dma.CIRCn
	incP = 1 << dma.PINCn
	incM = 1 << dma.MINCn

	prioM = 1 << dma.PLn
	prioH = 2 << dma.PLn
	prioV = 3 << dma.PLn
)

func setup(s *Stream, m Mode, ch Channel) {
	if ch&^7 != 0 {
		panic("dma: bad channel")
	}
	cr := dma.CR_Bits(ch)<<dma.CHSELn | dma.CR_Bits(m)
	mask := dma.CHSEL | dma.PL | dma.MINC | dma.PINC | dma.CIRC | dma.DIR
	s.raw.CR.StoreBits(mask, cr)
	// Enable FIFO with threshold set to 1/2.
	s.raw.FCR.Store(dma.DMDIS | dma.FTH_0)
}

func wordSize(s *Stream) (p, m uintptr) {
	cr := uintptr(s.raw.CR.Load())
	p = 1 << (cr >> 11 & 3)
	m = 1 << (cr >> 13 & 3)
	return
}

func setWordSize(s *Stream, p, m uintptr) {
	cr := p&6<<10 | m&6<<12
	s.raw.CR.U32.StoreBits(0x7800, uint32(cr))
}

func num(s *Stream) int {
	return int(s.raw.NDTR.Load() & 0xFFFF)
}

func setNum(s *Stream, n int) {
	s.raw.NDTR.U32.Store(uint32(n) & 0xFFFF)
}

func setAddrP(s *Stream, a unsafe.Pointer) {
	s.raw.PAR.U32.Store(uint32(uintptr(a)))
}

func setAddrM(s *Stream, a unsafe.Pointer) {
	s.raw.M0AR.U32.Store(uint32(uintptr(a)))
}
