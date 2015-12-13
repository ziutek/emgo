// Package systick gives an access to STM32 timers registers.
//
// Peripheral: Timer
// Instances:
//  TIM1:  0x40010000  APB2  Advanced
//  TIM2:  0x40000000  APB1  General 4ch
//  TIM3:  0x40000400  APB1  General 4ch
//  TIM4:  0x40000800  APB1  General 4ch
//  TIM5:  0x40000C00  APB1  General 4ch
//  TIM6:  0x40001000  APB1  Basic
//  TIM7:  0x40001400  APB1  Basic
//  TIM8:  0x40010400  APB2  Advanced
//  TIM9:  0x40014000  APB2  General 2ch
//  TIM10: 0x40014400  APB2  General 1ch
//  TIM11: 0x40014800  APB2  General 1ch
//  TIM12: 0x40001800  APB1  General 2ch
//  TIM13: 0x40001C00  APB1  General 1ch
//  TIM14: 0x40002000  APB1  General 1ch
// Registers:
//   0: CR1   Control register 1.
//   1: CR2   Control register 2.
//   2: SMCR  Slave mode control register.
//   3: DIER  DMA/interrupt enable register.
//   4: SR    Status register.
//   5: EGR   Event generation register.
//   6: CCMR1 Capture/compare mode register 1.
//   7: CCMR2 Capture/compare mode register 2.
//   8: CCER  Capture/compare enable register.
//   9: CNT   Counter.
//  10: PSC   Prescaler
//  11: ARR   Auto-reload register.
//  12: RCR   Repetition counter register.
//  13: CCR1  Capture/compare register 1.
//  14: CCR2  Capture/compare register 2.
//  15: CCR3  Capture/compare register 3.
//  16: CCR4  Capture/compare register 4.
//  17: BDTR  Break and dead-time register.
//  18: DCR   DMA control register.
//  19: DMAR  DMA address for full transfer.
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

const (
	CC1S  CCMR1_Bits = 3 << 0  // Capture/Compare 1 selection.
	OC1FE CCMR1_Bits = 1 << 2  // Output Compare 1 fast enable.
	OC1PE CCMR1_Bits = 1 << 3  // Output Compare 1 preload enable.
	OC1M  CCMR1_Bits = 7 << 4  // Output Compare 1 mode.
	OC1CE CCMR1_Bits = 1 << 7  // Output Compare 1 clear enable.
	CC2S  CCMR1_Bits = 3 << 8  // Capture/Compare 2 selection.
	OC2FE CCMR1_Bits = 1 << 10 // Output Compare 2 fast enable.
	OC2PE CCMR1_Bits = 1 << 11 // Output Compare 2 preload enable.
	OC2M  CCMR1_Bits = 7 << 12 // Output Compare 2 mode.
	OC2CE CCMR1_Bits = 1 << 15 // Output Compare 2 clear enable.
)
const (
	IC1PSC CCMR1_Bits = 3 << 2    // Input capture 1 prescaler.
	IC1F   CCMR1_Bits = 0xf << 4  // Input capture 1 filter.
	IC2PSC CCMR1_Bits = 3 << 10   // Input capture 2 prescaler.
	IC2F   CCMR1_Bits = 0xf << 12 // Input capture 2 filter.
)

const (
	CC3S  CCMR2_Bits = 3 << 0  // Capture/Compare 3 selection.
	OC3FE CCMR2_Bits = 1 << 2  // Output Compare 3 fast enable.
	OC3PE CCMR2_Bits = 1 << 3  // Output Compare 3 preload enable.
	OC3M  CCMR2_Bits = 7 << 4  // Output Compare 3 mode.
	OC3CE CCMR2_Bits = 1 << 7  // Output Compare 3 clear enable.
	CC4S  CCMR2_Bits = 3 << 8  // Capture/Compare 4 selection.
	OC4FE CCMR2_Bits = 1 << 10 // Output Compare 4 fast enable.
	OC4PE CCMR2_Bits = 1 << 11 // Output Compare 4 preload enable.
	OC4M  CCMR2_Bits = 7 << 12 // Output Compare 4 mode.
	OC4CE CCMR2_Bits = 1 << 15 // Output Compare 4 clear enable.
)
const (
	IC3PSC CCMR2_Bits = 3 << 2    // Input capture 3 prescaler.
	IC3F   CCMR2_Bits = 0xf << 4  // Input capture 3 filter.
	IC4PSC CCMR2_Bits = 3 << 10   // Input capture 4 prescaler.
	IC4F   CCMR2_Bits = 0xf << 12 // Input capture 4 filter.
)

const (
	CC1E  CCER_Bits = 1 << 0  // Capture/Compare 1 output enable.
	CC1P  CCER_Bits = 1 << 1  // Capture/Compare 1 output polarity.
	CC1NE CCER_Bits = 1 << 2  // Capture/Compare 1 complementary output enable.
	CC1NP CCER_Bits = 1 << 3  // Capture/Compare 1 complementary output polarity
	CC2E  CCER_Bits = 1 << 4  // Capture/Compare 2 output enable.
	CC2P  CCER_Bits = 1 << 5  // Capture/Compare 2 output polarity.
	CC2NE CCER_Bits = 1 << 6  // Capture/Compare 2 complementary output enable.
	CC2NP CCER_Bits = 1 << 7  // Capture/Compare 2 complementary output polarity
	CC3E  CCER_Bits = 1 << 8  // Capture/Compare 3 output enable.
	CC3P  CCER_Bits = 1 << 9  // Capture/Compare 3 output polarity.
	CC3NE CCER_Bits = 1 << 10 // Capture/Compare 3 complementary output enable.
	CC3NP CCER_Bits = 1 << 11 // Capture/Compare 3 complementary output polarity
	CC4E  CCER_Bits = 1 << 12 // Capture/Compare 4 output enable.
	CC4P  CCER_Bits = 1 << 13 // Capture/Compare 4 output polarity.
)

const (
	REP RCR_Bits = 0xff << 0 // Repetition counter value.
)

const (
	DTG  BDTR_Bits = 0xff << 0 // Dead-time generator setup.
	LOCK BDTR_Bits = 3 << 8    // Lock configuration.
	OSSI BDTR_Bits = 1 << 10   // Off-state selection for Idle mode.
	OSSR BDTR_Bits = 1 << 11   // Off-state selection for Run mode.
	BKE  BDTR_Bits = 1 << 12   // Break enable.
	BKP  BDTR_Bits = 1 << 13   // Break polarity.
	AOE  BDTR_Bits = 1 << 14   // Automatic output enable.
	MOE  BDTR_Bits = 1 << 15   // Main output enable.
)

const (
	DBA DCR_Bits = 0x1f << 0 // DMA base address.
	DBL DCR_Bits = 0x1f << 8 // DMA burst length.
)
