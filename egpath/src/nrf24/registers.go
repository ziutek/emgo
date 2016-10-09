package nrf24

import (
	"fmt"
	"io"
	"strconv"
)

func fflags(f fmt.State, format string, mask, b byte) {
	m := byte(0x80)
	k := 0
	for i, c := range format {
		if c == '+' {
			for mask&m == 0 {
				m >>= 1
			}
			if b&m == 0 {
				io.WriteString(f, format[k:i])
				f.Write([]byte{'-'})
			} else {
				io.WriteString(f, format[k:i+1])
			}
			m >>= 1
			k = i + 1
		}
	}
	if k < len(format) {
		io.WriteString(f, format[k:])
	}
}

func (d *Device) byteReg(addr byte) byte {
	var buf [1]byte
	d.Reg(addr, buf[:])
	return buf[0]
}

type Status byte

const (
	FullTx Status = 1 << 0 // Tx FIFO full flag.
	MaxRT  Status = 1 << 4 // Maximum number of Tx retransmits interrupt.
	TxDS   Status = 1 << 5 // Data Sent Tx FIFO interrupt.
	RxDR   Status = 1 << 6 // Data Ready Rx FIFO interrupt.
)

// RxPipe returns data pipe number for the payload available for reading from
// RxFifo or -1 if RxFifo is empty
func (s Status) RxPipe() int {
	n := int(s) & 0x0e
	if n == 0x0e {
		return -1
	}
	return n >> 1
}

func (s Status) Format(f fmt.State, _ rune) {
	fflags(f, "RxDR+ TxDS+ MaxRT+ FullTx+ RxPipe:", 0x71, byte(s))
	strconv.WriteInt(f, s.RxPipe(), 10, 1)
}

type Config byte

const (
	PrimRx    Config = 1 << 0 //  Rx/Tx control 1: PRX, 0: PTX.
	PwrUp     Config = 1 << 1 // 1: power up, 0: power down.
	CRCO      Config = 1 << 2 // CRC encoding scheme 0: one byte, 1: two bytes.
	EnCRC     Config = 1 << 3 // Enable CRC. Force 1 if one of bits in AA is 1.
	MaskMaxRT Config = 1 << 4 // If 1 then mask interrupt caused by MaxRT.
	MaskTxDS  Config = 1 << 5 // If 1 then mask interrupt caused by TxDS.
	MaskRxDR  Config = 1 << 6 // If 1 then mask interrupt caused by RxDR.
)

func (c Config) Format(f fmt.State, _ rune) {
	fflags(
		f, "Mask(RxDR+ TxDS+ MaxRT+) EnCRC+ CRCO+ PwrUp+ PrimRx+",
		0x7f, byte(c),
	)
}

// Config returns value of CONFIG register.
func (d *Device) Config() Config {
	return Config(d.byteReg(0))
}

// SetConfig sets value of CONFIG register.
func (d *Device) SetConfig(c Config) {
	d.SetReg(0, byte(c))
}

