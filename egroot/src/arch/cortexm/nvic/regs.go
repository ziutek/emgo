package nvic

import (
	"mmio"
	"unsafe"
)

type u8w struct {
	r [8]uint32
} //c:volatile

func (r *u8w) SetBit(e Exce) {
	val := uint32(1) << uint(e&31)
	e >>= 5
	r.r[e] = val
}

func (r *u8w) Bit(e Exce) bool {
	mask := uint32(1) << uint(e&31)
	e >>= 5
	return r.r[e]&mask != 0
}

type u60w struct {
	r [60]uint32 // Use uint32 because Cortex-M0 supports only word access.
} //c:volatile

func (r *u60w) SetByte(e Exce, b byte) {
	shift := uint(e&3) * 8
	val := uint32(b) << shift
	mask := uint32(0xff) << shift
	e >>= 2
	r.r[e] = r.r[e]&^mask | val
}

func (r *u60w) Byte(e Exce) byte {
	shift := uint(e&3) * 8
	return byte(r.r[e>>2] >> shift)
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
