package main

import (
	"stm32/hal/gpio"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"

	"stm32/hal/raw/rcc"
	"stm32/hal/raw/sdio"
)

func init() {
	system.Setup96(8) // Setups USB/SDIO/RNG clock to 48 MHz
	systick.Setup(2e6)

	// GPIO

	gpio.A.EnableClock(true)
	cmd := gpio.A.Pin(6)
	//d1 := gpio.A.Pin(8)
	//d2 := gpio.A.Pin(9)

	gpio.B.EnableClock(true)
	//d3 := gpio.B.Pin(5)
	d0 := gpio.B.Pin(7)
	clk := gpio.B.Pin(15)

	cfg := gpio.Config{Mode: gpio.Alt, Speed: gpio.High}
	clk.Setup(&cfg)
	clk.SetAltFunc(gpio.SDIO)
	cfg.Pull = gpio.PullUp
	for _, pin := range []gpio.Pin{cmd, d0 /*d1, d2, d3*/} {
		pin.Setup(&cfg)
		pin.SetAltFunc(gpio.SDIO)
	}

	rcc.RCC.SDIOEN().Set()
	sd := sdio.SDIO
	sd.CLKCR.Store(sdio.CLKEN | (48e6/400e3-2)<<sdio.CLKDIVn) // 400 kHz
	sd.POWER.Store(3)                                         // Power on.
}

func sdioCMD(cmd sdio.CMD, arg sdio.ARG) {
	sd := sdio.SDIO
	for i := 0; i < 10; i++ {
		sd.ICR.Store(0xFFFFFFFF)
		sd.ARG.Store(arg)
		sd.CMD.Store(sdio.CPSMEN | cmd)
		
		/*
		for sd.CMDACT().Load() != 0 {
			// Wait for transfer end.
		}
		*/
	}
}

func main() {

}
