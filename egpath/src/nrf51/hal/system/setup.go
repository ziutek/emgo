package system

import (
	"nrf51/hal/clock"
)

// Setup setups nRF51 to operate using specified HFCLK and LFCLK clock sources..
func Setup(hfsrc, lfsrc clock.SRC, lfena bool) {
	clkm := clock.Mgmt
	clkm.SetLFCLKSRC(lfsrc)
	if hfsrc == clock.Xtal {
		clkm.Task(clock.HFCLKSTART).Trigger()
	}
	if lfena {
		clkm.Task(clock.LFCLKSTART).Trigger()
	}
	for {
		src, run := clkm.HFCLKSTAT()
		if src == hfsrc && run {
			break
		}
	}
	for lfena {
		src, run := clkm.LFCLKSTAT()
		if src == lfsrc && run {
			break
		}
	}
}
