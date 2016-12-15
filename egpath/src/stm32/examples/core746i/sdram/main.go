package main

import (
	"delay"
	"fmt"
	"rtos"
	"sync/atomic"
	"unsafe"

	"stm32/hal/fmc/sdram"
	"stm32/hal/gpio"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"

	"stm32/hal/raw/fmc"
	"stm32/hal/raw/rcc"
	"stm32/hal/raw/syscfg"
)

var (
	leds       *gpio.Port
	led1, led2 gpio.Pins
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
	leds, led1, led2 = gpio.H, gpio.Pin3, gpio.Pin4
	sdnwe := gpio.Pin5
	sdne1 := gpio.Pin6
	sdcke1 := gpio.Pin7
	fmcH := sdnwe | sdne1 | sdcke1

	RCC := rcc.RCC

	// LEDs

	leds.Setup(led1|led2, &gpio.Config{Mode: gpio.Out, Speed: gpio.Low})

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

	// Remap SDRAM to External RAM region (bank1:0x60000000, bank2:0x70000000).
	RCC.SYSCFGEN().Set()
	syscfg.SYSCFG.SWP_FMC().Store(1 << syscfg.SWP_FMCn)
	RCC.SYSCFGEN().Clear()

	RCC.FMCEN().Set()

	/*
		sdcr1 := &fmc.FMC_Bank5_6.SDCR[0]
		sdcr2 := &fmc.FMC_Bank5_6.SDCR[1]

		// SDCLK period = 2 x HCLK period, burst read, zero clock delay.
		sdcr1.Store(2<<fmc.SDCLKn | fmc.RBURST | 0<<fmc.RPIPEn)

		// Col/row addr: 8/12 bit, 16bit data bus, 4 banks, CAS latency 2.
		sdcr2.Store(0<<fmc.NCn | 1<<fmc.NRn | 1<<fmc.SDMWIDn | fmc.NB | 2<<fmc.CASn)

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
	*/

	sdram.SetConf(&sdram.Conf{
		ClkDiv: 2, // 192 MHz / 2 = 96 MHz (period ≥ 10 ns)
		TRP:    2, // 20 ns > 15 ns ISSI
		TRC:    6, // 60 ns > 55 ns ISSI
		Banks: [2]sdram.BankConf{
			1: {
				BankNum: 4,
				RowAddr: 12,
				ColAddr: 8,
				Bits:    16,
				CAS:     2, // 96 MHz < 133 MHz ISSI
				TRCD:    2, // 20 ns > 15 ns ISSI
				TWR:     2, // 2CLK = 2CLK ISSI (TWR ≥ TRAS-TRCD and TWR ≥ TRC-TRCD-TRP)
				TRAS:    4, // 40 ns = 40 ns ISSI
				TXSR:    6, // 60 ns = 60 ns ISSI
				TMRD:    2, // 2CLK = 2CLK ISSI
			},
		},
	})

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
		casLatency2      = 0x20
		modeStandard     = 0
		writeBurstSingle = 0x200

		ldm = burstLength1 | burstSequential | casLatency2 | modeStandard |
			writeBurstSingle
	)
	sdcmr.Store(cmdTarget | loadMode | ldm<<fmc.MRDn)
	for sdsr.Bits(fmc.BUSY) != 0 {
	}
	sdrtr.Store(0x603 << 1)
}

func main() {
	sdram := (*[8 << 20 / 4]uint32)(unsafe.Pointer(uintptr(0x70000000)))

	for n := 0; ; n++ {
		for i := 0; i < len(sdram); i += 2 {
			sdram[i] = uint32(n ^ i)
			sdram[i+1] = uint32(n ^ i + 1)
		}
		leds.SetPins(led1)
		for i := 0; i < len(sdram); i += 2 {
			a := sdram[i]
			b := sdram[i+1]
			if a != uint32(n^i) || b != uint32(n^i+1) {
				leds.SetPins(led2)
				return
			}
		}
		leds.ClearPins(led1)

		t := rtos.Nanosec()
		for i := 0; i < len(sdram); i += 4 {
			atomic.StoreUint32(&sdram[i], uint32(i))
			atomic.StoreUint32(&sdram[i+1], uint32(i))
			atomic.StoreUint32(&sdram[i+2], uint32(i))
			atomic.StoreUint32(&sdram[i+3], uint32(i))
		}
		wt1 := (rtos.Nanosec() - t + 0.5e6) / 1e6
		t = rtos.Nanosec()
		for i := 0; i < len(sdram); i += 4 {
			atomic.LoadUint32(&sdram[i])
			atomic.LoadUint32(&sdram[i+1])
			atomic.LoadUint32(&sdram[i+2])
			atomic.LoadUint32(&sdram[i+3])
		}
		rt1 := (rtos.Nanosec() - t + 0.5e6) / 1e6
		t = rtos.Nanosec()
		for i := len(sdram) - 4; i > 0; i -= 4 {
			atomic.StoreUint32(&sdram[i-3], uint32(i))
			atomic.StoreUint32(&sdram[i-2], uint32(i))
			atomic.StoreUint32(&sdram[i-1], uint32(i))
			atomic.StoreUint32(&sdram[i], uint32(i))
		}
		wt2 := (rtos.Nanosec() - t + 0.5e6) / 1e6
		t = rtos.Nanosec()
		for i := len(sdram) - 4; i > 0; i -= 4 {
			atomic.LoadUint32(&sdram[i-3])
			atomic.LoadUint32(&sdram[i-2])
			atomic.LoadUint32(&sdram[i-1])
			atomic.LoadUint32(&sdram[i])
		}
		rt2 := (rtos.Nanosec() - t + 0.5e6) / 1e6
		t = rtos.Nanosec()
		for i := 0; i < len(sdram); i += 4 {
			v := atomic.LoadUint32(&sdram[i])
			atomic.StoreUint32(&sdram[i+1], v)
			v = atomic.LoadUint32(&sdram[i+2])
			atomic.StoreUint32(&sdram[i+3], v)
		}
		rw := (rtos.Nanosec() - t + 0.5e6) / 1e6
		fmt.Printf(
			"wr+: %d ms, rd+: %d ms, wr-: %d ms, rd-: %d ms, rw: %d ms\n",
			wt1, rt1, wt2, rt2, rw,
		)
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
