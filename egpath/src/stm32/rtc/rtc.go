// Package rtc gives an access to STM32 real time clock registers.
//
// Peripheral: Clock
// Instances:
//  RTC  0x40002800  APB1
// Registers:
//  0x00  TR       Time register.
//  0x04  DR       Date register.
//  0x08  CR       Control register.
//  0x0C  ISR      Initialization and status register.
//  0x10  PRER     Prescaler register.
//  0x14  WUTR     Wakeup timer register.
//  0x18  CALIBR   Calibration register.
//  0x1C  ALRMAR   Alarm A register.
//  0x20  ALRMBR   Alarm B register.
//  0x24  WPR      Write protection register.
//  0x28  SSR      Sub second register.
//  0x2C  SHIFTR   Shift control register.
//  0x30  TSTR     Time stamp time register.
//  0x34  TSDR     Time stamp date register.
//  0x38  TSSSR    Time stamp sub second register.
//  0x3C  CALR     Calibration register.
//  0x40  TAFCR    Tamper and alternate function configuration register.
//  0x44  ALRMASSR Alarm A sub second register.
//  0x48  ALRMBSSR Alarm B sub second register.
package rtc

import (
	"mmio"
	"unsafe"
)

const (
	SU  TR_Bits = 0xf << 0  // Second units in BCD format
	ST  TR_Bits = 0x7 << 4  // Second tens in BCD format
	MNU TR_Bits = 0xf << 8  // Minute units in BCD format
	MNT TR_Bits = 0x7 << 12 // Minute tens in BCD format
	HU  TR_Bits = 0xf << 16 // Hour units in BCD format
	HT  TR_Bits = 0x3 << 20 // Hour tens in BCD format
	PM  TR_Bits = 0x1 << 22 // AM/PM notation
)

const (
	DU  DR_Bits = 0xf << 0  // Date units in BCD format.
	DT  DR_Bits = 0x3 << 4  // Date tens in BCD format.
	MU  DR_Bits = 0xf << 8  // Month units in BCD format.
	MT  DR_Bits = 0x1 << 12 // Month tens in BCD format.
	WDU DR_Bits = 0x7 << 13 // Week day units.
	YU  DR_Bits = 0xf << 16 // Year units in BCD format.
	YT  DR_Bits = 0xf << 20 // Year tens in BCD format.
)

const (
	WUCKSEL CR_Bits = 7 << 0  // Wakeup clock selection
	TSEDGE  CR_Bits = 1 << 3  // Timestamp event active edge
	REFCKON CR_Bits = 1 << 4  // Reference clock detection enable (50 or 60 Hz)
	BYPSHAD CR_Bits = 1 << 5  //  Bypass the shadow registers
	FMT     CR_Bits = 1 << 6  // Hour format
	DCE     CR_Bits = 1 << 7  // Coarse digital calibration enable
	ALRAE   CR_Bits = 1 << 8  // Alarm A enable
	ALRBE   CR_Bits = 1 << 9  // Alarm B enable
	WUTE    CR_Bits = 1 << 10 // Wakeup timer enable.
	TSE     CR_Bits = 1 << 11 // Time stamp enable.
	ALRAIE  CR_Bits = 1 << 12 // Alarm A interrupt enable.
	ALRBIE  CR_Bits = 1 << 13 // Alarm B interrupt enable.
	WUTIE   CR_Bits = 1 << 14 // Wakeup timer interrupt enable.
	TSIE    CR_Bits = 1 << 15 // Timestamp interrupt enable.
	ADD1H   CR_Bits = 1 << 16 // Add 1 hour (summer time change).
	SUB1H   CR_Bits = 1 << 17 // Subtract 1 hour (winter time change).
	BKP     CR_Bits = 1 << 18 // Backup.
	COSEL   CR_Bits = 1 << 19 // Calibration output selection.
	POL     CR_Bits = 1 << 20 // Output polarit.
	OSEL    CR_Bits = 3 << 21 // Output selection.
	COE     CR_Bits = 1 << 23 //  Calibration output enable.
)

