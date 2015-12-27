// Peripheral: IWDG_Periph  Independent WATCHDOG.
// Instances:
//  IWDG  mmap.IWDG_BASE
// Registers:
//  0x00 32  KR  Key register.
//  0x04 32  PR  Prescaler register.
//  0x08 32  RLR Reload register.
//  0x0C 32  SR  Status register.
// Import:
//  stm32/o/l1xx_md/mmap
package iwdg

const (
	KEY KR_Bits = 0xFFFF << 0 //+ Key value (write only, read 0000h).
)

const (
	PR   PR_Bits = 0x07 << 0 //+ PR[2:0] (Prescaler divider).
	PR_0 PR_Bits = 0x01 << 0 //  Bit 0.
	PR_1 PR_Bits = 0x02 << 0 //  Bit 1.
	PR_2 PR_Bits = 0x04 << 0 //  Bit 2.
)

const (
	RL RLR_Bits = 0xFFF << 0 //+ Watchdog counter reload value.
)

const (
	PVU SR_Bits = 0x01 << 0 //+ Watchdog prescaler value update.
	RVU SR_Bits = 0x01 << 1 //+ Watchdog counter reload value update.
)
