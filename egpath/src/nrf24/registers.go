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

func (r *Radio) byteReg(addr byte) (byte, STATUS) {
	var buf [1]byte
	s := r.R_REGISTER(addr, buf[:])
	return buf[0], s
}

type CONFIG byte

const (
	PRIM_RX     CONFIG = 1 << 0 // Primary: 1: Receiver (PRX) / 0: Transmitter (PPTX).
	PWR_UP      CONFIG = 1 << 1 // Power: 1: up / 0: down.
	CRCO        CONFIG = 1 << 2 // CRC encoding scheme: 0: one byte / 1: two bytes.
	EN_CRC      CONFIG = 1 << 3 // Enable CRC. Forced 1 if one of bits in EN_AA is 1.
	MASK_MAX_RT CONFIG = 1 << 4 // MAX_RT: 1: does not assert / 0: asserts IRQ.
	MASK_TX_DS  CONFIG = 1 << 5 // TX_DS: 1: does not assert, 0: asserts IRQ.
	MASK_RX_DR  CONFIG = 1 << 6 // RX_DR: 1: does not assert, 0: asserts IRQ.
)

func (c CONFIG) Format(fs fmt.State, _ rune) {
	fflags(
		fs, "MASK(RX_DR+ TX_DS+ MAX_RT+) EN_CRC+ CRCO+ PWR_UP+ PRIM_RX+",
		0x7f, byte(c),
	)
}

// CONFIG returns the value of the Configuration Register.
func (r *Radio) CONFIG() (CONFIG, STATUS) {
	b, s := r.byteReg(0)
	return CONFIG(b), s
}

// Set_CONFIG sets the value of the Configuration Register.
func (r *Radio) Set_CONFIG(c CONFIG) STATUS {
	return r.W_REGISTER(0, byte(c))
}

// Pipes is a bitfield that represents the nRF24L01+ Rx data pipes.
type Pipes byte

const (
	P0    Pipes = 1 << 0
	P1    Pipes = 1 << 1
	P2    Pipes = 1 << 2
	P3    Pipes = 1 << 3
	P4    Pipes = 1 << 4
	P5    Pipes = 1 << 5
	P_ALL       = P0 | P1 | P2 | P3 | P4 | P5
)

func (p Pipes) Format(fs fmt.State, _ rune) {
	fflags(fs, "P5+ P4+ P3+ P2+ P1+ P0+", 0x3f, byte(p))
}

// EN_AA returns the value of the Enable ‘Auto Acknowledgment’ Function
// register.
func (r *Radio) EN_AA() (Pipes, STATUS) {
	b, s := r.byteReg(1)
	return Pipes(b), s
}

// Set_EN_AA sets the value of the Enable ‘Auto Acknowledgment’ Function
// register.
func (r *Radio) Set_EN_AA(p Pipes) STATUS {
	return r.W_REGISTER(1, byte(p))
}

// EN_RXADDR returns the value of the Enabled RX Addresses register.
func (r *Radio) EN_RXADDR() (Pipes, STATUS) {
	b, s := r.byteReg(2)
	return Pipes(b), s
}

// Set_EN_RXADDR sets the value of the Enabled RX Addresses) register.
func (r *Radio) Set_EN_RXADDR(p Pipes) STATUS {
	return r.W_REGISTER(2, byte(p))
}

// SETUP_AW returns the value of the Setup of Address Widths register increased
// by two, that is it returns the address length in bytes.
func (r *Radio) SETUP_AW() (int, STATUS) {
	b, s := r.byteReg(3)
	return int(b) + 2, s
}

// Set_SETUP_AW sets the value of the Setup of Address Widths register to
// (aw-2), that is it sets the address length to aw bytes (allowed values:
// 3, 4, 5).
func (r *Radio) Set_SETUP_AW(alen int) STATUS {
	if alen < 3 || alen > 5 {
		panic("alen<3 || alen>5")
	}
	return r.W_REGISTER(3, byte(alen-2))
}

// SETUP_RETR returns the value of the Setup of Automatic Retransmission
// register converted to number of retries and delay (µs) between retries.
func (r *Radio) SETUP_RETR() (cnt, dlyus int, s STATUS) {
	b, s := r.byteReg(4)
	cnt = int(b & 0xf)
	dlyus = (int(b>>4) + 1) * 250
	return cnt, dlyus, s
}

