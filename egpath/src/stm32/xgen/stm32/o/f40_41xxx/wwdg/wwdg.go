// Peripheral: WWDG_Periph  Window WATCHDOG.
// Instances:
//  WWDG  mmap.WWDG_BASE
// Registers:
//  0x00 32  CR  Control register.
//  0x04 32  CFR Configuration register.
//  0x08 32  SR  Status register.
// Import:
//  stm32/o/f40_41xxx/mmap
package wwdg

const (
	T    CR_Bits = 0x7F << 0 //+ T[6:0] bits (7-Bit counter (MSB to LSB)).
	T0   CR_Bits = 0x01 << 0 //  Bit 0.
	T1   CR_Bits = 0x02 << 0 //  Bit 1.
	T2   CR_Bits = 0x04 << 0 //  Bit 2.
	T3   CR_Bits = 0x08 << 0 //  Bit 3.
	T4   CR_Bits = 0x10 << 0 //  Bit 4.
	T5   CR_Bits = 0x20 << 0 //  Bit 5.
	T6   CR_Bits = 0x40 << 0 //  Bit 6.
	WDGA CR_Bits = 0x01 << 7 //+ Activation bit.
)

const (
	W      CFR_Bits = 0x7F << 0 //+ W[6:0] bits (7-bit window value).
	W0     CFR_Bits = 0x01 << 0 //  Bit 0.
	W1     CFR_Bits = 0x02 << 0 //  Bit 1.
	W2     CFR_Bits = 0x04 << 0 //  Bit 2.
	W3     CFR_Bits = 0x08 << 0 //  Bit 3.
	W4     CFR_Bits = 0x10 << 0 //  Bit 4.
	W5     CFR_Bits = 0x20 << 0 //  Bit 5.
	W6     CFR_Bits = 0x40 << 0 //  Bit 6.
	WDGTB  CFR_Bits = 0x03 << 7 //+ WDGTB[1:0] bits (Timer Base).
	WDGTB0 CFR_Bits = 0x01 << 7 //  Bit 0.
	WDGTB1 CFR_Bits = 0x02 << 7 //  Bit 1.
	EWI    CFR_Bits = 0x01 << 9 //+ Early Wakeup Interrupt.
)

const (
	EWIF SR_Bits = 0x01 << 0 //+ Early Wakeup Interrupt Flag.
)
