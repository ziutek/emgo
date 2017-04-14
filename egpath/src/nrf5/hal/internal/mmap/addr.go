package mmap

// Base addresses for peripherals:
const (
	APB_BASE uintptr = 0x40000000 // accessed by APB,
	AHB_BASE uintptr = 0x50000000 // accessed by AHB.

	FICR_BASE uintptr = 0x10000000
)
