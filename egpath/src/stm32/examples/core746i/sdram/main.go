package main

import (
	"delay"
	"fmt"
	"rtos"
	"sync/atomic"
	"unsafe"

	"stm32/hal/fmc"
	"stm32/hal/fmc/sdram"
	"stm32/hal/gpio"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"

	"stm32/hal/raw/rcc"
	"stm32/hal/raw/syscfg"
)

var (
	leds       *gpio.Port
	led1, led2 gpio.Pins
)

func init() {
	system.Setup192(8)
	systick.Setup(2e6)

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

	// Remap SDRAM to External RAM region (bank0:0x60000000, bank1:0x70000000).
	RCC.SYSCFGEN().Set()
	syscfg.SYSCFG.SWP_FMC().Store(1 << syscfg.SWP_FMCn)
	RCC.SYSCFGEN().Clear()

	fmc.EnableClock(true)

	// Parameters for IS42S16400J-7 SDRAM chip..
	sdram.Setup(&sdram.Conf{
		ClkDiv: 2, // SDRAMclk = AHBclk / 2
		TRPns:  15,
		TRCns:  70,
		TREFms: 64,
		Banks: [2]sdram.BankConf{
			1: {
				BankNum: 4,
				RowAddr: 12,
				ColAddr: 8,
				Bits:    16,
				CASL:    2, // SDRAMclk < 133 MHz
				TRCDns:  15,
				TWR:     2,
				TRASns:  42,
				TXSRns:  70,
				TMRD:    2,
			},
		},
	})

	// SDRAM

	sdram.ClockConfEna(sdram.Bank1)
	for sdram.Status().Busy() {
	}
	delay.Millisec(1)
	sdram.PrechargeAll(sdram.Bank1)
	for sdram.Status().Busy() {
	}
	sdram.AutoRefresh(sdram.Bank1, 8)
	for sdram.Status().Busy() {
	}
	sdram.LoadModeReg(sdram.Bank1, sdram.CASL2)
	for sdram.Status().Busy() {
	}
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

		t1 := rtos.Nanosec()
		for i := 0; i < len(sdram); i += 4 {
			atomic.StoreUint32(&sdram[i], uint32(i))
			atomic.StoreUint32(&sdram[i+1], uint32(i))
			atomic.StoreUint32(&sdram[i+2], uint32(i))
			atomic.StoreUint32(&sdram[i+3], uint32(i))
		}
		t2 := rtos.Nanosec()
		for i := 0; i < len(sdram); i += 4 {
			atomic.LoadUint32(&sdram[i])
			atomic.LoadUint32(&sdram[i+1])
			atomic.LoadUint32(&sdram[i+2])
			atomic.LoadUint32(&sdram[i+3])
		}
		t3 := rtos.Nanosec()
		for i := len(sdram) - 1; i > 0; i -= 4 {
			atomic.StoreUint32(&sdram[i], uint32(i))
			atomic.StoreUint32(&sdram[i-1], uint32(i))
			atomic.StoreUint32(&sdram[i-2], uint32(i))
			atomic.StoreUint32(&sdram[i-3], uint32(i))
		}
		t4 := rtos.Nanosec()
		for i := len(sdram) - 1; i > 0; i -= 4 {
			atomic.LoadUint32(&sdram[i])
			atomic.LoadUint32(&sdram[i-1])
			atomic.LoadUint32(&sdram[i-2])
			atomic.LoadUint32(&sdram[i-3])
		}
		t5 := rtos.Nanosec()
		for i := 0; i < len(sdram); i += 4 {
			v := atomic.LoadUint32(&sdram[i])
			atomic.StoreUint32(&sdram[i+1], v)
			v = atomic.LoadUint32(&sdram[i+2])
			atomic.StoreUint32(&sdram[i+3], v)
		}
		t6 := rtos.Nanosec()
		wt1 := (t2 - t1 + 0.5e6) / 1e6
		rt1 := (t3 - t2 + 0.5e6) / 1e6
		wt2 := (t4 - t3 + 0.5e6) / 1e6
		rt2 := (t5 - t4 + 0.5e6) / 1e6
		rw := (t6 - t5 + 0.5e6) / 1e6
		fmt.Printf(
			"wr+: %d ms, rd+: %d ms, wr-: %d ms, rd-: %d ms, rw: %d ms\n",
			wt1, rt1, wt2, rt2, rw,
		)
	}
}
