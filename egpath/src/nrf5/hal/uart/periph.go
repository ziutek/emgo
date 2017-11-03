package uart

import (
	"bits"
	"mmio"
	"unsafe"

	"nrf5/hal/gpio"
	"nrf5/hal/te"

	"nrf5/hal/internal/mmap"
)

type Periph struct {
	te.Regs

	_        [32]mmio.U32
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
var UART0 = (*Periph)(unsafe.Pointer(mmap.APB_BASE + 0x02000))

type Task byte

const (
	STARTRX Task = 0 // Start UART receiver.
	STOPRX  Task = 1 // Stop UART receiver.
	STARTTX Task = 2 // Start UART transmitter.
	STOPTX  Task = 3 // Stop UART transmitter.
	SUSPEND Task = 7 // Suspend UART. (nRF52)
)

type Event byte

const (
	CTS    Event = 0  // CTS is activated (set low). nRF51v3+.
	NCTS   Event = 1  // CTS is deactivated (set high). nRF51v3+.
	RXDRDY Event = 2  // Data received in RXD.
	TXDRDY Event = 7  // Data sent from TXD.
	ERROR  Event = 9  // Error detected.
	RXTO   Event = 17 // Receiver timeout.
)

func (p *Periph) Task(t Task) *te.Task    { return p.Regs.Task(int(t)) }
func (p *Periph) Event(e Event) *te.Event { return p.Regs.Event(int(e)) }

type Shorts uint32

const (
	CTS_STARTRX Shorts = 1 << 3
	NCTS_STOPRX Shorts = 1 << 4
)

func (p *Periph) LoadSHORTS() Shorts   { return Shorts(p.Regs.LoadSHORTS()) }
func (p *Periph) StoreSHORTS(s Shorts) { p.Regs.StoreSHORTS(uint32(s)) }

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

// LoadERRORSRC returns error flags.
func (p *Periph) LoadERRORSRC() Error {
	return Error(p.errorsrc.Load())
}

// ClearERRORSRC clears specfied error flags.
func (p *Periph) ClearERRORSRC(e Error) {
	p.errorsrc.Store(uint32(e))
}

// LoadENABLE reports whether the UART peripheral is enabled.
func (p *Periph) LoadENABLE() bool {
	return p.enable.Load()&4 != 0
}

// StoreENABLE enables or disables UART peripheral.
func (p *Periph) StoreENABLE(en bool) {
	p.enable.Store(uint32(bits.One(en)) << 2)
}

type Signal byte

const (
	SignalRTS Signal = 0
	SignalTXD Signal = 1
	SignalCTS Signal = 2
	SignalRXD Signal = 3
)

func (p *Periph) LoadPSEL(s Signal) gpio.Pin {
	return gpio.SelPin(int8(p.psel[s].Load()))
}
func (p *Periph) StorePSEL(s Signal, pin gpio.Pin) {
	p.psel[s].Store(uint32(pin.Sel()))
}

func (p *Periph) LoadRXD() byte {
	return byte(p.rxd.Load())
}

func (p *Periph) StoreTXD(b byte) {
	p.txd.Store(uint32(b))
}

type Baudrate uint32

const (
	Baud1200   Baudrate = 0x0004F000 // Actual rate: 1205 baud.
	Baud2400   Baudrate = 0x0009D000 // Actual rate: 2396 baud.
	Baud4800   Baudrate = 0x0013B000 // Actual rate: 4808 baud.
	Baud9600   Baudrate = 0x00275000 // Actual rate: 9598 baud.
	Baud14400  Baudrate = 0x003B0000 // Actual rate: 14414 baud.
	Baud19200  Baudrate = 0x004EA000 // Actual rate: 19208 baud.
	Baud28800  Baudrate = 0x0075F000 // Actual rate: 28829 baud.
	Baud31250  Baudrate = 0x00800000
	Baud38400  Baudrate = 0x009D5000 // Actual rate: 38462 baud.
	Baud57600  Baudrate = 0x00EBF000 // Actual rate: 55944 baud.
	Baud76800  Baudrate = 0x013A9000 // Actual rate: 57602 baud.
	Baud115200 Baudrate = 0x01D7E000 // Actual rate: 115204 baud.
	Baud230400 Baudrate = 0x03AFB000 // Actual rate: 230393 baud.
	Baud250k   Baudrate = 0x04000000
	Baud460800 Baudrate = 0x075F7000
	Baud921600 Baudrate = 0x0EBEE000 // Actual rate: 921585 baud.
	Baud1M     Baudrate = 0x10000000
)

func BR(baud int) Baudrate {
	return (Baudrate(uint64(baud)<<32/16e6) + 0x800) & 0xFFFFF000
}

func (br Baudrate) Baud() int {
	return int((uint64(br)*16e6 + 1<<31) >> 32)
}

// LoadBAUDRATE returns configured baudrate.
func (p *Periph) LoadBAUDRATE() Baudrate {
	return Baudrate(p.baudrate.Load())
}

// StoreBAUDRATE stores baudrate.
func (p *Periph) StoreBAUDRATE(br Baudrate) {
	p.baudrate.Store(uint32(br))
}
