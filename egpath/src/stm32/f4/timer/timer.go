// Package systick gives an access to STM32 timers registers.
//
// Instances of Timer:
//  TIM1: 0x40010000
//  TIM2: 0x40000000
//  TIM3: 0x40000400
//  TIM4: 0x40000800
// Registers:
//  0: CR1  Control register 1.
//  1: CR2  Control register 2.
//  2: SMCR Slave mode control register.
package timer

import "mmio"

const siz = mmio.Fsiz

const (
	CEN  CR1_Bits  = 1 << 0     // Counter enable.
	UDIS CR1_Bits  = 1 << 1     // Update disable.
	URS  CR1_Bits  = 1 << 2     // Update request source.
	OPM  CR1_Bits  = 1 << 3     // One pulse mode.
	DIR  CR1_Bits  = 1 << 4     // Direction.
	CMS  CR1_Field = 2<<siz + 5 // Center-aligned mode selection.
	ARPE CR1_Bits  = 1 << 7     // Auto-reload preload enable.
	CKD  CR1_Field = 2<<siz + 8 // Clock division
)

const (
	CCPC  CR2_Bits  = 1 << 0     // Capture/compare preloaded control.
	CCUS  CR2_Bits  = 1 << 2     // Capture/compare control update selection.
	CCDS  CR2_Bits  = 1 << 3     // Capture/compare DMS selection.
	MMS   CR2_Field = 3<<siz + 4 // Master mode selection.
	TI1S  CR2_Bits  = 1 << 7     // TI1 selection
	OIS1  CR2_Bits  = 1 << 8     // Output Idle state 1 (OC1 output).
	OIS1N CR2_Bits  = 1 << 9     // Output Idle state 1 (OC1N output).
	OIS2  CR2_Bits  = 1 << 10    // Output Idle state 2 (OC2 output).
	OIS2N CR2_Bits  = 1 << 11    // Output Idle state 2 (OC2N output).
	OIS3  CR2_Bits  = 1 << 12    // Output Idle state 3 (OC3 output).
	OIS3N CR2_Bits  = 1 << 13    // Output Idle state 3 (OC3N output).
	OIS4  CR2_Bits  = 1 << 14    // Output Idle state 4 (OC4 output).
)

const (
	SMS  SMCR_Field = 3<<siz + 0  // Slave mode selection.
	TS   SMCR_Field = 3<<siz + 4  // Trigger selection.
	MSM  SMCR_Bits  = 1 << 7      // Master/slave mode.
	ETF  SMCR_Field = 4<<siz + 8  // External trigger filter.
	ETPS SMCR_Field = 2<<siz + 12 // External trigger prescaler.
	ECE  SMCR_Bits  = 1 << 14     //  External clock enable.
	ETP  SMCR_Bits  = 1 << 15     // External trigger polarity.
)
