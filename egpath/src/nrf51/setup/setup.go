package setup

import (
	"nrf51/clock"
)

// Clocks setups HFCLK and LFCLK.
func Clocks(hfsrc, lfsrc clock.Src, lfena bool) {
	clkm := clock.Mgmt
	clkm.SetLFCLKSrc(lfsrc)
	if hfsrc == clock.Xtal {
		clkm.Task(clock.HFCLKSTART).Trig()
	}
	if lfena {
		clkm.Task(clock.LFCLKSTART).Trig()
	}
wait:
	src, run := clkm.HFCLKStat()
	if src != hfsrc || !run {
		goto wait
	}
	if lfena {
		src, run = clkm.LFCLKStat()
		if src != lfsrc || !run {
			goto wait
		}
	}
	sysclkChanged()
}
