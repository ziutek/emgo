// Package clock provides interface to manage nRF51 clocks source/generation.
package clock

import (
	"arch/cortexm/nvic"
	"mmio"
	"unsafe"

	"nrf5/hal/internal/mmap"
	"nrf5/hal/te"
)

type regs struct {
	te.Regs

	_            [2]mmio.U32
	hfclkrun     mmio.U32
	hfclkstat    mmio.U32
	_            mmio.U32
	lfclkrun     mmio.U32
	lfclkstat    mmio.U32
	lfclksrccopy mmio.U32
	_            [62]mmio.U32
	lfclksrc     mmio.U32
	_            [7]mmio.U32
	ctiv         mmio.U32
	_            [5]mmio.U32
	xtalfreq     mmio.U32
	_            [2]mmio.U32
	traceconfig  mmio.U32
}

func r() *regs {
	return (*regs)(unsafe.Pointer(mmap.APB_BASE))
}

func NVIC() nvic.IRQ {
	return r().NVIC()
}

func IRQEnabled() te.EventMask {
	return r().IRQEnabled()
}

func EnableIRQ(mask te.EventMask) {
	r().EnableIRQ(mask)
}

func DisableIRQ(mask te.EventMask) {
	r().DisableIRQ(mask)
}
