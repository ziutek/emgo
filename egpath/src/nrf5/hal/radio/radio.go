package radio

import (
	"mmio"
	"unsafe"

	"nrf5/hal/internal"
	"nrf5/hal/te"
)

type Periph struct {
	te.Regs

	crcstatus   mmio.U32
	_           mmio.U32
	rxmatch     mmio.U32
	rxcrc       mmio.U32
	dai         mmio.U32
	_           [60]mmio.U32
	packetptr   mmio.U32
	frequency   mmio.U32
	txpower     mmio.U32
	mode        mmio.U32
	pcnf0       mmio.U32
	pcnf1       mmio.U32
	base0       mmio.U32
	base1       mmio.U32
	prefix0     mmio.U32
	prefix1     mmio.U32
	txaddress   mmio.U32
	rxaddresses mmio.U32
	crccnf      mmio.U32
	crcpoly     mmio.U32
	crcinit     mmio.U32
	test        mmio.U32
	tifs        mmio.U32
	rssisample  mmio.U32
	_           mmio.U32
	state       mmio.U32
	datawhiteiv mmio.U32
	_           [2]mmio.U32
	bcc         mmio.U32
	_           [39]mmio.U32
	dab         [8]mmio.U32
	dap         [8]mmio.U32
	dacnf       mmio.U32
	_           [56]mmio.U32
	override    [5]mmio.U32
	_           [561]mmio.U32
	power       mmio.U32
}

//emgo:const
var RADIO = (*Periph)(unsafe.Pointer(internal.BaseAPB + 0x01000))

type Task byte

const (
	TXEN      Task = 0 // Enable radio in TX mode.
	RXEN      Task = 1 // Enable radio in RX mode.
	START     Task = 2 // Start radio.
	STOP      Task = 3 // Stop radio.
	DISABLE   Task = 4 // Disable radio.
	RSSISTART Task = 5 // Start measurement and take single sample of the RSSI.
	RSSISTOP  Task = 6 // Stop the RSSI measurement.
	BCSTART   Task = 7 // Start bit counter.
	BCSTOP    Task = 8 // Stop bit counter.
)

type Event byte

const (
	READY    Event = 0  // Ready event.
	ADDRESS  Event = 1  // Address event.
	PAYLOAD  Event = 2  // Payload event.
	END      Event = 3  // End event.
	DISABLED Event = 4  // Disabled event.
	DEVMATCH Event = 5  // An address match occurred on the last received pkt.
	DEVMISS  Event = 6  // No address match occurred on the last received pkt.
	RSSIEND  Event = 7  // A new RSSI sample is ready in RSSISAMPLE register.
	BCMATCH  Event = 10 // Bit counter reached bit count value specified in BCC.
)

type Shorts uint32

const (
	READY_START       Shorts = 1 << 0
	END_DISABLE       Shorts = 1 << 1
	DISABLED_TXEN     Shorts = 1 << 2
	DISABLED_RXEN     Shorts = 1 << 3
	ADDRESS_RSSISTART Shorts = 1 << 4
	END_START         Shorts = 1 << 5
	ADDRESS_BCSTART   Shorts = 1 << 6
	DISABLED_RSSISTOP Shorts = 1 << 8
)

func (p *Periph) Task(t Task) *te.Task      { return p.Regs.Task(int(t)) }
func (p *Periph) Event(e Event) *te.Event   { return p.Regs.Event(int(e)) }
func (p *Periph) Shorts(s Shorts) mmio.UM32 { return p.Regs.Shorts(uint32(s)) }

// CRCSTATUS returns CRC status of packet received.
func (p *Periph) CRCSTATUS() bool {
	return p.crcstatus.Bits(1) != 0
}

// RXMATCH returns logical address on which previous packet was received.
func (p *Periph) RXMATCH() int {
	return int(p.rxmatch.Bits(7))
}