// Set_SETUP_RETR sets the value of the Setup of Automatic Retransmission
// register using cnt as number of retries and dlyus as delay between retries.
func (r *Radio) Set_SETUP_RETR(cnt, dlyus int) STATUS {
	if uint(cnt) > 15 {
		panic("cnt<0 || cnt>15")
	}
	if dlyus < 250 || dlyus > 4000 {
		panic("dlyus<250 || dlyus>4000")
	}
	return r.W_REGISTER(4, byte(((dlyus+125)/250-1)<<4|cnt))
}

// RF_CH returns the value of the RF Channel register.
func (r *Radio) RF_CH() (int, STATUS) {
	b, s := r.byteReg(5)
	return int(b), s
}

// Set_RF_CH sets value of RF Channel register.
func (r *Radio) Set_RF_CH(ch int) STATUS {
	if uint(ch) > 127 {
		panic("ch<0 || ch>127")
	}
	return r.W_REGISTER(5, byte(ch))
}

type RF_SETUP byte

const (
	LNA_HCURR  RF_SETUP = 1 << 0 // LNA gain 0: -1.5 dB, -0.8 mA (nRF24L01 specific).
	RF_DR_HIGH RF_SETUP = 1 << 3 // High speed data rate 0: 1Mbps, 1: 2Mbps.
	PLL_LOCK   RF_SETUP = 1 << 4 // Force PLL lock signal. Only used in test.
	RF_DR_LOW  RF_SETUP = 1 << 5 // Set RF Data Rate to 250kbps.
	CONT_WAVE  RF_SETUP = 1 << 7 // Enable continuous carrier transmit.
)

// RF_PWR returns RF output power in Tx mode [dBm].
func (rf RF_SETUP) RF_PWR() int {
	return 3*int(rf&6) - 18
}

func RF_PWR(dbm int) RF_SETUP {
	switch {
	case dbm < -18:
		dbm = -18
	case dbm > 0:
		dbm = 0
	}
	return RF_SETUP((18+dbm)/3) & 6
}

func (rf RF_SETUP) Format(fs fmt.State, _ rune) {
	fflags(
		fs, "CONT_WAVE+ RF_DR_LOW+ PLL_LOCK+ RF_DR_HIGH+ LNA_HCURR+ RF_PWR:",
		0xb9, byte(rf),
	)
	strconv.WriteInt(fs, rf.RF_PWR(), 10, 0)
	io.WriteString(fs, " dBm")
}

// RF_SETUP returns the value of the RF Setup register.
func (r *Radio) RF_SETUP() (RF_SETUP, STATUS) {
	b, s := r.byteReg(6)
	return RF_SETUP(b), s
}

// Set_RF_SETUP sets the value of the RF Setup register.
func (r *Radio) Set_RF_SETUP(rf RF_SETUP) STATUS {
	return r.W_REGISTER(6, byte(rf))
}

type STATUS byte

const (
	FULL_TX STATUS = 1 << 0 // Tx FIFO full flag.
	MAX_RT  STATUS = 1 << 4 // Maximum number of Tx retransmits interrupt.
	TX_DS   STATUS = 1 << 5 // Data Sent Tx FIFO interrupt.
	RX_DR   STATUS = 1 << 6 // Data Ready Rx FIFO interrupt.
)

// RX_P_NO returns the data pipe number for the payload available for reading
// from Rx FIFO or -1 if Tx FIFO is empty.
func (s STATUS) RX_P_NO() int {
	n := int(s) & 0x0e
	if n == 0x0e {
		return -1
	}
	return n >> 1
}

func (s STATUS) Format(fs fmt.State, _ rune) {
	fflags(fs, "RX_DR+ TX_DS+ MAX_RT+ FULL_TX+ RX_P_NO:", 0x71, byte(s))
	strconv.WriteInt(fs, s.RX_P_NO(), 10, 0)
}

// ClearIRQ allow to clear the interrupt bits of the Status register.
func (r *Radio) ClearIRQ(s STATUS) STATUS {
	mask := RX_DR | TX_DS | MAX_RT
	return r.W_REGISTER(7, byte(s&mask))
}

// OBSERVE_TX returns the values of PLOS_CNT and ARC_CNT counters from the
// Transmit observe register.
func (r *Radio) OBSERVE_TX() (plos, arc int, s STATUS) {
	b, s := r.byteReg(8)
	arc = int(b & 0xf)
	plos = int(b >> 4)
	return plos, arc, s
}

