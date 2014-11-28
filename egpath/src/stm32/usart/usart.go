// Package usart supports USART/UART devices in most STM32 MCUs.
package usart

type Dev struct {
	s   uint32 `C:"volatile"`
	d   uint32 `C:"volatile"`
	br  uint32 `C:"volatile"`
	c1  uint32 `C:"volatile"`
	c2  uint32 `C:"volatile"`
	c3  uint32 `C:"volatile"`
	gtp uint32 `C:"volatile"`
}

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

func (u *Dev) SetBaudRate(br int, pclk uint) {
	div := uint32(pclk / uint(br))
	if uint(br) > pclk/16 {
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

	afterLastIRQ
)

func (u *Dev) EnabledIRQs() IRQ {
	return IRQ(u.c1>>4) & (afterLastIRQ - 1)
}

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

func (u *Dev) HalfDuplex() bool {
	return u.c3|(1<<3) != 0
}

func (u *Dev) SetHalfDuplex(hd bool) {
	if hd {
		u.c3 |= 1 << 3
	} else {
		u.c3 &^= 1 << 3
	}
}

func (u *Dev) Store(d uint32) {
	u.d = d
}

func (u *Dev) Load() uint32 {
	return u.d
}

// STM32F10xxx doesn't support:
// 1. 8x oversampling mode
// 2. Onebit sample method.
