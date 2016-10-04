package spi

import (
	"bits"
	"mmio"
	"unsafe"

	"stm32/hal/internal"
	"stm32/hal/system"

	"stm32/hal/raw/spi"
)

// Periph represents SPI peripheral.
type Periph struct {
	raw spi.SPI_Periph
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

type Conf uint16

const (
	CPHA0  = Conf(0)        // Sample on first edge.
	CPHA1  = Conf(spi.CPHA) // Sample on second edge.
	CPOL0  = Conf(0)        // Clock idle state is 0.
	CPOL1  = Conf(spi.CPOL) // Clock idle state is 1.
	Slave  = Conf(0)        // Slave mode.
	Master = Conf(spi.MSTR) // Master mode.

	BR2   = Conf(0)            // Baud rate = PCLK/2
	BR4   = Conf(1 << spi.BRn) // Baud rate = PCLK/4.
	BR8   = Conf(2 << spi.BRn) // Baud rate = PCLK/8.
	BR16  = Conf(3 << spi.BRn) // Baud rate = PCLK/16.
	BR32  = Conf(4 << spi.BRn) // Baud rate = PCLK/32.
	BR64  = Conf(5 << spi.BRn) // Baud rate = PCLK/64.
	BR128 = Conf(6 << spi.BRn) // Baud rate = PCLK/128.
	BR256 = Conf(7 << spi.BRn) // Baud rate = PCLK/256.

	MSBF = Conf(0)            // Most significant bit first.
	LSBF = Conf(spi.LSBFIRST) // Least significant bit first.

	HardSS = Conf(0)       // Hardware slave select.
	SoftSS = Conf(spi.SSM) // Software slave select (use ISSLow, ISSHigh).

	ISSLow  = Conf(0)       // Set NSS internally to low (requires SoftSS).
	ISSHigh = Conf(spi.SSI) // Set NSS internally to high (requires SoftSS).

	Frame8  = Conf(0)       // 8-bit frame (must disable Periph before change).
	Frame16 = Conf(spi.DFF) // 16-bit frame (must disable Periph before change).
)

// BR calculates baud rate bits of configuration. BR guarantees that returned
// value will set baud rate to the value closest to but not greater than
// baudrate. APB1 and APB2 clock in stm32/hal/system package must be set
// properly before use this function.
func (p *Periph) BR(baudrate int) Conf {
	pclk := p.Bus().Clock()
	div := pclk / uint(baudrate)
	if div < 2 {
		div = 2
	}
	br := 31 - bits.LeadingZeros32(uint32(div-1))
	return Conf(br << spi.BRn)
}

const cfgMask = ^uint16(spi.SPE)

// Conf returns the current configuration.
func (p *Periph) Conf() Conf {
	return Conf(p.raw.CR1.U16.Bits(cfgMask))
}

// SetConf configures p and enables or disables it.
func (p *Periph) SetConf(cfg Conf) {
	p.raw.CR1.U16.StoreBits(cfgMask, uint16(cfg))
}

type Event byte

const (
	RxNotEmpty = Event(spi.RXNE) // Receive buffer not empty.
	TxEmpty    = Event(spi.TXE)  // Transmit buffer empty.
	Busy       = Event(spi.BSY)  // Periph is busy (not a real event).

	eventMask = RxNotEmpty | TxEmpty | Busy
)

type Error byte

const (
	ErrUnderrun = Error(spi.UDR >> 3)
	ErrCRC      = Error(spi.CRCERR >> 3)
	ErrMode     = Error(spi.MODF >> 3)
	ErrOverrun  = Error(spi.OVR >> 3)

	errorMask = ErrUnderrun | ErrCRC | ErrMode | ErrOverrun
)

func (e Error) Error() string {
	var (
		s string
		d Error
	)
	switch {
	case e&ErrUnderrun != 0:
		d = ErrUnderrun
		s = "SPI underrun+"
	case e&ErrCRC != 0:
		d = ErrCRC
		s = "SPI CRC error+"
	case e&ErrMode != 0:
		d = ErrMode
		s = "SPI mode fault+"
	case e&ErrOverrun != 0:
		d = ErrOverrun
		s = "SPI overrun+"
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
	sr := p.raw.SR.Load()
	return Event(sr) & eventMask, Error(sr>>3) & errorMask
}

// Enable enables p.
func (p *Periph) Enable() {
	p.raw.SPE().Set()
}

// Disable disables p.
func (p *Periph) Disable() {
	p.raw.SPE().Clear()
}

func (p *Periph) StoreU16(v uint16) {
	p.raw.DR.U16.Store(v)
}

func (p *Periph) LoadU16() uint16 {
	return p.raw.DR.U16.Load()
}

func (p *Periph) StoreByte(v byte) {
	(*mmio.U8)(unsafe.Pointer(&p.raw.DR)).Store(v)
}

func (p *Periph) LoadByte() byte {
	return (*mmio.U8)(unsafe.Pointer(&p.raw.DR)).Load()
}