// RPD returns the value of the Received Power Detector register: true if
// RP > -64dBm, false otherwise. In case of nRF24L01 it returns the value of
// the CD (Carrier Detect) register.
func (r *Radio) RPD() (bool, STATUS) {
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

// Read_RX_ADDR reads the receive address of the data pipe pn.
func (r *Radio) Read_RX_ADDR(pn int, addr []byte) STATUS {
	checkPayNumAddr(pn, addr)
	return r.R_REGISTER(byte(0xa+pn), addr)
}

// Write_RX_ADDR sets the receive address of the data pipe pn.
func (r *Radio) Write_RX_ADDR(pn int, addr ...byte) STATUS {
	checkPayNumAddr(pn, addr)
	return r.W_REGISTER(byte(0xa+pn), addr...)
}

// Read_TX_ADDR reads value of Transmit address register into addr.
func (r *Radio) Read_TX_ADDR(addr []byte) STATUS {
	checkAddr(addr)
	return r.R_REGISTER(0x10, addr)
}

// Write_TX_ADDR sets value of Transmit address.
func (r *Radio) Set_TX_ADDR(addr ...byte) STATUS {
	checkAddr(addr)
	return r.W_REGISTER(0x10, addr...)
}

// RX_PW returns the Rx payload width for pipe pn.
func (r *Radio) RX_PW(pn int) (int, STATUS) {
	checkPayNum(pn)
	b, s := r.byteReg(byte(0x11 + pn))
	return int(b) & 0x3f, s
}

// Set_RX_PW sets the Rx payload width for pipe pn.
func (r *Radio) Set_RX_PW(pn, pw int) STATUS {
	checkPayNum(pn)
	if uint(pw) > 32 {
		panic("pw<0 || pw>32")
	}
	return r.W_REGISTER(byte(0x11+pn), byte(pw))
}

type FIFO_STATUS byte

const (
	RX_EMPTY FIFO_STATUS = 1 << 0 // 1: Rx FIFO empty, 0: Data in Rx FIFO.
	RX_FULL  FIFO_STATUS = 1 << 1 // 1: Rx FIFO full, 0: Available locations in Rx FIFO.
	TX_EMPTY FIFO_STATUS = 1 << 4 // 1: Tx FIFO empty, 0: Data in Tx FIFO.
	TX_FULL  FIFO_STATUS = 1 << 5 // 1: Tx FIFO full, 0: Available locations in Tx FIFO.
	TX_REUSE FIFO_STATUS = 1 << 6 // Set by TX_REUSE, cleared by W_TX_PAYLOAD, FLUSH_TX.
)

func (f FIFO_STATUS) Format(fs fmt.State, _ rune) {
	fflags(fs, "TX_REUSE+ TX_FULL+ TX_EMPTY+ RX_FULL+ RX_EMPTY+", 0x73, byte(f))
}

// FIFO_STATUS returns value of FIFO Status Register.
func (r *Radio) FIFO_STATUS() (FIFO_STATUS, STATUS) {
	b, s := r.byteReg(0x17)
	return FIFO_STATUS(b), s
}

// DYNPD returns the value of Enable dynamic payload length register.
func (r *Radio) DYNPD() (Pipes, STATUS) {
	b, s := r.byteReg(0x1c)
	return Pipes(b), s
}

// Set_DYNPD sets the value of Enable dynamic payload length register.
func (r *Radio) Set_DYNPD(p Pipes) STATUS {
	return r.W_REGISTER(0x1c, byte(p))
}

type FEATURE byte

const (
	EN_DYN_ACK FEATURE = 1 << 0 // 1: Enables the W_TX_PAYLOAD_NOACK command.
	EN_ACK_PAY FEATURE = 1 << 1 // 1: Enables payload with ACK
	EN_DPL     FEATURE = 1 << 2 // 1: Enables dynamic payload length
)

func (f FEATURE) Format(fs fmt.State, _ rune) {
	fflags(fs, "EN_DPL+ EN_ACK_PAY+ EN_DYN_ACK+", 7, byte(f))
}

// FEATURE returns value of Feature Register.
func (r *Radio) FEATURE() (FEATURE, STATUS) {
	b, s := r.byteReg(0x1d)
	return FEATURE(b), s
}

// Set_FEATURE sets value of FEATURE register.
func (r *Radio) Set_FEATURE(f FEATURE) STATUS {
	return r.W_REGISTER(0x1d, byte(f))
}
