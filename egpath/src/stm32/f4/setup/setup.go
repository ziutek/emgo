package setup

import (
	"stm32/f4/clock"
	"stm32/f4/flash"
)

// Performance setups MCU for best performance (168MHz, Flash prefetch and
// I/D cache on).
// It accepts one argument which is freqency of external resonator in MHz.
// Allowed values are multiples of 2 MHz, from 4 MHz to 26 MHz. Use 0 to
// select internal HSI oscilator as system clock source.
func Performance(osc int) {
	if osc < 0 || osc > 26 || osc == 2 || osc&1 != 0 {
		panic("wrong frequency of external resonator")
	}

	flash.SetLatency(5) // need for 2.7-3.6V and 150-168MHz
	flash.SetPrefetch(true)
	flash.SetICache(true)
	flash.SetDCache(true)

	// Be sure that flash latency is set before incrase frequency.
	for flash.Latency() != 5 {
	}

	// Reset clock subsystem
	clock.ResetCR()
	clock.ResetPLLCFGR()
	clock.ResetCFGR()
	clock.ResetCIR()

	if osc != 0 {
		clock.EnableHSE()
	}

	// Configure clocks for AHB, APB1, APB2 bus.
	clock.SetPrescalerAHB(clock.AHBDiv1)
	clock.SetPrescalerAPB1(clock.APBDiv4) // APB1Freq must be <= 42 MHz
	clock.SetPrescalerAPB2(clock.APBDiv2) // APB2Freq must be <= 84 MHz

	if osc != 0 {
		// Be sure that HSE is ready.
		for !clock.HSEReady() {
		}
		clock.SetPLLClock(clock.SrcHSE)
		clock.SetPLLInputDiv(osc / 2)
	} else {
		clock.SetPLLClock(clock.SrcHSI)
		clock.SetPLLInputDiv(16 / 2)
	}

	clock.SetMainPLLMul(168)     // 336 MHz
	clock.SetMainPLLSysDiv(2)    // 168 MHz
	clock.SetMainPLLPeriphDiv(7) // 48 MHz
	clock.EnableMainPLL()
	for !clock.MainPLLReady() {
	}

	// Set PLL as system clock source
	clock.SetSysClock(clock.PLL)
	for clock.SysClock() != clock.PLL {
	}

	if osc != 0 {
		clock.DisableHSI()
	}
}