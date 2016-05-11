package usart

import (
	"unsafe"

	"stm32/hal/internal"
	"stm32/hal/system"

	"stm32/hal/raw/usart"
)

// Periph represents USART peripheral.
type Periph struct {
	raw usart.USART_Periph
}

// Bus returns a bus to which p is connected.
func (p *Periph) Bus() system.Bus {
	return internal.Bus(unsafe.Pointer(p))
}

type Status uint16

const (
	ParityErr  = Status(usart.PE)   // Parity error.
	FramingErr = Status(usart.FE)   // Framing error.
	NoiseErr   = Status(usart.NE)   // Noise error flag.
	OverrunErr = Status(usart.ORE)  // Overrun error.
	Idle       = Status(usart.IDLE) // IDLE line detected.
	RxNotEmpty = Status(usart.RXNE) // Read data register not empty.
	TxDone     = Status(usart.TC)   // Transmission complete.
	TxEmpty    = Status(usart.TXE)  // Transmit data register empty.
	LINBreak   = Status(usart.LBD)  // LIN break detection flag.
	CTS        = Status(usart.CTS)  // CTS flag.
)

// EnableClock enables clock for p.
// lp determines whether the clock remains on in low power (sleep) mode.
func (p *Periph) EnableClock(lp bool) {
	addr := unsafe.Pointer(p)
	internal.APB_SetLPEnabled(addr, lp)
	internal.APB_SetEnabled(addr, true)
}

// DisableClock disables clock for p.
func (p *Periph) DisableClock() {
	internal.APB_SetEnabled(unsafe.Pointer(p), false)
}

// Reset resets p.
func (p *Periph) Reset() {
	internal.APB_Reset(unsafe.Pointer(p))
}

// Status return current status.
func (p *Periph) Status() Status {
	return Status(p.raw.SR.Load())
}

// SetBaudRate sets baudrate [sym/s]. APB1 and APB2 clock in stm32/hal/system
// package must be set properly before use this function.
func (p *Periph) SetBaudRate(baudrate int) {
	br := uint(baudrate)
	pclk := p.Bus().Clock()
	usartdiv := (pclk + br/2) / br
	if uint(br) > pclk/16 {
		// Oversampling = 8
		p.raw.OVER8().Set()
		usartdiv = usartdiv&^7<<1 | usartdiv&7
	} else {
		// Oversampling = 16
		p.raw.OVER8().Clear()
	}
	p.raw.BRR.U16.Store(uint16(usartdiv))
}

// Enable enables p.
func (p *Periph) Enable() {
	p.raw.UE().Set()
}

// Disable disables p at end of the current byte transfer.
func (p *Periph) Disable() {
	p.raw.UE().Clear()
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

func (p *Periph) SetConf(cfg Conf) {
	mask := uint16(RxEna | TxEna | ParOdd)
	p.raw.CR1.U16.StoreBits(mask, uint16(cfg))
	cfg >>= 16
	mask = uint16(Stop1b5 >> 16)
	p.raw.CR2.U16.StoreBits(mask, uint16(cfg))
}

type Mode uint32

const (
	HalfDuplex Mode = 1 << (16 + 3)
)

func (p *Periph) SetMode(mode Mode) {
	//mask :=
	//p.raw.CR2.U16.StoreBits(mask, uint16(mode))
	mode >>= 16
	mask := uint16(HalfDuplex >> 16)
	p.raw.CR3.U16.StoreBits(mask, uint16(mask))
}

type IRQ uint16

const (
	IdleIRQ       = IRQ(usart.IDLEIE)
	RxNotEmptyIRQ = IRQ(usart.RXNEIE)
	TxDoneIRQ     = IRQ(usart.TCIE)
	TxEmptyIRQ    = IRQ(usart.TXEIE)
	ParityErrIRQ  = IRQ(usart.PEIE)
)

func (p *Periph) IRQEnabled() IRQ {
	irqmask := IdleIRQ | RxNotEmptyIRQ | TxDoneIRQ | TxEmptyIRQ | ParityErrIRQ
	return IRQ(p.raw.CR1.U16.Bits(uint16(irqmask)))
}

func (p *Periph) EnableIRQ(irq IRQ) {
	p.raw.CR1.U16.SetBits(uint16(irq))
}

func (p *Periph) DisableIRQ(irq IRQ) {
	p.raw.CR1.U16.ClearBits(uint16(irq))
}

func (p *Periph) Store(d int) {
	p.raw.DR.U16.Store(uint16(d))
}

func (p *Periph) Load() int {
	return int(p.raw.DR.Load())
}
