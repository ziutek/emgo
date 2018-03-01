// +build f746xx
package sdram

import (
	"bits"

	"stm32/hal/raw/fmc"
	"stm32/hal/system"
)

type BankConf struct {
	// WP controls write protection.
	WP bool

	// BankNum sets the number of SDRAM internal banks.
	BankNum int8

	// RowAddr sets the row address width (bits).
	RowAddr int8

	// ColAddr sets the column address width (bits).
	ColAddr int8

	// Bits sets the data bus width / cell capacity (bits).
	Bits int8

	// CASL sets CAS latency (SDRAM clock cycles, 1 to 3).
	CASL int8

	// TRCDns sets row to column delay (nanoseconds).
	TRCDns int8

	// TWR+TWRns sets WRI to PRE delay (TWR>=TRAS-TRCD and TWR>=TRC-TRCD-TRP).
	TWR, TWRns int8

	// TRASns sets ACT to PRE min. delay (nanoseconds).
	TRASns int8

	// TXSRns sets exit Self-refresh to Active delay (nanoseconds).
	TXSRns int8

	// TMRD sets Load Mode Register to Active delay (SDRAM clock cycles).
	TMRD int8
}

type Conf struct {
	// ClkDiv sets SDRAM memory clock to AHBClk/ClkDiv. Alowed values: 2, 3.
	ClkDiv int8

	// ReadPipe sets read burst: -1: disable read burst; 0, 1, 2: enable read
	// burst and store data in Read FIFO during CAS latency + ReadPipe period.
	ReadPipe int8

	// TRPns sets PRE to other command dleay (nanoseconds).
	TRPns int8

	// TRC sets REF to ACT, REF to REF delay (nanoseconds).
	TRCns int8

	// TREFms sets refresh period per one row (milliseconds).
	TREFms int16

	// Banks contains configuration specific to two FMC SDRAM banks.
	Banks [2]BankConf
}

func nsclk(ns int8, kHz uint) fmc.SDTR {
	// Rounding up by adding 999424. It is slightly less than 999999 but fits in
	// the ADD instruction.
	return fmc.SDTR(uint(ns)*kHz+999424) / 1e6
}

func nsSDTR(ns int8, kHz uint, shift uint) fmc.SDTR {
	return (nsclk(ns, kHz) - 1) & 15 << shift
}

func Setup(c *Conf) {
	var (
		sdcr [2]fmc.SDCR
		sdtr [2]fmc.SDTR
	)

	kHz := system.AHB.Clock() / (1e3 * uint(c.ClkDiv)) // SDRAM clock

	sdcr[0] = fmc.SDCR(c.ClkDiv&3) << fmc.SDCLKn
	if c.ReadPipe >= 0 {
		sdcr[0] |= fmc.RBURST | fmc.SDCR(c.ReadPipe)&3<<fmc.RPIPEn
	}
	sdtr[0] = nsSDTR(c.TRPns, kHz, fmc.TRPn) | nsSDTR(c.TRCns, kHz, fmc.TRCn)

	for i := 0; i < 2; i++ {
		b := &c.Banks[i]
		sdcr[i] |= fmc.SDCR(bits.One(b.WP)) << fmc.WPn
		sdcr[i] |= fmc.SDCR(b.BankNum/4) & 1 << fmc.NBn
		sdcr[i] |= fmc.SDCR(b.RowAddr-11) & 3 << fmc.NRn
		sdcr[i] |= fmc.SDCR(b.ColAddr-8) & 3 << fmc.NCn
		sdcr[i] |= fmc.SDCR(b.Bits/16) & 3 << fmc.MWIDn
		sdcr[i] |= fmc.SDCR(b.CASL) & 3 << fmc.CASn

		sdtr[i] |= nsSDTR(b.TRCDns, kHz, fmc.TRCDn)
		sdtr[i] |= (fmc.SDTR(b.TWR) + nsclk(b.TWRns, kHz) - 1) & 15 << fmc.TWRn
		sdtr[i] |= nsSDTR(b.TRASns, kHz, fmc.TRASn)
		sdtr[i] |= nsSDTR(b.TXSRns, kHz, fmc.TXSRn)
		sdtr[i] |= fmc.SDTR(b.TMRD-1) & 15 << fmc.TMRDn

		fmc.FMC_Bank5_6.SDCR[i].Store(sdcr[i])
		fmc.FMC_Bank5_6.SDTR[i].Store(sdtr[i])
	}
	maxra := uint(c.Banks[0].RowAddr)
	if ra := uint(c.Banks[1].RowAddr); ra > maxra {
		maxra = ra
	}
	refclk := uint(c.TREFms)*kHz/1e3>>maxra - 20
	fmc.FMC_Bank5_6.SDRTR.Store(fmc.SDRTR(refclk-1) & 8191 << fmc.COUNTn)
}

type Banks byte

const (
	Bank0 = Banks(fmc.CTB1)
	Bank1 = Banks(fmc.CTB2)
)

type Mode byte

const (
	Normal      Mode = 0
	SelfRefresh Mode = 1
	PowerDown   Mode = 2
)

func SetMode(banks Banks, mode Mode) {
	var m uint32
	if mode != 0 {
		m = 4 + uint32(mode)
	}
	fmc.FMC_Bank5_6.SDCMR.U32.Store(m | uint32(banks))
}

func ClockConfEna(banks Banks) {
	fmc.FMC_Bank5_6.SDCMR.U32.Store(1 | uint32(banks))
}

func PrechargeAll(banks Banks) {
	fmc.FMC_Bank5_6.SDCMR.U32.Store(2 | uint32(banks))
}

func AutoRefresh(banks Banks, n int) {
	fmc.FMC_Bank5_6.SDCMR.U32.Store(3 | uint32(banks) | uint32(n-1)&15<<5)
}

type ModeReg uint16

const (
	Burst1   ModeReg = 0 << 0
	Burst2   ModeReg = 1 << 0
	Burst4   ModeReg = 2 << 0
	Burst8   ModeReg = 3 << 0
	BurstRow ModeReg = 7 << 0

	Interleaved ModeReg = 1 << 3

	CASL2 ModeReg = 2 << 4
	CASL3 ModeReg = 3 << 4

	SingleWrite ModeReg = 1 << 9
)

func LoadModeReg(banks Banks, mr ModeReg) {
	fmc.FMC_Bank5_6.SDCMR.U32.Store(4 | uint32(banks) | uint32(mr)<<9)
}

type ModeState struct {
	r fmc.SDSR
}

func (ms ModeState) Mode(bank int) Mode {
	return Mode(ms.r >> uint(fmc.MODES1n+2*bank) & 3)
}

func (ms ModeState) Busy() bool {
	return ms.r&fmc.BUSY != 0
}

func (ms ModeState) RefrErr() bool {
	return ms.r&fmc.RE != 0
}

func Status() ModeState {
	return ModeState{fmc.FMC_Bank5_6.SDSR.Load()}
}
