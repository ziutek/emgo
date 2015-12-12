// Package systick gives an access to STM32 timers registers.
//
// Peripheral: Timer
// Instances:
//  TIM1: 0x40010000
//  TIM2: 0x40000000
//  TIM3: 0x40000400
//  TIM4: 0x40000800
// Registers:
//   0: CR1  Control register 1.
//   1: CR2  Control register 2.
//   2: SMCR Slave mode control register.
//   3: DIER DMA/interrupt enable register.
//   4: SR   Status register.
//   5: EGR  Event generation register.
package timer

const (
	CEN  CR1_Bits = 1 << 0 // Counter enable.
	UDIS CR1_Bits = 1 << 1 // Update disable.
	URS  CR1_Bits = 1 << 2 // Update request source.
	OPM  CR1_Bits = 1 << 3 // One pulse mode.
	DIR  CR1_Bits = 1 << 4 // Direction.
	CMS  CR1_Bits = 3 << 5 // Center-aligned mode selection.
	ARPE CR1_Bits = 1 << 7 // Auto-reload preload enable.
	CKD  CR1_Bits = 3 << 8 // Clock division
)

const (
	CCPC  CR2_Bits = 1 << 0  // Capture/compare preloaded control.
	CCUS  CR2_Bits = 1 << 2  // Capture/compare control update selection.
	CCDS  CR2_Bits = 1 << 3  // Capture/compare DMS selection.
	MMS   CR2_Bits = 7 << 4  // Master mode selection.
	TI1S  CR2_Bits = 1 << 7  // TI1 selection
	OIS1  CR2_Bits = 1 << 8  // Output Idle state 1 (OC1 output).
	OIS1N CR2_Bits = 1 << 9  // Output Idle state 1 (OC1N output).
	OIS2  CR2_Bits = 1 << 10 // Output Idle state 2 (OC2 output).
	OIS2N CR2_Bits = 1 << 11 // Output Idle state 2 (OC2N output).
	OIS3  CR2_Bits = 1 << 12 // Output Idle state 3 (OC3 output).
	OIS3N CR2_Bits = 1 << 13 // Output Idle state 3 (OC3N output).
	OIS4  CR2_Bits = 1 << 14 // Output Idle state 4 (OC4 output).
)

const (
	SMS  SMCR_Bits = 7 << 0   // Slave mode selection.
	TS   SMCR_Bits = 7 << 4   // Trigger selection.
	MSM  SMCR_Bits = 1 << 7   // Master/slave mode.
	ETF  SMCR_Bits = 0xf << 8 // External trigger filter.
	ETPS SMCR_Bits = 3 << 12  // External trigger prescaler.
	ECE  SMCR_Bits = 1 << 14  // External clock enable.
	ETP  SMCR_Bits = 1 << 15  // External trigger polarity.
)

const (
	UIE   DIER_Bits = 1 << 0  // Update interrupt enable.
	CC1IE DIER_Bits = 1 << 1  // Capture/Compare 1 interrupt enable.
	CC2IE DIER_Bits = 1 << 2  // Capture/Compare 2 interrupt enable.
	CC3IE DIER_Bits = 1 << 3  // Capture/Compare 3 interrupt enable.
	CC4IE DIER_Bits = 1 << 4  // Capture/Compare 4 interrupt enable.
	COMIE DIER_Bits = 1 << 5  // COM interrupt enable.
	TIE   DIER_Bits = 1 << 6  // Trigger interrupt enable.
	BIE   DIER_Bits = 1 << 7  // Break interrupt enable.
	UDE   DIER_Bits = 1 << 8  // Update DMA request enable.
	CC1DE DIER_Bits = 1 << 9  // Capture/Compare 1 DMA request enable.
	CC2DE DIER_Bits = 1 << 10 // Capture/Compare 2 DMA request enable.
	CC3DE DIER_Bits = 1 << 11 // Capture/Compare 3 DMA request enable.
	CC4DE DIER_Bits = 1 << 12 // Capture/Compare 4 DMA request enable.
	COMDE DIER_Bits = 1 << 13 // COM DMA request enable.
	TDE   DIER_Bits = 1 << 14 // Trigger DMA request enable.
)

const (
	UIF   SR_Bits = 1 << 0  // Update interrupt flag.
	CC1IF SR_Bits = 1 << 1  // Capture/Compare 1 interrupt flag.
	CC2IF SR_Bits = 1 << 2  // Capture/Compare 2 interrupt flag.
	CC3IF SR_Bits = 1 << 3  // Capture/Compare 3 interrupt flag.
	CC4IF SR_Bits = 1 << 4  // Capture/Compare 4 interrupt flag.
	COMIF SR_Bits = 1 << 5  // COM interrupt flag.
	TIF   SR_Bits = 1 << 6  // Trigger interrupt flag.
	BIF   SR_Bits = 1 << 7  // Break interrupt flag.
	CC1OF SR_Bits = 1 << 9  // Capture/Compare 1 overcapture flag.
	CC2OF SR_Bits = 1 << 10 // Capture/Compare 2 overcapture flag.
	CC3OF SR_Bits = 1 << 11 // Capture/Compare 4 overcapture flag.
	CC4OF SR_Bits = 1 << 12 // Capture/Compare 4 overcapture flag.
)

const (
	UG   EGR_Bits = 1 << 0 // Update generation.
	CC1G EGR_Bits = 1 << 1 // Capture/Compare 1 generation.
	CC2G EGR_Bits = 1 << 2 // Capture/Compare 2 generation.
	CC3G EGR_Bits = 1 << 3 // Capture/Compare 3 generation.
	CC4G EGR_Bits = 1 << 4 // Capture/Compare 4 generation.
	COMG EGR_Bits = 1 << 5 // Capture/Compare control update generation.
	TG   EGR_Bits = 1 << 6 // Trigger generation.
	BG   EGR_Bits = 1 << 7 // Break generation.
)