/*
// Pipe is a bitfield that represents nRF24L01+ Rx data pipes.
type Pipe byte

const (
	P0 Pipe = 1 << iota
	P1
	P2
	P3
	P4
	P5
	PAll = P0 | P1 | P2 | P3 | P4 | P5
)

func (p Pipe) String() string {
	return flags("P5+ P4+ P3+ P2+ P1+ P0+", 0x3f, byte(p))
}

// AA returns value of EN_AA (Enable ‘Auto Acknowledgment’ Function) register.
func (d *Device) AA() Pipe {
	return Pipe(d.byteReg(1))
}

// SetAA sets value of EN_AA (Enable ‘Auto Acknowledgment’ Function) register.
func (d *Device) SetAA(p Pipe) {
	d.SetReg(1, byte(p))
}

// RxAE returns value of EN_RXADDR (Enabled RX Addresses) register.
func (d *Device) RxAE() Pipe {
	return Pipe(d.byteReg(2))
}

// SetRxAE sets value of EN_RXADDR (Enabled RX Addresses) register.
func (d *Device) SetRxAE(p Pipe) {
	d.SetReg(2, byte(p))
}

// AW returns value of SETUP_AW (Setup of Address Widths) register increased
// by two.
func (d *Device) AW() int {
	return int(d.byteReg(3)) + 2
}

// SetAW sets value of SETUP_AW (Setup of Address Widths) register to (alen-2).
func (d *Device) SetALen(alen int) {
	if alen < 3 || alen > 5 {
		panic("alen<3 || alen>5")
	}
	d.SetReg(3, byte(alen-2))
}

// Retr returns value of SETUP_RETR (Setup of Automatic Retransmission)
// register converted to number of retries and delay betwee retries.
func (d *Device) Retr() (cnt, dlyus int) {
	b := d.byteReg(4)
	cnt = int(b & 0xf)
	dlyus = (int(b>>4) + 1) * 250
	return
}

// SetRetr sets value of SETUP_RETR (Setup of Automatic Retransmission)
// register using cnt as number of retries and dlyus as delay betwee retries.
func (d *Device) SetRetr(cnt, dlyus int) {
	if uint(cnt) > 15 {
		panic("cnt<0 || cnt>15")
	}
	if dlyus < 250 || dlyus > 4000 {
		panic("dlyus<250 || dlyus>4000")
	}
	d.SetReg(4, byte((dlyus/250-1)<<4|cnt))
}

// Ch returns value of RF_CH (RF Channel) register.
func (d *Device) Ch() int {
	return int(d.byteReg(5))
}

// SetCh sets value of RF_CH (RF Channel) register.
func (d *Device) SetCh(ch int) {
	if uint(ch) > 127 {
		panic("ch<0 || ch>127")
	}
	d.SetReg(5, byte(ch))
}

type RF byte

const (
	LNAHC RF = 1 << iota // (nRF24L01.LNA_HCURR) Rx LNA gain 0: -1.5dB,-0.8mA.
	_
	_
	DRHigh // (RF_DR_HIGH) Select high speed data rate 0: 1Mbps, 1: 2Mbps.
	Lock   // (PLL_LOCK) Force PLL lock signal. Only used in test.
	DRLow  // (RF_DR_LOW) Set RF Data Rate to 250kbps.
	_
	Wave // (CONT_WAVE) Enables continuous carrier transmit when 1.
)

// Pwr returns RF output power in Tx mode [dBm].
func (rf RF) Pwr() int {
	return 3*int(rf&6) - 18
}

func Pwr(dbm int) RF {
	switch {
	case dbm < -18:
		dbm = -18
	case dbm > 0:
		dbm = 0
	}
	return RF((18+dbm)/3) & 6
}

func (rf RF) String() string {
	return flags("Wave+ DRLow+ Lock+ DRHigh+ LNAHC+ Pwr:", 0xb9, byte(rf)) +
		strconv.Itoa(rf.Pwr()) + "dBm"
}

// RF returns value of RF_SETUP register.
func (d *Device) RF() RF {
	return RF(d.byteReg(6))
}

// RF sets value of RF_SETUP register.
func (d *Device) SetRF(rf RF) {
	d.SetReg(6, byte(rf))
}

// Clear clears specified bits in STATUS register.
func (d *Device) Clear(stat Status) {
	d.SetReg(7, byte(stat))
}

// TxCnt returns values of PLOS and ARC counters from OBSERVE_TX register.
func (d *Device) TxCnt() (plos, arc int) {
	b := d.byteReg(8)
	arc = int(b & 0xf)
	plos = int(b >> 4)
	return
}

// RPD returns value of RPD (Received Power Detector) register (is RP > -64dBm).
// In case of nRF24L01 it returns value of.CD (Carrier Detect) register.
func (d *Device) RPD() bool {
	return d.byteReg(9)&1 != 0
}
*/

