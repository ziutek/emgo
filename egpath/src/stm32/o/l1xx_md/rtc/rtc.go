// +build l1xx_md

// Peripheral: RTC_Periph  Real-Time Clock.
// Instances:
//  RTC  mmap.RTC_BASE
// Registers:
//  0x00 32  TR         Time register.
//  0x04 32  DR         Date register.
//  0x08 32  CR         Control register.
//  0x0C 32  ISR        Initialization and status register.
//  0x10 32  PRER       Prescaler register.
//  0x14 32  WUTR       Wakeup timer register.
//  0x18 32  CALIBR     Calibration register.
//  0x1C 32  ALRMR[2]   Alarm A, B registers.
//  0x24 32  WPR        Write protection register.
//  0x28 32  SSR        Sub second register.
//  0x2C 32  SHIFTR     Shift control register.
//  0x30 32  TSTR       Time stamp time register.
//  0x34 32  TSDR       Time stamp date register.
//  0x38 32  TSSSR      Time-stamp sub second register.
//  0x3C 32  CALR       RRTC calibration register.
//  0x40 32  TAFCR      Tamper and alternate function configuration register.
//  0x44 32  ALRMSSR[2] Alarm A, B subsecond registers.
//  0x50 32  BKPR[32]   Backup registers.
// Import:
//  stm32/o/l1xx_md/mmap
package rtc

const (
	PM    TR_Bits = 0x01 << 22 //+
	HT    TR_Bits = 0x03 << 20 //+
	HT_0  TR_Bits = 0x01 << 20
	HT_1  TR_Bits = 0x02 << 20
	HU    TR_Bits = 0x0F << 16 //+
	HU_0  TR_Bits = 0x01 << 16
	HU_1  TR_Bits = 0x02 << 16
	HU_2  TR_Bits = 0x04 << 16
	HU_3  TR_Bits = 0x08 << 16
	MNT   TR_Bits = 0x07 << 12 //+
	MNT_0 TR_Bits = 0x01 << 12
	MNT_1 TR_Bits = 0x02 << 12
	MNT_2 TR_Bits = 0x04 << 12
	MNU   TR_Bits = 0x0F << 8 //+
	MNU_0 TR_Bits = 0x01 << 8
	MNU_1 TR_Bits = 0x02 << 8
	MNU_2 TR_Bits = 0x04 << 8
	MNU_3 TR_Bits = 0x08 << 8
	ST    TR_Bits = 0x07 << 4 //+
	ST_0  TR_Bits = 0x01 << 4
	ST_1  TR_Bits = 0x02 << 4
	ST_2  TR_Bits = 0x04 << 4
	SU    TR_Bits = 0x0F << 0 //+
	SU_0  TR_Bits = 0x01 << 0
	SU_1  TR_Bits = 0x02 << 0
	SU_2  TR_Bits = 0x04 << 0
	SU_3  TR_Bits = 0x08 << 0
)

const (
	YT    DR_Bits = 0x0F << 20 //+
	YT_0  DR_Bits = 0x01 << 20
	YT_1  DR_Bits = 0x02 << 20
	YT_2  DR_Bits = 0x04 << 20
	YT_3  DR_Bits = 0x08 << 20
	YU    DR_Bits = 0x0F << 16 //+
	YU_0  DR_Bits = 0x01 << 16
	YU_1  DR_Bits = 0x02 << 16
	YU_2  DR_Bits = 0x04 << 16
	YU_3  DR_Bits = 0x08 << 16
	WDU   DR_Bits = 0x07 << 13 //+
	WDU_0 DR_Bits = 0x01 << 13
	WDU_1 DR_Bits = 0x02 << 13
	WDU_2 DR_Bits = 0x04 << 13
	MT    DR_Bits = 0x01 << 12 //+
	MU    DR_Bits = 0x0F << 8  //+
	MU_0  DR_Bits = 0x01 << 8
	MU_1  DR_Bits = 0x02 << 8
	MU_2  DR_Bits = 0x04 << 8
	MU_3  DR_Bits = 0x08 << 8
	DT    DR_Bits = 0x03 << 4 //+
	DT_0  DR_Bits = 0x01 << 4
	DT_1  DR_Bits = 0x02 << 4
	DU    DR_Bits = 0x0F << 0 //+
	DU_0  DR_Bits = 0x01 << 0
	DU_1  DR_Bits = 0x02 << 0
	DU_2  DR_Bits = 0x04 << 0
	DU_3  DR_Bits = 0x08 << 0
)