const (
	ALRAWF  ISR_Bits = 1 << 0  // Alarm A write flag.
	ALRBWF  ISR_Bits = 1 << 1  // Alarm B write flag.
	WUTWF   ISR_Bits = 1 << 2  // Wakeup timer write flag.
	SHPF    ISR_Bits = 1 << 3  // Shift operation pending.
	INITS   ISR_Bits = 1 << 4  // Initialization status flag.
	RSF     ISR_Bits = 1 << 5  // Registers synchronization flag.
	INITF   ISR_Bits = 1 << 6  // Initialization flag.
	INIT    ISR_Bits = 1 << 7  // Initialization mode.
	ALRAF   ISR_Bits = 1 << 8  // Alarm A flag.
	ALRBF   ISR_Bits = 1 << 9  // Alarm B flag.
	WUTF    ISR_Bits = 1 << 10 // Wakeup timer flag.
	TSF     ISR_Bits = 1 << 11 // Timestamp flag.
	TSOVF   ISR_Bits = 1 << 12 // Timestamp overflow flag.
	TAMP1F  ISR_Bits = 1 << 13 // Tamper detection flag.
	TAMP2F  ISR_Bits = 1 << 14 // TAMPER2 detection flag.
	RECALPF ISR_Bits = 1 << 16 // Recalibration pending Flag.
)

const (
	PREDIV_S PRER_Bits = 0x7fff << 0 // Synchronous prescaler factor.
	PREDIV_A PRER_Bits = 0x7f << 16  // Asynchronous prescaler factor.
)

const (
	DC  CALIBR_Bits = 0x1f << 0 // Digital calibration.
	DCS CALIBR_Bits = 1 << 7    // Digital calibration sign.
)

const (
	ASU    ALRMAR_Bits = 0xf << 0  // Second units in BCD format.
	AST    ALRMAR_Bits = 0x7 << 4  // Second tens in BCD format.
	AMSK1  ALRMAR_Bits = 0x1 << 7  //  Alarm A seconds mask.
	AMNU   ALRMAR_Bits = 0xf << 8  //  Minute units in BCD format.
	AMNT   ALRMAR_Bits = 0x7 << 12 // Minute tens in BCD format.
	AMSK2  ALRMAR_Bits = 0x1 << 15 // Alarm A minutes mask.
	AHU    ALRMAR_Bits = 0xf << 16 // Hour units in BCD format.
	AHT    ALRMAR_Bits = 0x3 << 20 // Hour tens in BCD format.
	APM    ALRMAR_Bits = 0x1 << 22 // AM/PM notation.
	AMSK3  ALRMAR_Bits = 0x1 << 23 // Alarm A hours mask.
	ADU    ALRMAR_Bits = 0xf << 24 // Date units or day in BCD format.
	ADT    ALRMAR_Bits = 0x3 << 28 // Date tens in BCD format.
	AWDSEL ALRMAR_Bits = 0x1 << 30 // Week day selection.
	AMSK4  ALRMAR_Bits = 0x1 << 31 // Alarm A date mask
)

const (
	BSU    ALRMBR_Bits = 0xf << 0  // Second units in BCD format.
	BST    ALRMBR_Bits = 0x7 << 4  // Second tens in BCD format.
	BMSK1  ALRMBR_Bits = 0x1 << 7  //  Alarm A seconds mask.
	BMNU   ALRMBR_Bits = 0xf << 8  //  Minute units in BCD format.
	BMNT   ALRMBR_Bits = 0x7 << 12 // Minute tens in BCD format.
	BMSK2  ALRMBR_Bits = 0x1 << 15 // Alarm A minutes mask.
	BHU    ALRMBR_Bits = 0xf << 16 // Hour units in BCD format.
	BHT    ALRMBR_Bits = 0x3 << 20 // Hour tens in BCD format.
	BPM    ALRMBR_Bits = 0x1 << 22 // AM/PM notation.
	BMSK3  ALRMBR_Bits = 0x1 << 23 // Alarm A hours mask.
	BDU    ALRMBR_Bits = 0xf << 24 // Date units or day in BCD format.
	BDT    ALRMBR_Bits = 0x3 << 28 // Date tens in BCD format.
	BWDSEL ALRMBR_Bits = 0x1 << 30 // Week day selection.
	BMSK4  ALRMBR_Bits = 0x1 << 31 // Alarm A date mask
)

