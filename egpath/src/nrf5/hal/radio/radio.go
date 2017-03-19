package radio

import (
	"bits"
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
	_           [15]mmio.U32
	modecnf0    mmio.U32
	_           [40]mmio.U32
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

func (p *Periph) Task(t Task) *te.Task    { return p.Regs.Task(int(t)) }
func (p *Periph) Event(e Event) *te.Event { return p.Regs.Event(int(e)) }

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

func (p *Periph) SHORTS() Shorts     { return Shorts(p.Regs.SHORTS()) }
func (p *Periph) SetSHORTS(s Shorts) { p.Regs.SetSHORTS(uint32(s)) }

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
func (p *Periph) SetPACKETPTR(addr unsafe.Pointer) {
	p.packetptr.Store(uint32(uintptr(addr)))
}

type Freq uint32

const (
	CM2400_2500 Freq = 0      // nRF5x
	CM2360_2460 Freq = 1 << 8 // nRF52
)

func Channel(ch int) Freq {
	return Freq(ch & 0x7f)
}

func (f Freq) Channel() int {
	return int(f & 0x7f)
}

// FREQUENCY returns a bitmap that encodes current radio channel and channelmap.
func (p *Periph) FREQUENCY() Freq {
	return Freq(p.frequency.Load())
}

// SetFREQUENCY sets radio channel and channel map.
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

type PktConf0 uint32

const (
	Pre8b     PktConf0 = 0
	Pre16b    PktConf0 = 1 << 24
	S0_0b     PktConf0 = 0
	S0_8b     PktConf0 = 1 << 8
	S1AutoRAM PktConf0 = 0
	S1AlwsRAM PktConf0 = 1 << 20
)

func LenBitN(n int) PktConf0 {
	return PktConf0(n & 15)
}

func S1BitN(n int) PktConf0 {
	return PktConf0(n & 15 << 16)
}

// LenBitN returns number of bits used for LENGTH field.
func (c PktConf0) LenBitN() int {
	return int(c & 15)
}

// S1BitN returns number of bits used for S1 field.
func (c PktConf0) S1BitN() int {
	return int(c >> 16 & 15)
}

func (p *Periph) PCNF0() PktConf0 {
	return PktConf0(p.pcnf0.Load())
}

func (p *Periph) SetPCNF0(pcnf PktConf0) {
	p.pcnf0.Store(uint32(pcnf))
}

type PktConf1 uint32

const (
	LSBFirst PktConf1 = 0
	MSBFirst PktConf1 = 1 << 24
	WhiteDis PktConf1 = 0
	WhiteEna PktConf1 = 1 << 25
)

func MaxLen(n int) PktConf1 {
	return PktConf1(n & 0xff)
}

func StatLen(n int) PktConf1 {
	return PktConf1(n & 0xff << 8)
}

func BALen(n int) PktConf1 {
	return PktConf1(n & 7 << 16)
}

// MaxLen returns maximum length of packet payload in bytes.
func (c PktConf1) MaxLen() int {
	return int(c & 0xff)
}

// StatLen returns static length of payload in bytes.
func (c PktConf1) StatLen() int {
	return int(c >> 8 & 0xff)
}

// BALen returns number of bytes used as base address.
func (c PktConf1) BALen() int {
	return int(c >> 16 & 7)
}

func (p *Periph) PCNF1() PktConf1 {
	return PktConf1(p.pcnf1.Load())
}

func (p *Periph) SetPCNF1(pcnf PktConf1) {
	p.pcnf1.Store(uint32(pcnf))
}

// BASE0 returns radio base address 0.
func (p *Periph) BASE0() uint32 {
	return p.base0.Load()
}

// SetBASE0 sets radio base address 0.
func (p *Periph) SetBASE0(ba uint32) {
	p.base0.Store(ba)
}

// BASE1 returns radio base address 1.
func (p *Periph) BASE1() uint32 {
	return p.base1.Load()
}

// SetBASE1 sets radio base address 1.
func (p *Periph) SetBASE1(ba uint32) {
	p.base1.Store(ba)
}

// PREFIX0 returns AP3<<24 | AP2<<16 | AP1<<8 | AP0.
func (p *Periph) PREFIX0() uint32 {
	return p.prefix0.Load()
}

// SetPREFIX0 sets PREFIX0 to prefix = AP3<<24 | AP2<<16 | AP1<<8 | AP0.
func (p *Periph) SetPREFIX0(prefix uint32) {
	p.prefix0.Store(prefix)
}

// PREFIX1 returns AP7<<24 | AP6<<16 | AP5<<8 | AP4.
func (p *Periph) PREFIX1() uint32 {
	return p.prefix1.Load()
}

// SetPREFIX1 sets PREFIX1 to prefix = AP7<<24 | AP6<<16 | AP5<<8 | AP4.
func (p *Periph) SetPREFIX1(prefix uint32) {
	p.prefix1.Store(prefix)
}

// TXADDRESS returns logical address used when transmitting a packet.
func (p *Periph) TXADDRESS() int {
	return int(p.txaddress.Bits(7))
}

// SetTXADDRESS sets logical address to be used when transmitting a packet.
func (p *Periph) SetTXADDRESS(laddr int) {
	p.txaddress.StoreBits(7, uint32(laddr))
}

// RXADERESSES returns bit field where eache of 8 low significant bits enables or
// disables one logical addresses for receive.
func (p *Periph) RXADDERESSES() uint32 {
	return p.rxaddresses.Load()
}

// SetRXADDERESSES sets bit field where eache of 8 low significant bits enables or
// disables one logical addresses for receive.
func (p *Periph) SetRXADDERESSES(laddr int) {
	p.rxaddresses.StoreBits(7, uint32(laddr))
}

// CRCCNF returns number of bytes in CRC field and whether address field is
// skipped for CRC calculation.
func (p *Periph) CRCCNF() (length int, skipAddr bool) {
	crccnf := p.crccnf.Load()
	return int(crccnf & 3), crccnf>>8&1 != 0
}

// SetCRCCNF sets number of bytes in CRC field and whether address will be
// skipped for CRC calculation.
func (p *Periph) SetCRCCNF(length int, skipAddr bool) {
	p.crccnf.Store(uint32(bits.One(skipAddr))<<8 | uint32(length)&3)
}

// CRCPOLY returns CRC polynominal coefficients.
func (p *Periph) CRCPOLY() uint32 {
	return p.crcpoly.Load() | 1
}

//  SetCRCPOLY sets CRC polynominal coefficients.
func (p *Periph) SetCRCPOLY(crcpoly uint32) {
	p.crcpoly.Store(crcpoly)
}

// CRCINIT returns initial value for CRC calculation.
func (p *Periph) CRCINIT() uint32 {
	return p.crcinit.Load()
}

//  SetCRCINIT sets initial value for CRC calculation.
func (p *Periph) SetCRCINIT(crcinit uint32) {
	p.crcinit.Store(crcinit)
}

func (p *Periph) TEST() (constCarrier, pllLock bool) {
	test := p.test.Load()
	return test&1 != 0, test&2 != 0
}

func (p *Periph) SetTEST(constCarrier, pllLock bool) {
	p.test.Store(uint32(bits.One(pllLock))<<1 | uint32(bits.One(constCarrier)))
}

// TIFS returns interframe spacing as the number of microseconds from the end of
// the last bit of the previous packet to the start of the first bit of the
// subsequent packet.
func (p *Periph) TIFS() int {
	return int(p.tifs.Load() & 255)
}

// SetTIFS sets interframe spacing as the number of microseconds from the end of
// the last bit of the previous packet to the start of the first bit of the
// subsequent packet.
func (p *Periph) SetTIFS(us int) {
	p.tifs.Store(uint32(us) & 255)
}

// RSSISAMPLE returns received signal strength [dBm].
func (p *Periph) RSSISAMPLE() int {
	return -int(p.rssisample.Load() & 127)
}

type State byte

const (
	Disabled  State = 0  // RADIO is in the Disabled state
	RxRu      State = 1  // RADIO is in the RXRU state
	RxIdle    State = 2  // RADIO is in the RXIDLE state
	Rx        State = 3  // RADIO is in the RX state
	RxDisable State = 4  // ADIO is in the RXDISABLED state
	TxRu      State = 9  // RADIO is in the TXRU state
	TxIdle    State = 10 // RADIO is in the TXIDLE state
	Tx        State = 11 // RADIO is in the TX state
	TxDisable State = 12 // RADIO is in the TXDISABLED state
)

//emgo:const
var stateStr = [...]string{
	Disabled:  "Disabled",
	RxRu:      "RxRu",
	RxIdle:    "RxIdle",
	Rx:        "Rx",
	RxDisable: "RxDisable",
	TxRu:      "TxRu",
	TxIdle:    "TxIdle",
	Tx:        "Tx",
	TxDisable: "TxDisable",
}

func (s State) String() string {
	var name string
	if int(s) < len(stateStr) {
		name = stateStr[s]
	}
	if len(name) == 0 {
		name = "unknown"
	}
	return name
}

// STATE returns current radio state.
func (p *Periph) STATE() State {
	return State(p.state.Bits(0xf))
}

// DATAWHITEIV returns data whitening initial value.
func (p *Periph) DATAWHITEIV() uint32 {
	return p.datawhiteiv.Load()
}

// SetDATAWHITEIV sets data whitening initial value.
func (p *Periph) SetDATAWHITEIV(initVal uint32) {
	p.datawhiteiv.Store(initVal)
}

// BCC returns value of bit counter compare.
func (p *Periph) BCC() int {
	return int(p.bcc.Load())
}

// SetBCC sets value of bit counter compare.
func (p *Periph) SetBCC(bcc int) {
	p.bcc.Store(uint32(bcc))
}

// POWER reports whether radio is powered on.
func (p *Periph) POWER() bool {
	return p.power.Bits(1) != 0
}

// DAB returns n-th device address base segment (32 LSBits of device address).
func (p *Periph) DAB(n int) uint32 {
	return p.dab[n].Load()
}

// SetDAB sets n-th device address base segment (32 LSBits of device address).
func (p *Periph) SetDAB(n int, dab uint32) {
	p.dab[n].Store(dab)
}

// DAP returns n-th device address prefix (16 MSBits of device address).
func (p *Periph) DAP(n int) uint16 {
	return uint16(p.dap[n].Load())
}

// SetDAP sets n-th device address prefix (16 MSBits of device address).
func (p *Periph) SetDAP(n int, dap uint16) {
	p.dap[n].Store(uint32(dap))
}

// DACNF returns a dap and txadd bit fields. Dap is a bit field where eache of 8 low
// significant bits enables or disables one device adressess (DAP-DAB pairs) for
// matching.
func (p *Periph) DACNF() (match, txadd byte) {
	dacnf := p.dacnf.Load()
	return byte(dacnf), byte(dacnf >> 8)
}

// SetDACNF sets bitmask that lists device adressess (DAP-DAB pairs) enabled
// for matching and TxAdd bits
func (p *Periph) SetDACNF(match, txadd byte) {
	p.dacnf.Store(uint32(txadd)<<8 | uint32(match))
}

type ModeConf0 uint32

const (
	FastRU   ModeConf0 = 1 << 0
	Tx1      ModeConf0 = 0
	Tx0      ModeConf0 = 1 << 8
	TxCenter ModeConf0 = 2 << 8
)

// MODECNF0 (nRF52 only).
func (p *Periph) MODECNF0() ModeConf0 {
	return ModeConf0(p.modecnf0.Load())
}

// SetMODECNF0 (nRF52 only).
func (p *Periph) SetMODECNF0(c ModeConf0) {
	p.modecnf0.Store(uint32(c))
}

// POWER reports whether radio is power on.
func (p *Periph) SPOWER() bool {
	return p.power.Bits(1) != 0
}

// SetPOWER can set peripheral power on or off.
func (p *Periph) SetPOWER(on bool) {
	p.power.StoreBits(1, uint32(bits.One(on)))
}
