package system

import (
	"nrf5/hal/clock"
)

// Setup setups nRF51 to operate using specified HFCLK and LFCLK clock sources..
func Setup(hfsrc, lfsrc clock.Source, lfena bool) {
	clk := clock.CLOCK
	clk.StoreLFCLKSRC(lfsrc)
	if hfsrc == clock.XTAL {
		clk.Task(clock.HFCLKSTART).Trigger()
	}
	if lfena {
		clk.Task(clock.LFCLKSTART).Trigger()
	}
	for {
		src, run := clk.LoadHFCLKSTAT()
		if src == hfsrc && run {
			break
		}
	}
	for lfena {
		src, run := clk.LoadLFCLKSTAT()
		if src == lfsrc && run {
			break
		}
	}
}