const (
	SUBFS SHIFTR_Bits = 0x7fff << 0
	ADD1S SHIFTR_Bits = 1 << 31
)

const (
	TSU  TSTR_Bits = 0xf << 0  // Second units in BCD format.
	TST  TSTR_Bits = 0x7 << 4  // Second tens in BCD format.
	TMNU TSTR_Bits = 0xf << 8  // Minute units in BCD format.
	TMNT TSTR_Bits = 0x7 << 12 // Minute tens in BCD format.
	THU  TSTR_Bits = 0xf << 16 // Hour units in BCD format.
	THT  TSTR_Bits = 0x3 << 20 // Hour tens in BCD format.
	TPM  TSTR_Bits = 0x1 << 22 // AM/PM notation.
)

const (
	TDU  TSDR_Bits = 0xf << 0  // Date units in BCD format.
	TDT  TSDR_Bits = 0x3 << 4  // Date tens in BCD format.
	TMU  TSDR_Bits = 0xf << 8  // Month units in BCD format.
	TMT  TSDR_Bits = 0x1 << 12 // Month tens in BCD format.
	TWDU TSDR_Bits = 0x7 << 13 // Week day units.
)

const (
	CALM   CALR_Bits = 0x1ff << 0 // Calibration minus.
	CALW16 CALR_Bits = 1 << 13    // Use a 16-second calibration cycle period.
	CALW8  CALR_Bits = 1 << 14    // Use an 8-second calibration cycle period.
	CALP   CALR_Bits = 1 << 15    // Increase frequency of RTC by 488.5 ppm.
)

const (
	TAMP1E       TAFCR_Bits = 1 << 0  // Tamper 1 detection enable.
	TAMP1TRG     TAFCR_Bits = 1 << 1  // Active level for tamper 1.
	TAMPIE       TAFCR_Bits = 1 << 2  // Tamper interrupt enable.
	TAMP2E       TAFCR_Bits = 1 << 3  // Tamper 2 detection enable.
	TAMP2TRG     TAFCR_Bits = 1 << 4  // Active level for tamper 2.
	TAMPTS       TAFCR_Bits = 1 << 7  // Activate timestamp on tamper detection event.
	TAMPFREQ     TAFCR_Bits = 7 << 8  // Tamper sampling frequency.
	TAMPFLT      TAFCR_Bits = 3 << 11 // Tamper filter count.
	TAMPPRCH     TAFCR_Bits = 3 << 13 // Tamper precharge duration.
	TAMPPUDIS    TAFCR_Bits = 1 << 15 // TAMPER pull-up disable.
	TAMP1INSEL   TAFCR_Bits = 1 << 16 // TAMPER1 mapping.
	TSINSEL      TAFCR_Bits = 1 << 17 // TIMESTAMP mapping.
	ALARMOUTTYPE TAFCR_Bits = 1 << 18 // RTC_ALARM output type.
)

const (
	ASS     ALRMASSR_Bits = 0x7fff << 0 // Mask the most-significant bits starting at this bit.
	AMASKSS ALRMASSR_Bits = 0xf << 24   // Sub seconds value.
)

const (
	BSS     ALRMBSSR_Bits = 0x7fff << 0 // Mask the most-significant bits starting at this bit.
	BMASKSS ALRMBSSR_Bits = 0xf << 24   // Sub seconds value.
)

// BKPxR backup registers.
// BUG: 20 is the size for STM32F4. Other variants can have more or less
// backup registers.
var BKPxR = (*[20]mmio.U32)(unsafe.Pointer(uintptr(0x40002800 + 0x50)))
