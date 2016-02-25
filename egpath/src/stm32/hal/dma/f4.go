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

type channel uintptr

func getChannel(p *DMA, sn, cn int) Channel {
	if uint(sn) > 7 {
		panic(badStream)
	}
	if uint(cn) > 7 {
		panic("dma: bad channel")
	}
	return Channel{channel(unsafe.Pointer(&p.sts[sn])) | channel(cn)}
}

func sdma(ch Channel) *dma.DMA_Periph {
	addr := ch.channel &^ 0x3ff
	return (*dma.DMA_Periph)(unsafe.Pointer(addr))
}

func sraw(ch Channel) *dma.DMA_Stream_Periph {
	return (*dma.DMA_Stream_Periph)(unsafe.Pointer(ch.channel &^ 7))
}

func snum(ch Channel) uintptr {
	off := uintptr(ch.channel) & 0x3ff
	step := unsafe.Sizeof(dma.DMA_Stream_Periph{})
	return (off - unsafe.Sizeof(dma.DMA_Periph{})) / step
}

func cnum(ch Channel) int { return int(ch.channel & 7) }

const (
	FFERR Events = 1 << 0 // FIFO error.
	DMERR Events = 1 << 2 // Direct mode error.
	TRERR Events = 1 << 3 // Transfer error.

	err = FFERR | DMERR | TRERR
	hce = 1 << 4
	tce = 1 << 5
)

func events(ch Channel) Events {
	d := sdma(ch)
	n := snum(ch)
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

func clearEvents(ch Channel, e Events) {
	d := sdma(ch)
	n := snum(ch)
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

func enable(ch Channel) {
	sraw(ch).EN().Set()
}

func disable(ch Channel) {
	sraw(ch).EN().Clear()
}

func intEnabled(ch Channel) Events {
	return Events(sraw(ch).CR.Load() & 0x1e << 1)
}

func enableInt(ch Channel, e Events) {
	sraw(ch).CR.U32.SetBits(uint32(e) >> 1 & 0x1e)
}

func disableInt(ch Channel, e Events) {
	sraw(ch).CR.U32.ClearBits(uint32(e) >> 1 & 0x1e)
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

func setup(ch Channel, m Mode) {
	cr := dma.CR_Bits(cnum(ch))<<dma.CHSELn | dma.CR_Bits(m)
	mask := dma.CHSEL | dma.PL | dma.MINC | dma.PINC | dma.CIRC | dma.DIR
	st := sraw(ch)
	st.CR.StoreBits(mask, cr)
	// Enable FIFO with threshold set to 1/2.
	st.FCR.Store(dma.DMDIS | dma.FTH_0)
}

func wordSize(ch Channel) (p, m uintptr) {
	cr := uintptr(sraw(ch).CR.Load())
	p = 1 << (cr >> 11 & 3)
	m = 1 << (cr >> 13 & 3)
	return
}

func setWordSize(ch Channel, p, m uintptr) {
	cr := p&6<<10 | m&6<<12
	sraw(ch).CR.U32.StoreBits(0x7800, uint32(cr))
}

func length(ch Channel) int {
	return int(sraw(ch).NDTR.Load() & 0xFFFF)
}

func setLen(ch Channel, n int) {
	sraw(ch).NDTR.U32.Store(uint32(n) & 0xFFFF)
}

func setAddrP(ch Channel, a unsafe.Pointer) {
	sraw(ch).PAR.U32.Store(uint32(uintptr(a)))
}

func setAddrM(ch Channel, a unsafe.Pointer) {
	sraw(ch).M0AR.U32.Store(uint32(uintptr(a)))
}
