package main

import (
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

	sdcr := fmc.FMC_Bank5_6.SDCR[1]
	_ = sdcr

}

func main() {
}
