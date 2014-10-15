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
	SysClk = 32e6

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
		mul = clock.PLLMul3
	default:
		panic("wrong frequency of external resonator")
	}

	// Set MSI as system clock source
	/*clock.EnableMSI()
	clock.SetSysClock(clock.MSI)
	for clock.SysClock() != clock.MSI {
	}*/

	flash.SetAcc64(true)
	for !flash.Acc64() {
	}
	flash.SetLatency(1) // need for 2.0-3.6V and 16-32MHz
	flash.SetPrefetch(true)

	// Be sure that flash latency is set before incrase frequency.
	for flash.Latency() != 1 {
	}

	// Reset clock subsystem
	clock.ResetCR()
	clock.ResetICSCR()
	clock.ResetCFGR()
	clock.ResetCIR()

	if osc != 0 {
		clock.EnableHSE()
	}

	// Configure maximum clocks frequency (32 MHz) for AHB, APB1, APB2 bus.
	clock.SetPrescalerAHB(clock.AHBDiv1)
	clock.SetPrescalerAPB1(clock.APBDiv1)
	clock.SetPrescalerAPB2(clock.APBDiv1)
	APB1Clk = SysClk
	APB2Clk = SysClk

	if osc == 0 {
		for !clock.HSIReady() {
		}
		clock.SetPLLClock(clock.SrcHSI)
	} else {
		for !clock.HSEReady() {
		}
		clock.SetPLLClock(clock.SrcHSE)
	}
	clock.SetPLLMul(mul)
	clock.SetPLLDiv(2)
	clock.EnableMainPLL()
	for !clock.MainPLLReady() {
	}

	// Set PLL as system clock source
	clock.SetSysClock(clock.PLL)
	for clock.SysClock() != clock.PLL {
	}

	clock.DisableMSI()

	sysClkChanged()
}
