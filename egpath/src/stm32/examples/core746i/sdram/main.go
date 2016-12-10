package main

import (
	"delay"
	"unsafe"

	"stm32/hal/gpio"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"

	"stm32/hal/raw/fmc"
	"stm32/hal/raw/rcc"
)

func init() {
	system.Setup192(8)
	systick.Setup()

	gpio.D.EnableClock(false)
	d2 := gpio.Pin0
	d3 := gpio.Pin1
	d13 := gpio.Pin8
	d14 := gpio.Pin9
	d15 := gpio.Pin10
	d0 := gpio.Pin14
	d1 := gpio.Pin15
	fmcD := d2 | d3 | d13 | d14 | d15 | d0 | d1

	gpio.E.EnableClock(false)
	nbl0 := gpio.Pin0
	nbl1 := gpio.Pin1
	d4 := gpio.Pin7
	d5 := gpio.Pin8
	d6 := gpio.Pin9
	d7 := gpio.Pin10
	d8 := gpio.Pin11
	d9 := gpio.Pin12
	d10 := gpio.Pin13
	d11 := gpio.Pin14
	d12 := gpio.Pin15
	fmcE := nbl0 | nbl1 | d4 | d5 | d6 | d7 | d8 | d9 | d10 | d11 | d12

	gpio.F.EnableClock(false)
	a0 := gpio.Pin0
	a1 := gpio.Pin1
	a2 := gpio.Pin2
	a3 := gpio.Pin3
	a4 := gpio.Pin4
	a5 := gpio.Pin5
	sdnras := gpio.Pin11
	a6 := gpio.Pin12
	a7 := gpio.Pin13
	a8 := gpio.Pin14
	a9 := gpio.Pin15
	fmcF := a0 | a1 | a2 | a3 | a4 | a5 | sdnras | a6 | a7 | a8 | a9

	gpio.G.EnableClock(false)
	a10 := gpio.Pin0
	a11 := gpio.Pin1
	ba0 := gpio.Pin4
	ba1 := gpio.Pin5
	sdclk := gpio.Pin8
	sdncas := gpio.Pin15
	fmcG := a10 | a11 | ba0 | ba1 | sdclk | sdncas

	gpio.H.EnableClock(false)
	sdnwe := gpio.Pin5
	sdne1 := gpio.Pin6
	sdcke1 := gpio.Pin7
	fmcH := sdnwe | sdne1 | sdcke1

	// FMC

	gpio.D.Setup(fmcD, &gpio.Config{Mode: gpio.Alt, Speed: gpio.VeryHigh})
	gpio.D.SetAltFunc(fmcD, gpio.FMC)
	gpio.E.Setup(fmcE, &gpio.Config{Mode: gpio.Alt, Speed: gpio.VeryHigh})
	gpio.E.SetAltFunc(fmcE, gpio.FMC)
	gpio.F.Setup(fmcF, &gpio.Config{Mode: gpio.Alt, Speed: gpio.VeryHigh})
	gpio.F.SetAltFunc(fmcF, gpio.FMC)
	gpio.G.Setup(fmcG, &gpio.Config{Mode: gpio.Alt, Speed: gpio.VeryHigh})
	gpio.G.SetAltFunc(fmcG, gpio.FMC)
	gpio.H.Setup(fmcH, &gpio.Config{Mode: gpio.Alt, Speed: gpio.VeryHigh})
	gpio.H.SetAltFunc(fmcH, gpio.FMC)

	rcc.RCC.FMCEN().Set()

	sdcr1 := &fmc.FMC_Bank5_6.SDCR[0]
	sdcr2 := &fmc.FMC_Bank5_6.SDCR[1]

	// SDCLK period = 2 x HCLK period, burst read, one clock delay.
	sdcr1.Store(2<<fmc.SDCLKn | fmc.RBURST&0 | 1<<fmc.RPIPEn)

	// Col/row addr: 8/12 bit, 16bit data bus, 4 banks, CAS latency 3.
	sdcr2.Store(0<<fmc.NCn | 1<<fmc.NRn | 1<<fmc.SDMWIDn | fmc.NB | 3<<fmc.CASn)

	sdtr1 := &fmc.FMC_Bank5_6.SDTR[0]
	sdtr2 := &fmc.FMC_Bank5_6.SDTR[1]

	const (
		rowCycleDly     = 6
		rowPrechDly     = 2
		loadToActDly    = 2
		exitSelfRefrDly = 7
		selfRefrTime    = 4
		recoveryDly     = 2
		rowToColDly     = 2
	)
	sdtr1.Store((rowCycleDly-1)<<fmc.TRCn | (rowPrechDly-1)<<fmc.TRPn)
	sdtr2.Store(
		(loadToActDly-1)<<fmc.TMRDn | (exitSelfRefrDly-1)<<fmc.TXSRn |
			(selfRefrTime-1)<<fmc.TRASn | (recoveryDly-1)<<fmc.TWRn |
			(rowToColDly-1)<<fmc.TRCDn,
	)

	// SDRAM

	sdcmr := &fmc.FMC_Bank5_6.SDCMR
	sdsr := &fmc.FMC_Bank5_6.SDSR
	sdrtr := &fmc.FMC_Bank5_6.SDRTR

	const (
		cmdTarget    = fmc.CTB2
		clkConfEna   = 1 << fmc.MODEn
		prechargeAll = 2 << fmc.MODEn
		autoRefresh  = 3 << fmc.MODEn
		loadMode     = 4 << fmc.MODEn
	)

	sdcmr.Store(cmdTarget | clkConfEna)
	for sdsr.Bits(fmc.BUSY) != 0 {
	}
	delay.Millisec(1)
	sdcmr.Store(cmdTarget | prechargeAll)
	for sdsr.Bits(fmc.BUSY) != 0 {
	}
	sdcmr.Store(cmdTarget | autoRefresh | (8-1)<<fmc.NRFSn)
	for sdsr.Bits(fmc.BUSY) != 0 {
	}
	const (
		burstLength1     = 0
		burstSequential  = 0
		casLatency3      = 0x30
		modeStandard     = 0
		writeBurstSingle = 0x200

		ldm = burstLength1 | burstSequential | casLatency3 | modeStandard |
			writeBurstSingle
	)
	sdcmr.Store(cmdTarget | loadMode | ldm<<fmc.MRDn)
	for sdsr.Bits(fmc.BUSY) != 0 {
	}
	sdrtr.Store(0x603 << 1)
}

