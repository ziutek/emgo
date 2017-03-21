package uart

import (
	"bits"
	"mmio"
	"unsafe"

	"nrf5/hal/gpio"
	"nrf5/hal/te"

	"nrf5/hal/internal/mmap"
	"nrf5/hal/internal/psel"
)

type Periph struct {
	te.Regs

	_        [31]mmio.U32
	errorsrc mmio.U32
	_        [31]mmio.U32
	enable   mmio.U32
	_        mmio.U32
	psel     [4]mmio.U32
	rxd      mmio.U32
	txd      mmio.U32
	_        mmio.U32
	baudrate mmio.U32
	_        [17]mmio.U32
	config   mmio.U32
}

//emgo:const
var UART0 = (*Periph)(unsafe.Pointer(mmap.BaseAPB + 0x02000))

type Task byte

const (
	STARTRX Task = 0  // Start UART receiver.
	STOPRX  Task = 1  // Stop UART receiver.
	STARTTX Task = 2  // Start UART transmitter.
	STOPTX  Task = 3  // Stop UART transmitter.
	SUSPEND Task = 19 // Suspend UART. (nRF52)
)

type Event byte

const (
	CTS    Event = 0  // CTS is activated (set low). Available: nRF51v3+.
	NCTS   Event = 1  // CTS is deactivated (set high). Available: nRF51v3+.
	RXDRDY Event = 2  // Data received in RXD.
	TXDRDY Event = 3  // Data sent from TXD.
	ERROR  Event = 11 // Error detected.
	RXTO   Event = 43 // Receiver timeout.
)

func (p *Periph) Task(t Task) *te.Task    { return p.Regs.Task(int(t)) }
func (p *Periph) Event(e Event) *te.Event { return p.Regs.Event(int(e)) }

type Shorts uint32

const (
	CTS_STARTRX Shorts = 1 << 3
	NCTS_STOPRX Shorts = 1 << 4
)

func (p *Periph) SHORTS() Shorts     { return Shorts(p.Regs.SHORTS()) }
func (p *Periph) SetSHORTS(s Shorts) { p.Regs.SetSHORTS(uint32(s)) }

// Error is a bitfield that describes detected errors.
type Error byte

const (
	ErrOverrun Error = 1 << 0
	ErrParity  Error = 1 << 1
	ErrFraming Error = 1 << 2
	ErrBreak   Error = 1 << 3
	ErrAll           = ErrOverrun | ErrParity | ErrFraming | ErrBreak
)

func (e Error) Error() string {
	var (
		s string
		d Error
	)
	switch {
	case e&ErrOverrun != 0:
		d = ErrOverrun
		s = "UART overrun+"
	case e&ErrFraming != 0:
		d = ErrFraming
		s = "UART framing+"
	case e&ErrParity != 0:
		d = ErrParity
		s = "UART parity+"
	case e&ErrBreak != 0:
		d = ErrBreak
		s = "UART break+"
	default:
		return ""
	}
	if e&^d == 0 {
		s = s[:len(s)-1]
	}
	return s
}

// ERRORSRC returns error source.
func (p *Periph) ERRORSRC() Error {
	return Error(p.errorsrc.Load())
}

// SetERRORSRC sets value of ERRORSRC. Usually used to reset all error flags:
// p.SetERRORSRC(0).
func (p *Periph) SetERRORSRC(e Error) {
	p.errorsrc.Store(uint32(e))
}

// ENABLE reports whether the p UART peripheral is enabled.
func (p *Periph) ENABLE() bool {
	return p.enable.Load()&4 != 0
}

func (p *Periph) SetENABLE(en bool) {
	p.enable.Store(uint32(bits.One(en)) << 2)
}

type Signal byte

const (
	SignalRTS Signal = 0
	SignalTXD Signal = 1
	SignalCTS Signal = 2
	SignalRXD Signal = 3
)

func (p *Periph) PSEL(s Signal) gpio.Pin {
	return psel.Pin(p.psel[s].Load())
}
func (p *Periph) SetPSEL(s Signal, pin gpio.Pin) {
	p.psel[s].Store(psel.Sel(pin))
}

func (p *Periph) RXD() byte {
	return byte(p.rxd.Load())
}

func (p *Periph) SetTXD(b byte) {
	p.txd.Store(uint32(b))
}

func (p *Periph) SetBAUDRATE(br uint32) {
	p.baudrate.Store(br)
}
