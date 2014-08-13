package usart

import (
	"unsafe"

	"stm32/f4/setup"
)

type Dev struct {
	s   uint32 `C:"volatile"`
	d   uint32 `C:"volatile"`
	br  uint32 `C:"volatile"`
	c1  uint32 `C:"volatile"`
	c2  uint32 `C:"volatile"`
	c3  uint32 `C:"volatile"`
	gtp uint32 `C:"volatile"`
}

var (
	USART1 = (*Dev)(unsafe.Pointer(uintptr(0x40011000)))
	USART2 = (*Dev)(unsafe.Pointer(uintptr(0x40004400)))
	USART3 = (*Dev)(unsafe.Pointer(uintptr(0x40004800)))
	UART4  = (*Dev)(unsafe.Pointer(uintptr(0x40004C00)))
	UART5  = (*Dev)(unsafe.Pointer(uintptr(0x40005000)))
	USART6 = (*Dev)(unsafe.Pointer(uintptr(0x40011400)))
	UART7  = (*Dev)(unsafe.Pointer(uintptr(0x40007800)))
	UART8  = (*Dev)(unsafe.Pointer(uintptr(0x40007C00)))
)

type Status uint32

const (
	ParityErr Status = 1 << iota
	FramingErr
	Noise
	OverrunErr
	Idle
	RxNotEmpty
	TxDone
	TxEmpty

	LINBreak
	CTS
)

func (u *Dev) Status() Status {
	return Status(u.s)
}

func (u *Dev) SetBaudRate(br int) {
	div := uint32(setup.APB1Clk / uint(br))
	if uint(br) > setup.APB1Clk/16 {
		// Oversampling = 8
		u.c1 |= 1 << 15
		u.br = div&7<<1 | div&7
	} else {
		// Oversampling = 16
		u.c1 &^= 1 << 15
		u.br = div
	}

}

func (u *Dev) Enable() {
	u.c1 |= 1 << 13
}

func (u *Dev) Disable() {
	u.c1 &^= 1 << 13
}

type WordLen byte

const (
	Bits8 WordLen = iota
	Bits9
)

func (u *Dev) SetWordLen(l WordLen) {
	if l == Bits8 {
		u.c1 &^= 1 << 12
	} else {
		u.c1 |= 1 << 12
	}
}

type Parity byte

const (
	None Parity = iota
	_
	Even
	Odd
)

func (u *Dev) SetParity(p Parity) {
	u.c1 = u.c1&^(3<<9) | uint32(p)<<9
}

type IRQ byte

const (
	IdleIRQ IRQ = 1 << iota
	RxNotEmptyIRQ
	TxDoneIRQ
	TxEmptyIRQ
	ParityErrIRQ
)

func (u *Dev) EnableIRQs(irqs IRQ) {
	u.c1 |= uint32(irqs) << 4
}

func (u *Dev) DisableIRQs(irqs IRQ) {

	u.c1 &^= uint32(irqs) << 4
}

type Mode byte

const (
	Rx = 1 << 2
	Tx = 1 << 3
)

func (u *Dev) Mode() Mode {
	return Mode(u.c1 & (3 << 2))
}

func (u *Dev) SetMode(m Mode) {
	u.c1 = u.c1&^(3<<2) | uint32(m)
}

type StopBits byte

const (
	Stop1b StopBits = iota
	Stop0b5
	Stop2b
	Stop1b5
)

func (u *Dev) SetStopBits(sb StopBits) {
	u.c2 = u.c2&^(3<<12) | uint32(sb)<<12
}

func (u *Dev) Store(b byte) {
	u.d = uint32(b)
}

func (u *Dev) Load() byte {
	return byte(u.d)
}

func (u *Dev) Ready() (tx, rx bool) {
	s := u.Status()
	return s&TxEmpty != 0, s&RxNotEmpty != 0
}
