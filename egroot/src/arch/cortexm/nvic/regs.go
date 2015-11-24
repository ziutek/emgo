package nvic

import (
	"mmio"
	"unsafe"
)

type u8w struct {
	r [8]uint32
} //c:volatile

func (r *u8w) SetBit(irq IRQ) {
	val := uint32(1) << (irq & 31)
	irq >>= 5
	r.r[irq] = val
}

func (r *u8w) Bit(irq IRQ) bool {
	mask := uint32(1) << (irq & 31)
	irq >>= 5
	return r.r[irq]&mask != 0
}

type u60w struct {
	r [60]uint32 // Use uint32 because Cortex-M0 supports only word access.
} //c:volatile

func (r *u60w) SetByte(irq IRQ, b byte) {
	shift := uint(irq&3) * 8
	val := uint32(b) << shift
	mask := uint32(0xff) << shift
	irq >>= 2
	r.r[irq] = r.r[irq]&^mask | val
}

func (r *u60w) Byte(irq IRQ) byte {
	shift := uint(irq&3) * 8
	return byte(r.r[irq>>2] >> shift)
}

type regs struct {
	ISER u8w
	_    [24]uint32
	ICER u8w
	_    [24]uint32
	ISPR u8w
	_    [24]uint32
	ICPR u8w
	_    [24]uint32
	IABR u8w
	_    [56]uint32
	IPR  u60w
}

var (
	r   = (*regs)(unsafe.Pointer(uintptr(0xe000e100)))
	sti = mmio.PtrU32(unsafe.Pointer(uintptr(0xe000ef00)))
)
