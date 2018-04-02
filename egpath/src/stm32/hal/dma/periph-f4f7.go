// +build f40_41xxx f411xe f746xx

package dma

import (
	"mmio"
	"unsafe"

	"stm32/hal/raw/dma"
)

type dmaperiph struct {
	raw dma.DMA_Periph
	sts [8]dma.DMA_Stream_Periph
}

type channel struct {
	_ [1<<31 - 1]byte // Prevent allocation.
}

func (p *DMA) getChannel(sn, cn int) *Channel {
	if uint(sn) > 7 {
		panic(badStream)
	}
	if uint(cn) > 7 {
		panic("dma: bad channel")
	}
	addr := uintptr(unsafe.Pointer(&p.sts[sn])) | uintptr(cn)
	return (*Channel)(unsafe.Pointer(addr))
}

func addr(ch *Channel) uintptr {
	return uintptr(unsafe.Pointer(ch))
}

func sdma(ch *Channel) *dma.DMA_Periph {
	addr := addr(ch) &^ 0x3ff
	return (*dma.DMA_Periph)(unsafe.Pointer(addr))
}

func sraw(ch *Channel) *dma.DMA_Stream_Periph {
	addr := addr(ch) &^ 7
	return (*dma.DMA_Stream_Periph)(unsafe.Pointer(addr))
}

func snum(ch *Channel) uintptr {
	off := addr(ch) & 0x3ff
	step := unsafe.Sizeof(dma.DMA_Stream_Periph{})
	return (off - unsafe.Sizeof(dma.DMA_Periph{})) / step
}

func cnum(ch *Channel) int {
	return int(addr(ch) & 7)
}

const (
	fferr = 1 << 0
	dmerr = 1 << 2
	trerr = 1 << 3

	htce = 1 << 4
	trce = 1 << 5
)

func (ch *Channel) status() byte {
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
	return byte(isr.Load() >> n & 0x3d)
}

func (ch *Channel) clear(flags byte) {
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
	ifcr.Store(uint32(flags) & 0x3d << n)
}

func (ch *Channel) enable() {
	sraw(ch).EN().Set()
}

func (ch *Channel) disable() {
	sraw(ch).EN().Clear()
}

func (ch *Channel) enabled() bool {
	return sraw(ch).EN().Load() != 0
}

func (ch *Channel) irqEnabled() byte {
	st := sraw(ch)
	ev := byte(st.CR.Load()&0x1e<<1) | byte(st.FCR.Load()>>7&1)
	return ev
}

func (ch *Channel) enableIRQ(flags byte) {
	st := sraw(ch)
	st.CR.U32.SetBits(uint32(flags) >> 1 & 0x1e)
	//st.FCR.U32.SetBits(uint32(flags) & 1 << 7) Do not use
}

func (ch *Channel) disableIRQ(flags byte) {
	st := sraw(ch)
	st.CR.U32.ClearBits(uint32(flags) >> 1 & 0x1e)
	st.FCR.U32.ClearBits(uint32(flags) & 1 << 7)
}

const (
	fifo_1_4 = 4 << 0
	fifo_2_4 = 5 << 0
	fifo_3_4 = 6 << 0
	fifo_4_4 = 7 << 0

	mtp = 1 << 6
	mtm = 2 << 6

	circ = 1 << 8
	incP = 1 << 9
	incM = 1 << 10
)

func (ch *Channel) setup(m Mode) {
	cr := dma.CR(cnum(ch))<<dma.CHSELn | dma.CR(m)
	mask := dma.CHSEL | dma.PL | dma.MINC | dma.PINC | dma.CIRC | dma.DIR
	st := sraw(ch)
	st.CR.StoreBits(mask, cr)
	st.FCR.StoreBits(dma.DMDIS|dma.FTH, dma.FCR(m))
}

const (
	prioM = 1
	prioH = 2
	prioV = 3
)

func (ch *Channel) setPrio(prio Prio) {
	sraw(ch).PL().Store(dma.CR(prio) << dma.PLn)
}

func (ch *Channel) prio() Prio {
	return Prio(sraw(ch).PL().Load() >> dma.PLn)
}

func (ch *Channel) wordSize() (p, m uintptr) {
	cr := uintptr(sraw(ch).CR.Load())
	p = 1 << (cr >> 11 & 3)
	m = 1 << (cr >> 13 & 3)
	return
}

func (ch *Channel) setWordSize(p, m uintptr) {
	cr := p&6<<10 | m&6<<12
	sraw(ch).CR.U32.StoreBits(0x7800, uint32(cr))
}

func (ch *Channel) len() int {
	return int(sraw(ch).NDTR.Load() & 0xFFFF)
}

func (ch *Channel) setLen(n int) {
	sraw(ch).NDTR.U32.Store(uint32(n) & 0xFFFF)
}

func (ch *Channel) setAddrP(a unsafe.Pointer) {
	sraw(ch).PAR.U32.Store(uint32(uintptr(a)))
}

func (ch *Channel) setAddrM(a unsafe.Pointer) {
	sraw(ch).M0AR.U32.Store(uint32(uintptr(a)))
}

func (ch *Channel) request() Request {
	return -1
}

func (ch *Channel) setRequest(_ Request) {
}
