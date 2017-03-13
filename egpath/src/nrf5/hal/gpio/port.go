package gpio

import (
	"mmio"
	"unsafe"

	"nrf5/hal/internal"
)

type Port struct {
	_      mmio.U32
	out    mmio.U32
	outset mmio.U32
	outclr mmio.U32
	in     mmio.U32
	dir    mmio.U32
	dirset mmio.U32
	dirclr mmio.U32
	_      [120]mmio.U32
	pincnf [32]mmio.U32
}

//emgo:const
var P0 = (*Port)(unsafe.Pointer(internal.BaseAHB + 0x500))

// Pins is a bitmask which represents the pins of GPIO port.
type Pins uint32

const (
	Pin0 Pins = 1 << iota
	Pin1
	Pin2
	Pin3
	Pin4
	Pin5
	Pin6
	Pin7
	Pin8
	Pin9
	Pin10
	Pin11
	Pin12
	Pin13
	Pin14
	Pin15
	Pin16
	Pin17
	Pin18
	Pin19
	Pin20
	Pin21
	Pin22
	Pin23
	Pin24
	Pin25
	Pin26
	Pin27
	Pin28
	Pin29
	Pin30
	Pin31
)

func (p *Port) Pin(index int) Pin {
	ptr := uintptr(unsafe.Pointer(p))
	return Pin{ptr | uintptr(index&0x1f)}
}
