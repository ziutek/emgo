package gpio

import (
	"mmio"
	"unsafe"

	"arch/cortexm/bitband"

	"stm32/hal/raw/mmap"
)

func portnum(p *Port) int {
	return int(uintptr(unsafe.Pointer(p))-mmap.GPIOA_BASE) / 0x400
}

func bit(p *Port, reg *mmio.U32, portAbitn int) bitband.Bit {
	return bitband.Alias32(reg).Bit(portAbitn + portnum(p))
}
