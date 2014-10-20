package setup

import (
	"stm32/l1/clock"
	"stm32/l1/flash"
)

var (
	SysClk  uint // System clock [Hz]
	APB1Clk uint // APB1 clock [Hz]
	APB2Clk uint // APB2 clock [Hz]
)

// Performance setups MCU for best performance (32MHz, Flash prefetch and
// 64-bit access on).
// It accepts one argument which is a freqency of external resonator in MHz.
// Allowed values: 2, 3, 4, 6, 8, 12, 16, 24. Use 0 to select internal HSI
// oscilator as system clock source.
func Performance(osc int) {
	// Reset clock subsystem
	clock.ResetCR()
	clock.ResetICSCR()
	clock.ResetCFGR()
	clock.ResetCIR()

	// Set HSI as temporary system clock source.
	clock.EnableHSI()
	for !clock.HSIReady() {
	}
	clock.SetSysClock(clock.HSI)
	for clock.SysClock() != clock.HSI {
	}
	clock.DisableMSI()

	// Set mul to obtain PLLVCO=96MHz (need by USB)
	var mul clock.PLLMul
	switch osc {
	case 2:
		mul = clock.PLLMul48
	case 3:
		mul = clock.PLLMul32
	case 4:
		mul = clock.PLLMul24
	case 6:
		mul = clock.PLLMul16
	case 8:
		mul = clock.PLLMul12
	case 12:
		mul = clock.PLLMul8
	case 16, 0:
		mul = clock.PLLMul6
	case 24:
		mul = clock.PLLMul4
	default:
		panic("wrong frequency of external resonator")
	}

	// HSE needs milliseconds to stabilize, so enable it now.
	if osc != 0 {
		clock.EnableHSE()
	}

	flash.SetAcc64(true)
	for !flash.Acc64() {
	}
	flash.SetLatency(1) // need for 2.0-3.6V and 16-32MHz
	flash.SetPrefetch(true)

	// Be sure that flash latency is set before incrase frequency.
	for flash.Latency() != 1 {
	}

	// Configure maximum clocks frequency (32 MHz) for AHB, APB1, APB2 bus.
	clock.SetPrescalerAHB(clock.AHBDiv1)
	clock.SetPrescalerAPB1(clock.APBDiv1)
	clock.SetPrescalerAPB2(clock.APBDiv1)
	SysClk = 32e6
	APB1Clk = SysClk
	APB2Clk = SysClk

	if osc == 0 {
		clock.SetPLLClock(clock.SrcHSI)
	} else {
		for !clock.HSEReady() {
		}
		clock.SetPLLClock(clock.SrcHSE)
	}
	clock.SetPLLMul(mul)
	clock.SetPLLDiv(3)
	clock.EnableMainPLL()
	for !clock.MainPLLReady() {
	}

	// Set PLL as system clock source
	clock.SetSysClock(clock.PLL)
	for clock.SysClock() != clock.PLL {
	}

	sysClkChanged()
}
