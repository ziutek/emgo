package gpiote

import (
	"mmio"
	"unsafe"

	"arch/cortexm/nvic"

	"nrf5/hal/te"

	"nrf5/hal/internal/mmap"
)

type regs struct {
	te.Regs

	_      [68]mmio.U32
	config [8]mmio.U32
}

func r() *regs {
	return (*regs)(unsafe.Pointer(mmap.APB_BASE + 0x06000))
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
