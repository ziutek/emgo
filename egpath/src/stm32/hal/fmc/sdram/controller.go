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

	// CAS sets CAS latency (SDRAM clock cycles, 1 to 3).
	CAS int8

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

func nsclk(ns int8, kHz uint) fmc.SDTR_Bits {
	// Rounding up by adding 999424. It is slightly less than 999999 but fits in
	// the ADD instruction.
	return fmc.SDTR_Bits(uint(ns)*kHz+999424) / 1e6
}

//emgo:noinline
func nsSDTR(ns int8, kHz uint) fmc.SDTR_Bits {
	return (nsclk(ns, kHz) - 1) & 15
}

func SetController(c *Conf) {
	var (
		sdcr [2]fmc.SDCR_Bits
		sdtr [2]fmc.SDTR_Bits
	)

	kHz := system.AHB.Clock() / (1e3 * uint(c.ClkDiv)) // SDRAM clock

	sdcr[0] = fmc.SDCR_Bits(c.ClkDiv&3) << fmc.SDCLKn
	if c.ReadPipe >= 0 {
		sdcr[0] |= fmc.RBURST | fmc.SDCR_Bits(c.ReadPipe)&3<<fmc.RPIPEn
	}
	sdtr[0] = nsSDTR(c.TRPns, kHz)<<fmc.TRPn | nsSDTR(c.TRCns, kHz)<<fmc.TRCn

	for i := 0; i < 2; i++ {
		b := &c.Banks[i]
		sdcr[i] |= fmc.SDCR_Bits(bits.One(b.WP)) << fmc.WPn
		sdcr[i] |= fmc.SDCR_Bits(b.BankNum/4) & 1 << fmc.NBn
		sdcr[i] |= fmc.SDCR_Bits(b.RowAddr-11) & 3 << fmc.NRn
		sdcr[i] |= fmc.SDCR_Bits(b.ColAddr-8) & 3 << fmc.NCn
		sdcr[i] |= fmc.SDCR_Bits(b.Bits/16) & 3 << fmc.MWIDn
		sdcr[i] |= fmc.SDCR_Bits(b.CAS) & 3 << fmc.CASn

		sdtr[i] |= nsSDTR(b.TRCDns, kHz) << fmc.TRCDn
		sdtr[i] |= (fmc.SDTR_Bits(b.TWR) + nsclk(b.TWRns, kHz) - 1) & 15 << fmc.TWRn
		sdtr[i] |= nsSDTR(b.TRASns, kHz) << fmc.TRASn
		sdtr[i] |= nsSDTR(b.TXSRns, kHz) << fmc.TXSRn
		sdtr[i] |= fmc.SDTR_Bits(b.TMRD-1) & 15 << fmc.TMRDn

		fmc.FMC_Bank5_6.SDCR[i].Store(sdcr[i])
		fmc.FMC_Bank5_6.SDTR[i].Store(sdtr[i])
	}
	maxra := uint(c.Banks[0].RowAddr)
	if ra := uint(c.Banks[1].RowAddr); ra > maxra {
		maxra = ra
	}
	refclk := uint(c.TREFms)*kHz/1e3>>maxra - 20
	fmc.FMC_Bank5_6.SDRTR.Store(fmc.SDRTR_Bits(refclk-1) & 8191 << fmc.COUNTn)
}
