// +build f10x_ld f10x_ld_vl f10x_md f10x_md_vl f10x_hd f10x_hd_vl f10x_xl f10x_cl l1xx_md l1xx_mdp l1xx_hd l1xx_xl

package dma

import (
	"unsafe"

	"stm32/hal/raw/dma"
	"stm32/hal/raw/mmap"
	"stm32/hal/raw/rcc"
)

var (
	DMA1 = (*DMA)(unsafe.Pointer(mmap.DMA1_BASE))
	DMA2 = (*DMA)(unsafe.Pointer(mmap.DMA2_BASE))
)

type dmaregs struct {
	raw dma.DMA_Periph
	chs [7]struct {
		raw dma.DMA_Channel_Periph
		_   uint32
	}
}

func pnum(p *DMA) int {
	return int(uintptr(unsafe.Pointer(p))-mmap.AHBPERIPH_BASE) / 0x400
}

func enableClock(p *DMA, _ bool) {
	bit := bit(p, &rcc.RCC.AHBENR.U32)
	bit.Set()
	bit.Load() // RCC delay (workaround for silicon bugs).
}

func disableClock(p *DMA) {
	bit(p, &rcc.RCC.AHBENR.U32).Clear()
}

type chanregs struct {
	raw dma.DMA_Channel_Periph
}

func getChannel(p *DMA, n int) *Channel {
	pn := pnum(p)
	n--
	if pn == 0 && uint(n) > 6 || uint(n) > 4 {
		panic("dma: bad channel")
	}
	return (*Channel)(unsafe.Pointer(&p.chs[n].raw))
}

// cnum returns channel number - 1, eg. 0 for fisrt channel.
func cnum(c *Channel) uintptr {
	off := uintptr(unsafe.Pointer(c)) & 0x3ff
	step := unsafe.Sizeof(dmaregs{}.chs[0])
	return (off - unsafe.Sizeof(dma.DMA_Periph{})) / step
}

func cdma(c *Channel) *dma.DMA_Periph {
	addr := uintptr(unsafe.Pointer(c)) &^ 0x3ff
	return (*dma.DMA_Periph)(unsafe.Pointer(addr))
}

func events(c *Channel) Events {
	isr := cdma(c).ISR.U32.Load()
	return Events(isr >> (cnum(c) * 4) & 0xf)
}

func clearEvents(c *Channel, e Events) {
	mask := uint32(e&0xf) << (cnum(c) * 4)
	cdma(c).IFCR.U32.Store(mask)
}

func enable(c *Channel) {
	c.raw.EN().Set()
}

func disable(c *Channel) {
	c.raw.EN().Clear()
}

func intEnabled(c *Channel, e Events) bool {
	return c.raw.CCR.U32.Load()&uint32(e)&0xe != 0
}

func enableInt(c *Channel, e Events) {
	c.raw.CCR.U32.SetBits(uint32(e) & 0xe)
}

func disableInt(c *Channel, e Events) {
	c.raw.CCR.U32.ClearBits(uint32(e) & 0xe)
}

const modeMask = 0x70f0

func mode(c *Channel) Mode {
	return Mode(c.raw.CCR.Bits(modeMask))
}

func setMode(c *Channel, m Mode) {
	c.raw.CCR.U32.StoreBits(modeMask, uint32(m))
}

/*
	P8  = 0 << 8 // Peripheral word size: 8-bits.
	P16 = 1 << 8 // Peripheral word size: 16-bits.
	P32 = 2 << 8 // Peripheral word size: 32-bits.

	M8  = 0 << 10 // Memory word size: 8-bits.
	M16 = 1 << 10 // Memory word size: 16-bits.
	M32 = 2 << 10 // Memory word size: 32-bits.
*/

func wordSize(c *Channel) (p, m uintptr) {
	ccr := uintptr(c.raw.CCR.Load())
	p = 1 << (ccr >> 8 & 3)
	m = 1 << (ccr >> 10 & 3)
	return
}

func setWordSize(c *Channel, p, m uintptr) {
	ccr := p&6<<7 | m&6<<9
	c.raw.CCR.U32.StoreBits(0xf00, uint32(ccr))
}

func num(c *Channel) int {
	return int(c.raw.CNDTR.Load())
}

func setNum(c *Channel, n int) {
	c.raw.NDT().UM32.Store(uint32(n))
}

func setAddrP(c *Channel, a unsafe.Pointer) {
	c.raw.CPAR.U32.Store(uint32(uintptr(a)))
}

func setAddrM(c *Channel, a unsafe.Pointer) {
	c.raw.CMAR.U32.Store(uint32(uintptr(a)))
}
