package nrf24

import (
	"fmt"
	"io"
	"strconv"
)

func fflags(fs fmt.State, format string, mask, b byte) {
	m := byte(0x80)
	k := 0
	for i := 0; i < len(format); i++ {
		c := format[i]
		if c == '+' {
			for mask&m == 0 {
				m >>= 1
			}
			if b&m == 0 {
				io.WriteString(fs, format[k:i])
				fs.Write([]byte{'-'})
			} else {
				io.WriteString(fs, format[k:i+1])
			}
			m >>= 1
			k = i + 1
		}
	}
	if k < len(format) {
		io.WriteString(fs, format[k:])
	}
}

func (r *Radio) byteReg(addr byte) (byte, Status) {
	var buf [1]byte
	s := r.Reg(addr, buf[:])
	return buf[0], s
}

type Status byte

const (
	FullTx Status = 1 << 0 // Tx FIFO full flag.
	MaxRT  Status = 1 << 4 // Maximum number of Tx retransmits interrupt.
	TxDS   Status = 1 << 5 // Data Sent Tx FIFO interrupt.
	RxDR   Status = 1 << 6 // Data Ready Rx FIFO interrupt.
)

// RxPipe returns the data pipe number for the payload available for reading
// from RxFifo or -1 if RxFifo is empty
func (s Status) RxPipe() int {
	n := int(s) & 0x0e
	if n == 0x0e {
		return -1
	}
	return n >> 1
}

