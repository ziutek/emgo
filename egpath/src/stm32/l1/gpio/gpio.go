package gpio

import "unsafe"

// Common for STM32F4xx and STM32L1xx

// GPIO represents registers of one GPIO port
type Port struct {
	moder   uint32 `C:"volatile"`
	otyper  uint32 `C:"volatile"`
	ospeedr uint32 `C:"volatile"`
	pupdr   uint32 `C:"volatile"`
	idr     uint32 `C:"volatile"`
	odr     uint32 `C:"volatile"`
	bsrr    uint32 `C:"volatile"`
	lckr    uint32 `C:"volatile"`
	afrl    uint32 `C:"volatile"`
	afrh    uint32 `C:"volatile"`
}

const (
	base uintptr = 0x40020000
	step uintptr = 0x400
)

var (
	A = (*Port)(unsafe.Pointer(base + step*0))
	B = (*Port)(unsafe.Pointer(base + step*1))
	C = (*Port)(unsafe.Pointer(base + step*2))
	D = (*Port)(unsafe.Pointer(base + step*3))
	E = (*Port)(unsafe.Pointer(base + step*4))
	F = (*Port)(unsafe.Pointer(base + step*5))
	G = (*Port)(unsafe.Pointer(base + step*6))
	H = (*Port)(unsafe.Pointer(base + step*7))
	I = (*Port)(unsafe.Pointer(base + step*8))
	J = (*Port)(unsafe.Pointer(base + step*9))
	K = (*Port)(unsafe.Pointer(base + step*10))
)

type Mode byte

const (
	In Mode = iota
	Out
	Alt
	Analog
)

// Mode returns I/O mode for n-th bit
func (g *Port) Mode(n int) Mode {
	n *= 2
	return Mode(g.moder>>uint(n)) & 3
}

// SetMode sets I/O mode for n-th bit
func (g *Port) SetMode(n int, mode Mode) {
	n *= 2
	g.moder = g.moder&^(3<<uint(n)) | uint32(mode)<<uint(n)
}

type OutType byte

const (
	PushPullOut OutType = iota
	OpenDrainOut
)

// OutType returns current type of n-th output bit
func (g *Port) OutType(n int) OutType {
	return OutType(g.otyper>>uint(n)) & 1
}

// SetOuttype sets type for n-th output bit
func (g *Port) SetOutType(n int, ot OutType) {
	g.otyper = g.otyper&^(1<<uint(n)) | uint32(ot)<<uint(n)
}

type Speed byte

// OutSpeed return current speed for n-th output bit
func (g *Port) OutSpeed(n int) Speed {
	n *= 2
	return Speed(g.ospeedr>>uint(n)) & 3
}

// SetOutSpeed sets speed for n-th output bit
func (g *Port) SetOutSpeed(n int, speed Speed) {
	n *= 2
	g.ospeedr = g.ospeedr&^(3<<uint(n)) | uint32(speed)<<uint(n)
}

// SetBit sets n-th output bit to 1
func (g *Port) SetBit(n int) {
	g.bsrr = uint32(1) << uint(n)
}

// ClearBit sets n-th output bit to 0
func (g *Port) ClearBit(n int) {
	g.bsrr = uint32(0x10000) << uint(n)
}

// SetBits sets output bits on positions specified by bits to 1
func (g *Port) SetBits(bits uint16) {
	g.bsrr = uint32(bits)
}

// ClearBits sets output bits on positions specified by bits to 0
func (g *Port) ClearBits(bits uint16) {
	g.bsrr = uint32(bits) << 16
}

// SetBSRR sets whole BSRR register.
// High 16 bits in bssr specifies which bits should be 0.
// Low 16 bits in bss specifies which bits should be 1.
func (g *Port) SetBSRR(bsrr uint32) {
	g.bsrr = bsrr
}

func (g *Port) Write(bits uint16) {
	g.odr = uint32(bits)
}

func (g *Port) Read() uint16 {
	return uint16(g.idr)
}

func (g *Port) Bit(n int) bool {
	return g.idr&(uint32(1)<<uint(n)) != 0
}