// RXCRC returns CRC field of previously received packet.
func (p *Periph) RXCRC() uint32 {
	return p.rxcrc.Bits(0xffffff)
}

// DAI returns index(n) of device address, see DAB[n] and DAP[n], that got an
// address match.
func (p *Periph) DAI() int {
	return int(p.dai.Bits(7))
}

// PACKETPTR returns packet address to be used for the next transmission or
// reception. When transmitting, the packet pointed to by this address will be
// transmitted and when receiving, the received packet will be written to this
// address.
func (p *Periph) PACKETPTR() uintptr {
	return uintptr(p.packetptr.Load())
}

// SetPACKETPTR sets PACKETPTR. See PACKETPTR for more information.
func (p *Periph) SetPACKETPTR(addr uintptr) {
	p.packetptr.Store(uint32(addr))
}

type Freq uint16

// MakeFreq returns Freq for given freqMHz and low. nRF51 requires low=false.
func MakeFreq(freqMHz int, low bool) Freq {
	var f Freq
	if low {
		f = 0x100 | Freq(freqMHz-2360)&0x7f
	} else {
		f = Freq(freqMHz-2400) & 0x7f
	}
	return f
}

func (f Freq) Low() bool {
	return f&0x100 != 0
}

func (f Freq) MHz() int {
	mhz := int(f) & 0x7f
	if f.Low() {
		mhz += 2360
	} else {
		mhz += 2400
	}
	return mhz
}

// FREQUENCY returns radio channel frequency.
func (p *Periph) FREQUENCY() Freq {
	return Freq(p.frequency.Load())
}

// SetFREQUENCY sets FREQUENCY. See FREQUENCY for more information.
func (p *Periph) SetFREQUENCY(f Freq) {
	p.frequency.Store(uint32(f))
}

// TXPOWER returns RADIO output power in dBm.
func (p *Periph) TXPOWER() int {
	return int(int8(p.txpower.Load()))
}

// SetTXPOWER sets TXPOWER. See TXPOWER for more information.
func (p *Periph) SetTXPOWER(dBm int) {
	p.txpower.StoreBits(0xff, uint32(dBm))
}

type Mode byte

const (
	NRF_1M   = 0
	NRF_2M   = 1
	NRF_250K = 2
	BLE_1M   = 3
	BLE_2M   = 4
)

// MODE returns radio data rate and modulation setting. The radio supports
// Frequency-shift Keying (FSK) modulation, which depending on setting are
// compatible with either Nordic Semiconductor’s proprietary radios, or
// Bluetooth low energy.
func (p *Periph) MODE() Mode {
	return Mode(p.mode.Bits(0xf))
}

// SetMODE sets MODE. See MODE for more information.
func (p *Periph) SetMODE(mode Mode) {
	p.mode.StoreBits(0xf, uint32(mode))
}

type Config0 uint32

func MakeConfig0(prLen, s0Len, lenLen, s1Len int, s1AlwaysRAM bool) Config0 {
	c := uint32(prLen)&16<<20 | uint32(s1Len)&15<<16 | uint32(s0Len)&8<<5 |
		uint32(lenLen)&15
	if s1AlwaysRAM {
		c |= 1 << 20
	}
	return Config0(c)
}

// PrLen returns number of bits used for preamble.
func (c Config0) PrLen() int {
	return int(c >> 20 & 16)
}

// S0Len returns number of bits used for S0 field.
func (c Config0) S0Len() int {
	return int(c >> 5 & 8)
}

// LenLen returns number of bits used for LENGTH field.
func (c Config0) LenLen() int {
	return int(c & 15)
}

// S1Len returns number of bits used for S1 field.
func (c Config0) S1Len() int {
	return int(c >> 16 & 15)
}

// S1AlwaysRAM reports whether S1 is always included in RAM independent of S1Len
func (c Config0) S1AlwaysRAM() bool {
	return c>>20&1 != 0
}

func (p *Periph) PCNF0() Config0 {
	return Config0(p.pcnf0.Load())
}