func (s Status) Format(fs fmt.State, _ rune) {
	fflags(fs, "RxDR+ TxDS+ MaxRT+ FullTx+ RxPipe:", 0x71, byte(s))
	strconv.WriteInt(fs, s.RxPipe(), 10, 0)
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

func (c Config) Format(fs fmt.State, _ rune) {
	fflags(
		fs, "Mask(RxDR+ TxDS+ MaxRT+) EnCRC+ CRCO+ PwrUp+ PrimRx+",
		0x7f, byte(c),
	)
}

// Config returns the value of the CONFIG register.
func (d *Radio) Config() (Config, Status) {
	b, s := d.byteReg(0)
	return Config(b), s
}

// SetConfig sets the value of the CONFIG register.
func (r *Radio) SetConfig(c Config) Status {
	return r.SetReg(0, byte(c))
}

// Pipes is a bitfield that represents the nRF24L01+ Rx data pipes.
type Pipes byte

const (
	P0   Pipes = 1 << 0
	P1   Pipes = 1 << 1
	P2   Pipes = 1 << 2
	P3   Pipes = 1 << 3
	P4   Pipes = 1 << 4
	P5   Pipes = 1 << 5
	PAll       = P0 | P1 | P2 | P3 | P4 | P5
)

func (p Pipes) Format(fs fmt.State, _ rune) {
	fflags(fs, "P5+ P4+ P3+ P2+ P1+ P0+", 0x3f, byte(p))
}

// AA returns the value of the EN_AA (Enable ‘Auto Acknowledgment’ Function)
// register.
func (r *Radio) AA() (Pipes, Status) {
	b, s := r.byteReg(1)
	return Pipes(b), s
}

// SetAA sets the value of the EN_AA (Enable ‘Auto Acknowledgment’ Function)
// register.
func (r *Radio) SetAA(p Pipes) Status {
	return r.SetReg(1, byte(p))
}

// RxAEn returns the value of the EN_RXADDR (Enabled RX Addresses) register.
func (r *Radio) RxAEn() (Pipes, Status) {
	b, s := r.byteReg(2)
	return Pipes(b), s
}

// SetRxAEn sets the value of the EN_RXADDR (Enabled RX Addresses) register.
func (r *Radio) SetRxAEn(p Pipes) Status {
	return r.SetReg(2, byte(p))
}

// AW returns the value of the SETUP_AW (Setup of Address Widths) register
// increased by two, that is it returns the address length in bytes.
func (r *Radio) AW() (int, Status) {
	b, s := r.byteReg(3)
	return int(b) + 2, s
}

// SetAW sets the value of the SETUP_AW (Setup of Address Widths) register to
// (alen-2), that is it sets the address length to alen bytes.
func (r *Radio) SetAW(alen int) Status {
	if alen < 3 || alen > 5 {
		panic("alen<3 || alen>5")
	}
	return r.SetReg(3, byte(alen-2))
}

// Retr returns the value of the SETUP_RETR (Setup of Automatic Retransmission)
// register converted to number of retries and delay betwee retries.
func (r *Radio) Retr() (cnt, dlyus int, s Status) {
	b, s := r.byteReg(4)
	cnt = int(b & 0xf)
	dlyus = (int(b>>4) + 1) * 250
	return cnt, dlyus, s
}

// SetRetr sets the value of the SETUP_RETR (Setup of Automatic Retransmission)
// register using cnt as number of retries and dlyus as delay between retries.
func (r *Radio) SetRetr(cnt, dlyus int) Status {
	if uint(cnt) > 15 {
		panic("cnt<0 || cnt>15")
	}
	if dlyus < 250 || dlyus > 4000 {
		panic("dlyus<250 || dlyus>4000")
	}
	return r.SetReg(4, byte(((dlyus+125)/250-1)<<4|cnt))
}

// Ch returns the value of the RF_CH (RF Channel) register.
func (r *Radio) Ch() (int, Status) {
	b, s := r.byteReg(5)
	return int(b), s
}

// SetCh sets value of RF_CH (RF Channel) register.
func (r *Radio) SetCh(ch int) Status {
	if uint(ch) > 127 {
		panic("ch<0 || ch>127")
	}
	return r.SetReg(5, byte(ch))
}

type RF byte

const (
	LNAHC  RF = 1 << 0 // (nRF24L01.LNA_HCURR) Rx LNA gain 0: -1.5dB,-0.8mA.
	DRHigh RF = 1 << 3 // (RF_DR_HIGH) High speed data rate 0: 1Mbps, 1: 2Mbps.
	Lock   RF = 1 << 4 // (PLL_LOCK) Force PLL lock signal. Only used in test.
	DRLow  RF = 1 << 5 // (RF_DR_LOW) Set RF Data Rate to 250kbps.
	Wave   RF = 1 << 7 // (CONT_WAVE) Enable continuous carrier transmit.
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

func (rf RF) Format(fs fmt.State, _ rune) {
	fflags(fs, "Wave+ DRLow+ Lock+ DRHigh+ LNAHC+ Pwr:", 0xb9, byte(rf))
	strconv.WriteInt(fs, rf.Pwr(), 10, 0)
	io.WriteString(fs, " dBm")
}

// RF returns the value of the RF_SETUP register.
func (r *Radio) RF() (RF, Status) {
	b, s := r.byteReg(6)
	return RF(b), s
}

// RF sets the value of the RF_SETUP register.
func (r *Radio) SetRF(rf RF) Status {
	return r.SetReg(6, byte(rf))
}

// Clear clears the specified bits in the STATUS register.
func (r *Radio) Clear(s Status) Status {
	mask := RxDR | TxDS | MaxRT
	return r.SetReg(7, byte(s&mask))
}

// ObserveTx returns the values of PLOS and ARC counters from the OBSERVE_TX
// register.
func (r *Radio) ObserveTx() (plos, arc int, s Status) {
	b, s := r.byteReg(8)
	arc = int(b & 0xf)
	plos = int(b >> 4)
	return plos, arc, s
}

// RPD returns the value of the RPD (Received Power Detector) register
// (true if RP > -64dBm, false otherwise).
// In case of nRF24L01 it returns the value of the CD (Carrier Detect) register.
func (r *Radio) RPD() (bool, Status) {
	b, s := r.byteReg(9)
	return b&1 != 0, s
}

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
func (r *Radio) RxAddr(pn int, addr []byte) Status {
	checkPayNumAddr(pn, addr)
	return r.Reg(byte(0xa+pn), addr)
}

// RxAddrLSB returns least significant byte of address assigned to Rx pipe pn.
func (r *Radio) RxAddrLSB(pn int) (byte, Status) {
	checkPayNum(pn)
	return r.byteReg(byte(0xa + pn))
}

// SetRxAddr sets address assigned to Rx pipe pn to addr.
func (r *Radio) SetRxAddr(pn int, addr ...byte) Status {
	checkPayNumAddr(pn, addr)
	return r.SetReg(byte(0xa+pn), addr...)
}

// TxAddr reads value of TX_ADDR (Transmit address) into addr.
func (r *Radio) TxAddr(addr []byte) Status {
	checkAddr(addr)
	return r.Reg(0x10, addr)
}

// SetTxAddr sets value of TX_ADDR (Transmit address).
func (r *Radio) SetTxAddr(addr ...byte) Status {
	checkAddr(addr)
	return r.SetReg(0x10, addr...)
}

// RxPW returns Rx payload width set for pipe pn.
func (r *Radio) RxPW(pn int) (int, Status) {
	checkPayNum(pn)
	b, s := r.byteReg(byte(0x11 + pn))
	return int(b) & 0x3f, s
}

// SetRxPW sets Rx payload width for pipe pn.
func (r *Radio) SetRxPW(pn, pw int) Status {
	checkPayNum(pn)
	if uint(pw) > 32 {
		panic("pw<0 || pw>32")
	}
	return r.SetReg(byte(0x11+pn), byte(pw))
}

type FIFO byte

const (
	RxEmpty FIFO = 1 << 0 // 1: Rx FIFO empty, 0: Data in Rx FIFO.
	RxFull  FIFO = 1 << 1 // 1: Rx FIFO full, 0: Avail.locations in Rx FIFO.
	TxEmpty FIFO = 1 << 4 // 1: Tx FIFO empty, 0: Data in Tx FIFO.
	TxFull  FIFO = 1 << 5 // 1: Tx FIFO full, 0: Available locations in Tx FIFO.
	TxReuse FIFO = 1 << 6 // 1: Reuse last transmitted payload.
)

func (f FIFO) Format(fs fmt.State, _ rune) {
	fflags(fs, "TxReuse+ TxFull+ TxEmpty+ RxFull+ RxEmpty+", 0x73, byte(f))
}

// FIFO returns value of FIFO_STATUS register.
func (r *Radio) FIFO() (FIFO, Status) {
	b, s := r.byteReg(0x17)
	return FIFO(b), s
}

// DPL returns value of DYNPD (Enable dynamic payload length) register.
func (r *Radio) DPL() (Pipes, Status) {
	b, s := r.byteReg(0x1c)
	return Pipes(b), s
}

// SetDPL sets value of DYNPD (Enable dynamic payload length) register.
func (r *Radio) SetDPL(p Pipes) Status {
	return r.SetReg(0x1c, byte(p))
}

type Feature byte

const (
	DynAck Feature = 1 << 0 // 1: Enables the W_TX_PAYLOAD_NOACK command.
	AckPay Feature = 1 << 1 // 1: Enables payload with ACK
	DPL    Feature = 1 << 2 // 1: Enables dynamic payload length
)

func (f Feature) Format(fs fmt.State, _ rune) {
	fflags(fs, "DPL+ AckPay+ DynAck+", 7, byte(f))
}

// Feature returns value of FEATURE register.
func (r *Radio) Feature() (Feature, Status) {
	b, s := r.byteReg(0x1d)
	return Feature(b), s
}

// SetFeature sets value of FEATURE register.
func (r *Radio) SetFeature(f Feature) Status {
	return r.SetReg(0x1d, byte(f))
}
