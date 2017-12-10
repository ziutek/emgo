package ft81

import (
	"display/eve"
)

const (
	ACTIVE  eve.HostCmd = 0x00 // Switch mode to Active.
	STANDBY eve.HostCmd = 0x41 // Switch mode to Standby: PLL and Oscillator on.
	SLEEP   eve.HostCmd = 0x42 // Switch mode to Sleep: PLL and Oscillator off.
	PWRDOWN eve.HostCmd = 0x43 // Switch off LDO, Clock, PLL and Oscillator.
	PD_ROMS eve.HostCmd = 0x49 // Power down individual ROMs.

	CLKEXT eve.HostCmd = 0x44 // Select PLL external clock source.
	CLKINT eve.HostCmd = 0x48 // Select PLL internal clock source.
	CLKSEL eve.HostCmd = 0x61 // Select PLL multiple.

	RST_PULSE eve.HostCmd = 0x68 // Send reset pulse to FT81x core.

	PINDRIVE     eve.HostCmd = 0x70 // Set pins drive strength.
	PIN_PD_STATE eve.HostCmd = 0x71 // Set pins state in PwrDown mode.
)
