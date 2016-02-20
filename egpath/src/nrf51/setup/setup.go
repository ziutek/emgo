package setup

import (
	"nrf51/clock"
)

// Clocks setups HFCLK and LFCLK.
func Clocks(hfsrc, lfsrc clock.SRC, lfena bool) {
	clkm := clock.Mgmt
	clkm.SetLFCLKSRC(lfsrc)
	if hfsrc == clock.Xtal {
		clkm.TASK(clock.HFCLKSTART).Trigger()
	}
	if lfena {
		clkm.TASK(clock.LFCLKSTART).Trigger()
	}
wait:
	src, run := clkm.HFCLKSTAT()
	if src != hfsrc || !run {
		goto wait
	}
	if lfena {
		src, run = clkm.LFCLKSTAT()
		if src != lfsrc || !run {
			goto wait
		}
	}
	sysclkChanged()
}
