package ft80

import (
	"display/eve"
)

const (
	ACTIVE  eve.HostCmd = 0x00 // Switch mode to Active.
	STANDBY eve.HostCmd = 0x41 // Switch mode to Standby: PLL and Oscillator on.
	SLEEP   eve.HostCmd = 0x42 // Switch mode to Sleep: PLL and Oscillator off.
	PWRDOWN eve.HostCmd = 0x50 // Switch off LDO, Clock, PLL and Oscillator.

	CLKEXT eve.HostCmd = 0x44 // Select PLL external clock source.
	CLK48M eve.HostCmd = 0x62 // Switch PLL output to 48 MHz (default).
	CLK36M eve.HostCmd = 0x61 // Switch PLL output to 36 MHz.

	CORERST eve.HostCmd = 0x68 // Send reset pulse to FT80x core.
)