const (
	COE       CR_Bits = 0x01 << 23 //+
	OSEL      CR_Bits = 0x03 << 21 //+
	OSEL_0    CR_Bits = 0x01 << 21
	OSEL_1    CR_Bits = 0x02 << 21
	POL       CR_Bits = 0x01 << 20 //+
	COSEL     CR_Bits = 0x01 << 19 //+
	BCK       CR_Bits = 0x01 << 18 //+
	SUB1H     CR_Bits = 0x01 << 17 //+
	ADD1H     CR_Bits = 0x01 << 16 //+
	TSIE      CR_Bits = 0x01 << 15 //+
	WUTIE     CR_Bits = 0x01 << 14 //+
	ALRBIE    CR_Bits = 0x01 << 13 //+
	ALRAIE    CR_Bits = 0x01 << 12 //+
	TSE       CR_Bits = 0x01 << 11 //+
	WUTE      CR_Bits = 0x01 << 10 //+
	ALRBE     CR_Bits = 0x01 << 9  //+
	ALRAE     CR_Bits = 0x01 << 8  //+
	DCE       CR_Bits = 0x01 << 7  //+
	FMT       CR_Bits = 0x01 << 6  //+
	BYPSHAD   CR_Bits = 0x01 << 5  //+
	REFCKON   CR_Bits = 0x01 << 4  //+
	TSEDGE    CR_Bits = 0x01 << 3  //+
	WUCKSEL   CR_Bits = 0x07 << 0  //+
	WUCKSEL_0 CR_Bits = 0x01 << 0
	WUCKSEL_1 CR_Bits = 0x02 << 0
	WUCKSEL_2 CR_Bits = 0x04 << 0
)

const (
	RECALPF ISR_Bits = 0x01 << 16 //+
	TAMP3F  ISR_Bits = 0x01 << 15 //+
	TAMP2F  ISR_Bits = 0x01 << 14 //+
	TAMP1F  ISR_Bits = 0x01 << 13 //+
	TSOVF   ISR_Bits = 0x01 << 12 //+
	TSF     ISR_Bits = 0x01 << 11 //+
	WUTF    ISR_Bits = 0x01 << 10 //+
	ALRBF   ISR_Bits = 0x01 << 9  //+
	ALRAF   ISR_Bits = 0x01 << 8  //+
	INIT    ISR_Bits = 0x01 << 7  //+
	INITF   ISR_Bits = 0x01 << 6  //+
	RSF     ISR_Bits = 0x01 << 5  //+
	INITS   ISR_Bits = 0x01 << 4  //+
	SHPF    ISR_Bits = 0x01 << 3  //+
	WUTWF   ISR_Bits = 0x01 << 2  //+
	ALRBWF  ISR_Bits = 0x01 << 1  //+
	ALRAWF  ISR_Bits = 0x01 << 0  //+
)

const (
	PREDIV_A PRER_Bits = 0x7F << 16  //+
	PREDIV_S PRER_Bits = 0x7FFF << 0 //+
)

const (
	WUT WUTR_Bits = 0xFFFF << 0 //+
)

const (
	DCS CALIBR_Bits = 0x01 << 7 //+
	DC  CALIBR_Bits = 0x1F << 0 //+
)

