// +build f746xx
package sdram

import (
	"bits"

	"stm32/hal/raw/fmc"
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

	// CAS sets CAS latency (mem. clock cycles, 1 to 3).
	CAS int8

	// TRCD sets row to column delay (mem. clock cycles, max. 16).
	TRCD int8

	// TWR sets WRI to PRE delay (mem. clock cycles, max. 16).
	TWR int8

	// TRAS sets ACT to PRE min. delay (mem. clock cycles, max. 16).
	TRAS int8

	// TXSR sets exit Self-refresh to Active delay (mem. clock cycles, max. 16).
	TXSR int8

	// TMRD sets Load Mode Register to Active delay (mem. clock cycles, max. 16).
	TMRD int8
}

type Conf struct {
	// ClkDiv sets SDRAM memory clock to AHBClk/ClkDiv. Alowed values: 2, 3.
	ClkDiv int8

	// ReadPipe sets read burst: -1 disable read burst; 0, 1, 2: enable read
	// burst and store data in Read FIFO during CAS latency + ReadPipe period.
	ReadPipe int8

	// TRP sets PRE to other command dleay (mem. clock cycles, max. 16).
	TRP int8

	// TRC sets REF to ACT, REF to REF delay (mem. clock cycles, max. 16).
	TRC int8

	// Banks contains configuration specific to two FMC SDRAM banks..
	Banks [2]BankConf
}

func SetConf(c *Conf) {
	var (
		sdcr [2]fmc.SDCR_Bits
		sdtr [2]fmc.SDTR_Bits
	)

	sdcr[0] = fmc.SDCR_Bits(c.ClkDiv&3) << fmc.SDCLKn
	if c.ReadPipe >= 0 {
		sdcr[0] |= fmc.RBURST | fmc.SDCR_Bits(c.ReadPipe)&3<<fmc.RPIPEn
	}
	sdtr[0] = fmc.SDTR_Bits(c.TRP-1) & 15 << fmc.TRPn
	sdtr[0] |= fmc.SDTR_Bits(c.TRC-1) & 15 << fmc.TRCn

	for i := 0; i < 2; i++ {
		b := &c.Banks[i]
		sdcr[i] |= fmc.SDCR_Bits(bits.One(b.WP)) << fmc.WPn
		sdcr[i] |= fmc.SDCR_Bits(b.BankNum/4) & 1 << fmc.NBn
		sdcr[i] |= fmc.SDCR_Bits(b.RowAddr-11) & 3 << fmc.NRn
		sdcr[i] |= fmc.SDCR_Bits(b.ColAddr-8) & 3 << fmc.NCn
		sdcr[i] |= fmc.SDCR_Bits(b.Bits/16) & 3 << fmc.MWIDn
		sdcr[i] |= fmc.SDCR_Bits(b.CAS) & 3 << fmc.CASn

		sdtr[i] |= fmc.SDTR_Bits(b.TRCD-1) & 15 << fmc.TRCDn
		sdtr[i] |= fmc.SDTR_Bits(b.TWR-1) & 15 << fmc.TWRn
		sdtr[i] |= fmc.SDTR_Bits(b.TRAS-1) & 15 << fmc.TRASn
		sdtr[i] |= fmc.SDTR_Bits(b.TXSR-1) & 15 << fmc.TXSRn
		sdtr[i] |= fmc.SDTR_Bits(b.TMRD-1) & 15 << fmc.TMRDn

		fmc.FMC_Bank5_6.SDCR[i].Store(sdcr[i])
		fmc.FMC_Bank5_6.SDTR[i].Store(sdtr[i])
	}
}
