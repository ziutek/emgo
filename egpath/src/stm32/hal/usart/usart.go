package usart

import (
	"mmio"
	"unsafe"

	"stm32/hal/setup"
)

// USART represents USART device.
type USART struct {
	sr   mmio.U32
	dr   mmio.U32
	brr  mmio.U32
	cr1  mmio.U32
	cr2  mmio.U32
	cr3  mmio.U32
	gtpr mmio.U32
}

func (u *USART) BaseAddr() uintptr {
	return uintptr(unsafe.Pointer(u))
}

type Status uint16

const (
	ParityErr  Status = 1 << 0 // Parity error.
	FramingErr Status = 1 << 1 // Framing error.
	NoiseErr   Status = 1 << 2 // Noise error flag.
	OverrunErr Status = 1 << 3 // Overrun error.
	Idle       Status = 1 << 4 // IDLE line detected.
	RxNotEmpty Status = 1 << 5 // Read data register not empty.
	TxDone     Status = 1 << 6 // Transmission complete.
	TxEmpty    Status = 1 << 7 // Transmit data register empty.
	LINBreak   Status = 1 << 8 // LIN break detection flag.
	CTS        Status = 1 << 9 // CTS flag.
)

// EnableClock enables clock for USART.
// lp determines whether the clock remains on in low power (sleep) mode.
func (u *USART) EnableClock(lp bool) {
	setLPEnabled(u, lp)
	setEnabled(u, true)
}

// DisableClock disables clock for u.
func (u *USART) DisableClock() {
	setEnabled(u, false)
}

// Reset resets USART u.
func (u *USART) Reset() {
	reset(u)
}

// Status return current status.
func (u *USART) Status() Status {
	return Status(u.sr.Load())
}

// SetBaudRate sets baudrate [sym/s]. APB1Clk, APB2Clk in stm32/hal/setup package
// must be set properly before use this function.
func (u *USART) SetBaudRate(baudrate int) {
	br := uint(baudrate)
	pclk := setup.PeriphClk(u.BaseAddr())
	usartdiv := (pclk + br/2) / br
	const over8 = 1 << 15
	if uint(br) > pclk/16 {
		// Oversampling = 8
		u.cr1.SetBits(over8)
		usartdiv = usartdiv&^7<<1 | usartdiv&7
	} else {
		// Oversampling = 16
		u.cr1.ClearBits(over8)
	}
	u.brr.Store(uint32(usartdiv))
}

// Enable enables u.
func (u *USART) Enable() {
	u.cr1.SetBits(1 << 13)
}

// Disable disables u at end of the current byte transfer.
func (u *USART) Disable() {
	u.cr1.ClearBits(1 << 13)
}

type Conf uint32

const (
	RxEna   Conf = 1 << 2         // Receiver enabled.
	TxEna   Conf = 1 << 3         // Transmiter enabled.
	ParEven Conf = 2 << 9         // Parity control enabled: even.
	ParOdd  Conf = 3 << 9         // Parity control enabled: odd.
	Word9b  Conf = 1 << 12        // Use 9 bit word instead of 8 bit.
	Stop0b5 Conf = 1 << (16 + 12) // Use 0.5 stop bits insted of 1.
	Stop2b  Conf = 2 << (16 + 12) // Use 2 stop bits instead of 1.
	Stop1b5 Conf = 3 << (16 + 12) // Use 1.5 stop bits instead of 1.
)

func (u *USART) SetConf(cfg Conf) {
	mask := RxEna | TxEna | ParOdd
	u.cr1.StoreBits(uint32(mask), uint32(cfg))
	cfg >>= 16
	mask = Stop1b5 >> 16
	u.cr2.StoreBits(uint32(mask), uint32(cfg))
}

type Mode uint32

const (
	HalfDuplex Mode = 1 << (16 + 3)
)

func (u *USART) SetMode(mode Mode) {
	//mask :=
	//u.cr2.StoreBits(uint32(mask), uint32(mode))
	mode >>= 16
	mask := HalfDuplex >> 16
	u.cr3.StoreBits(uint32(mask), uint32(mask))
}

type IRQs uint16

const (
	IdleIRQ       IRQs = 1 << 4
	RxNotEmptyIRQ IRQs = 1 << 5
	TxDoneIRQ     IRQs = 1 << 6
	TxEmptyIRQ    IRQs = 1 << 7
	ParityErrIRQ  IRQs = 1 << 8
)

func (u *USART) IRQsEnabled() IRQs {
	const irqmask = IdleIRQ | RxNotEmptyIRQ | TxDoneIRQ | TxEmptyIRQ | ParityErrIRQ
	return IRQs(u.cr1.Bits(uint32(irqmask)))
}

func (u *USART) EnableIRQs(irqs IRQs) {
	u.cr1.SetBits(uint32(irqs))
}

func (u *USART) DisableIRQs(irqs IRQs) {
	u.cr1.ClearBits(uint32(irqs))
}

func (u *USART) Store(d int) {
	u.dr.Store(uint32(d))
}

func (u *USART) Load() int {
	return int(u.dr.Load())
}
