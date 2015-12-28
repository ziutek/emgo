// +build f40_41xxx

// Peripheral: RNG_Periph  RNG.
// Instances:
//  RNG  mmap.RNG_BASE
// Registers:
//  0x00 32  CR Control register.
//  0x04 32  SR Status register.
//  0x08 32  DR Data register.
// Import:
//  stm32/o/f40_41xxx/mmap
package rng

const (
	RNGEN CR_Bits = 0x01 << 2 //+
	IE    CR_Bits = 0x01 << 3 //+
)

const (
	DRDY SR_Bits = 0x01 << 0 //+
	CECS SR_Bits = 0x01 << 1 //+
	SECS SR_Bits = 0x01 << 2 //+
	CEIS SR_Bits = 0x01 << 5 //+
	SEIS SR_Bits = 0x01 << 6 //+
)