const (
	AMSK4  ALRMR_Bits = 0x01 << 31 //+
	AWDSEL ALRMR_Bits = 0x01 << 30 //+
	ADT    ALRMR_Bits = 0x03 << 28 //+
	ADT_0  ALRMR_Bits = 0x01 << 28
	ADT_1  ALRMR_Bits = 0x02 << 28
	ADU    ALRMR_Bits = 0x0F << 24 //+
	ADU_0  ALRMR_Bits = 0x01 << 24
	ADU_1  ALRMR_Bits = 0x02 << 24
	ADU_2  ALRMR_Bits = 0x04 << 24
	ADU_3  ALRMR_Bits = 0x08 << 24
	AMSK3  ALRMR_Bits = 0x01 << 23 //+
	APM    ALRMR_Bits = 0x01 << 22 //+
	AHT    ALRMR_Bits = 0x03 << 20 //+
	AHT_0  ALRMR_Bits = 0x01 << 20
	AHT_1  ALRMR_Bits = 0x02 << 20
	AHU    ALRMR_Bits = 0x0F << 16 //+
	AHU_0  ALRMR_Bits = 0x01 << 16
	AHU_1  ALRMR_Bits = 0x02 << 16
	AHU_2  ALRMR_Bits = 0x04 << 16
	AHU_3  ALRMR_Bits = 0x08 << 16
	AMSK2  ALRMR_Bits = 0x01 << 15 //+
	AMNT   ALRMR_Bits = 0x07 << 12 //+
	AMNT_0 ALRMR_Bits = 0x01 << 12
	AMNT_1 ALRMR_Bits = 0x02 << 12
	AMNT_2 ALRMR_Bits = 0x04 << 12
	AMNU   ALRMR_Bits = 0x0F << 8 //+
	AMNU_0 ALRMR_Bits = 0x01 << 8
	AMNU_1 ALRMR_Bits = 0x02 << 8
	AMNU_2 ALRMR_Bits = 0x04 << 8
	AMNU_3 ALRMR_Bits = 0x08 << 8
	AMSK1  ALRMR_Bits = 0x01 << 7 //+
	AST    ALRMR_Bits = 0x07 << 4 //+
	AST_0  ALRMR_Bits = 0x01 << 4
	AST_1  ALRMR_Bits = 0x02 << 4
	AST_2  ALRMR_Bits = 0x04 << 4
	ASU    ALRMR_Bits = 0x0F << 0 //+
	ASU_0  ALRMR_Bits = 0x01 << 0
	ASU_1  ALRMR_Bits = 0x02 << 0
	ASU_2  ALRMR_Bits = 0x04 << 0
	ASU_3  ALRMR_Bits = 0x08 << 0
)

const (
	KEY WPR_Bits = 0xFF << 0 //+
)

const (
	SS SSR_Bits = 0xFFFF << 0 //+
)

const (
	SUBFS SHIFTR_Bits = 0x7FFF << 0 //+
	ADD1S SHIFTR_Bits = 0x01 << 31  //+
)

const (
	TPM    TSTR_Bits = 0x01 << 22 //+
	THT    TSTR_Bits = 0x03 << 20 //+
	THT_0  TSTR_Bits = 0x01 << 20
	THT_1  TSTR_Bits = 0x02 << 20
	THU    TSTR_Bits = 0x0F << 16 //+
	THU_0  TSTR_Bits = 0x01 << 16
	THU_1  TSTR_Bits = 0x02 << 16
	THU_2  TSTR_Bits = 0x04 << 16
	THU_3  TSTR_Bits = 0x08 << 16
	TMNT   TSTR_Bits = 0x07 << 12 //+
	TMNT_0 TSTR_Bits = 0x01 << 12
	TMNT_1 TSTR_Bits = 0x02 << 12
	TMNT_2 TSTR_Bits = 0x04 << 12
	TMNU   TSTR_Bits = 0x0F << 8 //+
	TMNU_0 TSTR_Bits = 0x01 << 8
	TMNU_1 TSTR_Bits = 0x02 << 8
	TMNU_2 TSTR_Bits = 0x04 << 8
	TMNU_3 TSTR_Bits = 0x08 << 8
	TST    TSTR_Bits = 0x07 << 4 //+
	TST_0  TSTR_Bits = 0x01 << 4
	TST_1  TSTR_Bits = 0x02 << 4
	TST_2  TSTR_Bits = 0x04 << 4
	TSU    TSTR_Bits = 0x0F << 0 //+
	TSU_0  TSTR_Bits = 0x01 << 0
	TSU_1  TSTR_Bits = 0x02 << 0
	TSU_2  TSTR_Bits = 0x04 << 0
	TSU_3  TSTR_Bits = 0x08 << 0
)

