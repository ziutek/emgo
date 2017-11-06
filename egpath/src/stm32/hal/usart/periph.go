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

// Bus returns a bus to which p is connected to.
func (p *Periph) Bus() system.Bus {
	return internal.Bus(unsafe.Pointer(p))
}

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

// Event is bitmask that describes events in USART peripheral.
type Event byte

const (
	Idle       = Event(usart.IDLE >> 4) // IDLE line detected.
	RxNotEmpty = Event(usart.RXNE >> 4) // Read data register not empty.
	TxDone     = Event(usart.TC >> 4)   // Transmission complete.
	TxEmpty    = Event(usart.TXE >> 4)  // Transmit data register empty.
	LINBreak   = Event(lbd >> 4)        // LIN break detection flag.
	CTS        = Event(usart.CTS >> 4)  // Change on CTS status line

	EvAll = Idle | RxNotEmpty | TxDone | TxEmpty | LINBreak | CTS
)

// Error is bitmask that describes errors that can be detected by USART hardware
// when receiving data.
type Error byte

const (
	ErrParity  = Error(usart.PE)  // Parity error.
	ErrFraming = Error(usart.FE)  // Framing error.
	ErrNoise   = Error(usart.NE)  // Noise error flag.
	ErrOverrun = Error(usart.ORE) // Overrun error.
)

func (e Error) Error() string {
	var (
		s string
		d Error
	)
	switch {
	case e&ErrOverrun != 0:
		d = ErrOverrun
		s = "USART overrun+"
	case e&ErrNoise != 0:
		d = ErrNoise
		s = "USART noise+"
	case e&ErrFraming != 0:
		d = ErrFraming
		s = "USART framing+"
	case e&ErrParity != 0:
		d = ErrParity
		s = "USART parity+"
	default:
		return ""
	}
	if e&^d == 0 {
		s = s[:len(s)-1]
	}
	return s
}

// Status return current status of p.
func (p *Periph) Status() (Event, Error) {
	return p.status()
}

// Clear clears events ev and errors err. For MCUs that have no USART_ICR
// register (before L0, L4, F7 series) only RxNotEmpty, TxDone, LINBreak and CTS
// events can be cleared this way. Other events can be cleared only by specific
// sequence of reading status register and read or write data register.
func (p *Periph) Clear(ev Event, err Error) {
	p.clear(ev, err)
}

// EnableIRQ enables generating of IRQ by events e.
func (p *Periph) EnableIRQ(e Event) {
	if cr1e := e & (Idle | RxNotEmpty | TxDone | TxEmpty); cr1e != 0 {
		p.raw.CR1.SetBits(usart.CR1_Bits(cr1e) << 4)
	}
	if e&LINBreak != 0 {
		p.raw.LBDIE().Set()
	}
	if e&CTS != 0 {
		p.raw.CTSIE().Set()
	}
}

// DisableIRQ disables generating of IRQ by events e.
func (p *Periph) DisableIRQ(e Event) {
	if cr1e := e & (Idle | RxNotEmpty | TxDone | TxEmpty); cr1e != 0 {
		p.raw.CR1.ClearBits(usart.CR1_Bits(cr1e) << 4)
	}
	if e&LINBreak != 0 {
		p.raw.LBDIE().Clear()
	}
	if e&CTS != 0 {
		p.raw.CTSIE().Clear()
	}
}

// EnableErrorIRQ enables generating of IRQ by ErrNoise, ErrOverrun, ErrFraming
// errors when DMA is used to handle incoming data.
func (p *Periph) EnableErrorIRQ() {
	p.raw.EIE().Set()
}

func (p *Periph) DisableErrorIRQ() {
	p.raw.EIE().Clear()
}

// SetBaudRate sets baud rate [sym/s]. APB1 and APB2 clock in stm32/hal/system
// package must be set properly before use this function.
func (p *Periph) SetBaudRate(baudrate int) {
	br := uint(baudrate)
	pclk := p.Bus().Clock()
	usartdiv := (pclk + br/2) / br
	if uint(br) > pclk/16 {
		// Oversampling = 8
		p.raw.OVER8().Set() // BUG: Not supported by F1xx.
		usartdiv = usartdiv&^7<<1 | usartdiv&7
	} else {
		// Oversampling = 16
		p.raw.OVER8().Clear()
	}
	p.raw.BRR.Store(usart.BRR_Bits(usartdiv))
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

func (p *Periph) Conf() Conf {
	return Conf(p.raw.CR1.Bits(usart.CR1_Bits(RxEna|TxEna|ParOdd))) |
		Conf(p.raw.CR2.Bits(usart.CR2_Bits(Stop1b5>>16))<<16)
}

func (p *Periph) SetConf(cfg Conf) {
	p.raw.CR1.StoreBits(usart.CR1_Bits(RxEna|TxEna|ParOdd), usart.CR1_Bits(cfg))
	p.raw.CR2.StoreBits(usart.CR2_Bits(Stop1b5>>16), usart.CR2_Bits(cfg>>16))
}

type Mode uint32

const (
	// HalfDuplex enables half-duplx operation.
	HalfDuplex Mode = 1 << (16 + 3)

	// OneBit sets sampling method to one bit and disables Noise error
	// detection.
	OneBit Mode = 1 << (16 + 11)
)

func (p *Periph) SetMode(mode Mode) {
	//p.raw.CR2.StoreBits(mask, usart.CR2_Bits(mode))
	p.raw.CR3.StoreBits(
		usart.CR3_Bits((HalfDuplex|OneBit)>>16),
		usart.CR3_Bits(mode>>16),
	)
}

func (p *Periph) Store(d int) {
	p.store(d)
}

func (p *Periph) Load() int {
	return p.load()
}
