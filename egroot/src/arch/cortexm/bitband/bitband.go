package bitband

import (
	"mmio"
	"unsafe"
)

type Bit struct {
	a *mmio.U32
}

func (b Bit) Load() int {
	return int(b.a.Load())
}

func (b Bit) Store(v int) {
	b.a.Store(uint32(v))
}

func (b Bit) Set() {
	b.Store(1)
}

func (b Bit) Clear() {
	b.Store(0)
}

// 0x20000000 - 0x200FFFFF: SRAM bit-band region.
// 0x22000000 - 0x23FFFFFF: SRAM bit-band alias.
//
// 0x40000000 - 0x400FFFFF: peripheral bit-band region.
// 0x42000000 - 0x43FFFFFF: peripheral bit-band alias.
func bitAlias(addr unsafe.Pointer) unsafe.Pointer {
	a := uintptr(addr)
	base := a &^ 0xfffff
	if base != 0x40000000 && base != 0x20000000 {
		panic("bitband: not in region")
	}
	base += 0x2000000
	offset := a & 0xfffff
	return unsafe.Pointer(base + offset*32)
}

type Bits8 struct {
	a *[8]mmio.U32
}

func (b Bits8) Bit(n int) Bit {
	return Bit{&b.a[n]}
}

func Alias8(r *mmio.U8) Bits8 {
	return Bits8{(*[8]mmio.U32)(bitAlias(unsafe.Pointer(r)))}
}

type Bits16 struct {
	a *[16]mmio.U32
}

func (b Bits16) Bit(n int) Bit {
	return Bit{&b.a[n]}
}

func Alias16(r *mmio.U16) Bits16 {
	return Bits16{(*[16]mmio.U32)(bitAlias(unsafe.Pointer(r)))}
}

type Bits32 struct {
	a *[32]mmio.U32
}

func (b Bits32) Bit(n int) Bit {
	return Bit{&b.a[n]}
}

func Alias32(r *mmio.U32) Bits32 {
	return Bits32{(*[32]mmio.U32)(bitAlias(unsafe.Pointer(r)))}
}