const (
	TWDU   TSDR_Bits = 0x07 << 13 //+
	TWDU_0 TSDR_Bits = 0x01 << 13
	TWDU_1 TSDR_Bits = 0x02 << 13
	TWDU_2 TSDR_Bits = 0x04 << 13
	TMT    TSDR_Bits = 0x01 << 12 //+
	TMU    TSDR_Bits = 0x0F << 8  //+
	TMU_0  TSDR_Bits = 0x01 << 8
	TMU_1  TSDR_Bits = 0x02 << 8
	TMU_2  TSDR_Bits = 0x04 << 8
	TMU_3  TSDR_Bits = 0x08 << 8
	TDT    TSDR_Bits = 0x03 << 4 //+
	TDT_0  TSDR_Bits = 0x01 << 4
	TDT_1  TSDR_Bits = 0x02 << 4
	TDU    TSDR_Bits = 0x0F << 0 //+
	TDU_0  TSDR_Bits = 0x01 << 0
	TDU_1  TSDR_Bits = 0x02 << 0
	TDU_2  TSDR_Bits = 0x04 << 0
	TDU_3  TSDR_Bits = 0x08 << 0
)

const (
	TSS TSSSR_Bits = 0xFFFF << 0 //+
)

const (
	CALP   CALR_Bits = 0x01 << 15 //+
	CALW8  CALR_Bits = 0x01 << 14 //+
	CALW16 CALR_Bits = 0x01 << 13 //+
	CALM   CALR_Bits = 0x1FF << 0 //+
	CALM_0 CALR_Bits = 0x01 << 0
	CALM_1 CALR_Bits = 0x02 << 0
	CALM_2 CALR_Bits = 0x04 << 0
	CALM_3 CALR_Bits = 0x08 << 0
	CALM_4 CALR_Bits = 0x10 << 0
	CALM_5 CALR_Bits = 0x20 << 0
	CALM_6 CALR_Bits = 0x40 << 0
	CALM_7 CALR_Bits = 0x80 << 0
	CALM_8 CALR_Bits = 0x100 << 0
)

const (
	ALARMOUTTYPE TAFCR_Bits = 0x01 << 18 //+
	TAMPPUDIS    TAFCR_Bits = 0x01 << 15 //+
	TAMPPRCH     TAFCR_Bits = 0x03 << 13 //+
	TAMPPRCH_0   TAFCR_Bits = 0x01 << 13
	TAMPPRCH_1   TAFCR_Bits = 0x02 << 13
	TAMPFLT      TAFCR_Bits = 0x03 << 11 //+
	TAMPFLT_0    TAFCR_Bits = 0x01 << 11
	TAMPFLT_1    TAFCR_Bits = 0x02 << 11
	TAMPFREQ     TAFCR_Bits = 0x07 << 8 //+
	TAMPFREQ_0   TAFCR_Bits = 0x01 << 8
	TAMPFREQ_1   TAFCR_Bits = 0x02 << 8
	TAMPFREQ_2   TAFCR_Bits = 0x04 << 8
	TAMPTS       TAFCR_Bits = 0x01 << 7 //+
	TAMP3TRG     TAFCR_Bits = 0x01 << 6 //+
	TAMP3E       TAFCR_Bits = 0x01 << 5 //+
	TAMP2TRG     TAFCR_Bits = 0x01 << 4 //+
	TAMP2E       TAFCR_Bits = 0x01 << 3 //+
	TAMPIE       TAFCR_Bits = 0x01 << 2 //+
	TAMP1TRG     TAFCR_Bits = 0x01 << 1 //+
	TAMP1E       TAFCR_Bits = 0x01 << 0 //+
)

const (
	AMASKSS   ALRMSSR_Bits = 0x0F << 24 //+
	AMASKSS_0 ALRMSSR_Bits = 0x01 << 24
	AMASKSS_1 ALRMSSR_Bits = 0x02 << 24
	AMASKSS_2 ALRMSSR_Bits = 0x04 << 24
	AMASKSS_3 ALRMSSR_Bits = 0x08 << 24
	ASS       ALRMSSR_Bits = 0x7FFF << 0 //+
)