func main() {
	sdram := (*[8 * 1024 * 1024 / 4]uint32)(unsafe.Pointer(uintptr(0xD0000000)))

	for i := range sdram {
		sdram[i] = uint32(i)
	}
}

/*
#define SDRAM_MODEREG_BURST_LENGTH_1             ((uint16_t)0x0000)
#define SDRAM_MODEREG_BURST_LENGTH_2             ((uint16_t)0x0001)
#define SDRAM_MODEREG_BURST_LENGTH_4             ((uint16_t)0x0002)
#define SDRAM_MODEREG_BURST_LENGTH_8             ((uint16_t)0x0004)
#define SDRAM_MODEREG_BURST_TYPE_SEQUENTIAL      ((uint16_t)0x0000)
#define SDRAM_MODEREG_BURST_TYPE_INTERLEAVED     ((uint16_t)0x0008)
#define SDRAM_MODEREG_CAS_LATENCY_2              ((uint16_t)0x0020)
#define SDRAM_MODEREG_CAS_LATENCY_3              ((uint16_t)0x0030)
#define SDRAM_MODEREG_OPERATING_MODE_STANDARD    ((uint16_t)0x0000)
#define SDRAM_MODEREG_WRITEBURST_MODE_PROGRAMMED ((uint16_t)0x0000)
#define SDRAM_MODEREG_WRITEBURST_MODE_SINGLE     ((uint16_t)0x0200)
*/
