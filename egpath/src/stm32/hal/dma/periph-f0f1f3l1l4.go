// +build f030x6 f030x8 f10x_ld f10x_ld_vl f10x_md f10x_md_vl f10x_hd f10x_hd_vl f10x_xl f10x_cl f303xe l1xx_md l1xx_mdp l1xx_hd l1xx_xl l476xx

package dma

import (
	"mmio"
	"unsafe"

	"stm32/hal/raw/dma"
)

type chanregs struct {
	raw dma.DMA_Channel_Periph
	_   uint32
}

type dmaperiph struct {
	raw   dma.DMA_Periph
	chs   [7]chanregs
	_     [5]uint32
	cselr mmio.U32
}

type channel struct {
	chanregs
	_ [1<<31 - unsafe.Sizeof(chanregs{}) - 4]byte // Prevent allocation.
}

func (p *DMA) getChannel(n, _ int) *Channel {
	n--
	dman := dmanum(p)
	if dman == 0 && uint(n) > 6 || dman != 0 && uint(n) > 4 {
		panic(badStream)
	}
	return (*Channel)(unsafe.Pointer(&p.chs[n]))
}

func sdma(ch *Channel) *dmaperiph {
	addr := uintptr(unsafe.Pointer(ch)) &^ 0x3ff
	return (*dmaperiph)(unsafe.Pointer(addr))
}

// snum returns stream number - 1, eg. 0 for fisrt stream.
func snum(ch *Channel) uintptr {
	off := uintptr(unsafe.Pointer(ch)) & 0x3ff
	step := unsafe.Sizeof(chanregs{})
	return (off - unsafe.Sizeof(dma.DMA_Periph{})) / step
}

const (
	trce = 1 << 1
	htce = 1 << 2

	trerr = 1 << 3
	fferr = 0
	dmerr = 0
)

func (ch *Channel) status() byte {
	isr := sdma(ch).raw.ISR.U32.Load()
	return byte(isr >> (snum(ch) * 4) & 0xf)
}

func (ch *Channel) clear(flags byte) {
	mask := uint32(flags&0xf) << (snum(ch) * 4)
	sdma(ch).raw.IFCR.U32.Store(mask)
}

func (ch *Channel) enable() {
	ch.raw.EN().Set()
}

func (ch *Channel) disable() {
	ch.raw.EN().Clear()
}

func (ch *Channel) enabled() bool {
	return ch.raw.EN().Load() != 0
}

func (ch *Channel) irqEnabled() byte {
	return byte(ch.raw.CCR.U32.Load() & 0xe)
}

func (ch *Channel) enableIRQ(flags byte) {
	ch.raw.CCR.U32.SetBits(uint32(flags) & 0xe)
}

func (ch *Channel) disableIRQ(flags byte) {
	ch.raw.CCR.U32.ClearBits(uint32(flags) & 0xe)
}

const (
	mtp = 1 << dma.DIRn
	mtm = 1 << dma.MEM2MEMn

	circ = 1 << dma.CIRCn
	incP = 1 << dma.PINCn
	incM = 1 << dma.MINCn

	fifo_1_4 = 0
	fifo_2_4 = 0
	fifo_3_4 = 0
	fifo_4_4 = 0

	pfc = 0
)

func (ch *Channel) setup(m Mode) {
	mask := dma.DIR | dma.MEM2MEM | dma.CIRC | dma.PINC | dma.MINC | dma.PL
	ch.raw.CCR.StoreBits(mask, dma.CCR(m))
}

const (
	prioM = 1
	prioH = 2
	prioV = 3
)

func (ch *Channel) setPrio(prio Prio) {
	ch.raw.PL().Store(dma.CCR(prio) << dma.PLn)
}

func (ch *Channel) prio() Prio {
	return Prio(ch.raw.PL().Load() >> dma.PLn)
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

func (ch *Channel) burst() (p, m int) { return 1, 1 }
func (ch *Channel) setBurst(p, m int) {}

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

func (ch *Channel) request() Request {
	n := snum(ch) * 4
	return Request(sdma(ch).cselr.Bits(0xf << n))
}

func (ch *Channel) setRequest(req Request) {
	n := snum(ch) * 4
	sdma(ch).cselr.AtomicStoreBits(0xf<<n, uint32(req)<<n)
}
