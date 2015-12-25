// Package mmap provides base memory adresses for all peripherals.
package mmap

const (
	FLASH_BASE     uintptr = 0x08000000 // FLASH base address in the alias region
	SRAM_BASE      uintptr = 0x20000000 // SRAM base address in the alias region
	PERIPH_BASE    uintptr = 0x40000000 // Peripheral base address in the alias region
	SRAM_BB_BASE   uintptr = 0x22000000 // SRAM base address in the bit-band region
	PERIPH_BB_BASE uintptr = 0x42000000 // Peripheral base address in the bit-band region
	FSMC_R_BASE    uintptr = 0xA0000000 // FSMC registers base address
)