func checkPayNum(pn int) {
	if uint(pn) > 5 {
		panic("pn<0 || pn>5")
	}
}

func checkAddr(addr []byte) {
	if len(addr) > 5 {
		panic("len(addr)>5")
	}
}

func checkPayNumAddr(pn int, addr []byte) {
	checkPayNum(pn)
	checkAddr(addr)
	if pn > 1 && len(addr) > 1 {
		panic("pn>1 && len(addr)>1")
	}
}

// RxAddr reads address assigned to Rx pipe pn into addr.
func (d *Device) RxAddr(pn int, addr []byte) {
	checkPayNumAddr(pn, addr)
	d.Reg(byte(0xa+pn), addr)
}

// RxAddr0 returns least significant byte of address assigned to Rx pipe pn.
func (d *Device) RxAddr0(pn int) byte {
	checkPayNum(pn)
	return d.byteReg(byte(0xa + pn))
}

// SetRxAddr sets address assigned to Rx pipe pn to addr.
func (d *Device) SetRxAddr(pn int, addr ...byte) {
	checkPayNumAddr(pn, addr)
	d.SetReg(byte(0xa+pn), addr...)
}

// TxAddr reads value of TX_ADDR (Transmit address) into addr.
func (d *Device) TxAddr(addr []byte) {
	checkAddr(addr)
	d.Reg(0x10, addr)
}

// SetTxAddr sets value of TX_ADDR (Transmit address).
func (d *Device) SetTxAddr(addr ...byte) {
	checkAddr(addr)
	d.SetReg(0x10, addr...)
}

// RxPW returns Rx payload width set for pipe pn.
func (d *Device) RxPW(pn int) int {
	checkPayNum(pn)
	return int(d.byteReg(byte(0x11+pn))) & 0x3f
}

// SetRxPW sets Rx payload width for pipe pn.
func (d *Device) SetRxPW(pn, pw int) {
	checkPayNum(pn)
	if uint(pw) > 32 {
		panic("pw<0 || pw>32")
	}
	d.SetReg(byte(0x11+pn), byte(pw))
}

/*
type FIFO byte

const (
	RxEmpty FIFO = 1 << iota // 1: Rx FIFO empty, 0: Data in Rx FIFO.
	RxFull                   // 1: Rx FIFO full, 0: Avail.locations in Rx FIFO.
	_
	_
	TxEmpty // 1: Tx FIFO empty, 0: Data in Tx FIFO.
	TxFull  // 1: Tx FIFO full, 0: Available locations in Tx FIFO.
	TxReuse // 1: Reuse last transmitted payload.
)

func (f FIFO) String() string {
	return flags("TxReuse+ TxFull+ TxEmpty+ RxFull+ RxEmpty+", 0x73, byte(f))
}

// FIFO returns value of FIFO_STATUS register.
func (d *Device) FIFO() FIFO {
	return FIFO(d.byteReg(0x17))
}

// DynPD returns value of DYNPD (Enable dynamic payload length) register.
func (d *Device) DynPD() Pipe {
	return Pipe(d.byteReg(0x1c))
}

// SetDynPD sets value of DYNPD (Enable dynamic payload length) register.
func (d *Device) SetDynPD(p Pipe) {
	d.SetReg(0x1c, byte(p))
}

type Feature byte

const (
	DynAck Feature = 1 << iota // 1: Enables the W_TX_PAYLOAD_NOACK command.
	AckPay                     // 1: Enables payload with ACK
	DPL                        // 1: Enables dynamic payload length
)

func (f Feature) String() string {
	return flags("DPL+ AckPay+ DynAck+", 7, byte(f))
}

// Feature returns value of FEATURE register.
func (d *Device) Feature() Feature {
	return Feature(d.byteReg(0x1d))
}

// SetFeature sets value of FEATURE register.
func (d *Device) SetFeature(f Feature) {
	d.SetReg(0x1d, byte(f))
}
*/
