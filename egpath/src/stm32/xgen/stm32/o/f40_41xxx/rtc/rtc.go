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
//  0x3C 32  CALR       Calibration register.
//  0x40 32  TAFCR      Tamper and alternate function configuration register.
//  0x44 32  ALRMSSR[2] Alarm A, B subsecond registers.
//  0x50 32  BKPR[20]   Backup registers.
// Import:
//  stm32/o/f40_41xxx/mmap
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
	PREDIV_S PRER_Bits = 0x1FFF << 0 //+
)

const (
	WUT WUTR_Bits = 0xFFFF << 0 //+
)

const (
	DCS CALIBR_Bits = 0x01 << 7 //+
	DC  CALIBR_Bits = 0x1F << 0 //+
)

const (
	MSK4  ALRMR_Bits = 0x01 << 31 //+
	WDSEL ALRMR_Bits = 0x01 << 30 //+
	DT    ALRMR_Bits = 0x03 << 28 //+
	DT_0  ALRMR_Bits = 0x01 << 28
	DT_1  ALRMR_Bits = 0x02 << 28
	DU    ALRMR_Bits = 0x0F << 24 //+
	DU_0  ALRMR_Bits = 0x01 << 24
	DU_1  ALRMR_Bits = 0x02 << 24
	DU_2  ALRMR_Bits = 0x04 << 24
	DU_3  ALRMR_Bits = 0x08 << 24
	MSK3  ALRMR_Bits = 0x01 << 23 //+
	PM    ALRMR_Bits = 0x01 << 22 //+
	HT    ALRMR_Bits = 0x03 << 20 //+
	HT_0  ALRMR_Bits = 0x01 << 20
	HT_1  ALRMR_Bits = 0x02 << 20
	HU    ALRMR_Bits = 0x0F << 16 //+
	HU_0  ALRMR_Bits = 0x01 << 16
	HU_1  ALRMR_Bits = 0x02 << 16
	HU_2  ALRMR_Bits = 0x04 << 16
	HU_3  ALRMR_Bits = 0x08 << 16
	MSK2  ALRMR_Bits = 0x01 << 15 //+
	MNT   ALRMR_Bits = 0x07 << 12 //+
	MNT_0 ALRMR_Bits = 0x01 << 12
	MNT_1 ALRMR_Bits = 0x02 << 12
	MNT_2 ALRMR_Bits = 0x04 << 12
	MNU   ALRMR_Bits = 0x0F << 8 //+
	MNU_0 ALRMR_Bits = 0x01 << 8
	MNU_1 ALRMR_Bits = 0x02 << 8
	MNU_2 ALRMR_Bits = 0x04 << 8
	MNU_3 ALRMR_Bits = 0x08 << 8
	MSK1  ALRMR_Bits = 0x01 << 7 //+
	ST    ALRMR_Bits = 0x07 << 4 //+
	ST_0  ALRMR_Bits = 0x01 << 4
	ST_1  ALRMR_Bits = 0x02 << 4
	ST_2  ALRMR_Bits = 0x04 << 4
	SU    ALRMR_Bits = 0x0F << 0 //+
	SU_0  ALRMR_Bits = 0x01 << 0
	SU_1  ALRMR_Bits = 0x02 << 0
	SU_2  ALRMR_Bits = 0x04 << 0
	SU_3  ALRMR_Bits = 0x08 << 0
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
	PM    TSTR_Bits = 0x01 << 22 //+
	HT    TSTR_Bits = 0x03 << 20 //+
	HT_0  TSTR_Bits = 0x01 << 20
	HT_1  TSTR_Bits = 0x02 << 20
	HU    TSTR_Bits = 0x0F << 16 //+
	HU_0  TSTR_Bits = 0x01 << 16
	HU_1  TSTR_Bits = 0x02 << 16
	HU_2  TSTR_Bits = 0x04 << 16
	HU_3  TSTR_Bits = 0x08 << 16
	MNT   TSTR_Bits = 0x07 << 12 //+
	MNT_0 TSTR_Bits = 0x01 << 12
	MNT_1 TSTR_Bits = 0x02 << 12
	MNT_2 TSTR_Bits = 0x04 << 12
	MNU   TSTR_Bits = 0x0F << 8 //+
	MNU_0 TSTR_Bits = 0x01 << 8
	MNU_1 TSTR_Bits = 0x02 << 8
	MNU_2 TSTR_Bits = 0x04 << 8
	MNU_3 TSTR_Bits = 0x08 << 8
	ST    TSTR_Bits = 0x07 << 4 //+
	ST_0  TSTR_Bits = 0x01 << 4
	ST_1  TSTR_Bits = 0x02 << 4
	ST_2  TSTR_Bits = 0x04 << 4
	SU    TSTR_Bits = 0x0F << 0 //+
	SU_0  TSTR_Bits = 0x01 << 0
	SU_1  TSTR_Bits = 0x02 << 0
	SU_2  TSTR_Bits = 0x04 << 0
	SU_3  TSTR_Bits = 0x08 << 0
)

const (
	WDU   TSDR_Bits = 0x07 << 13 //+
	WDU_0 TSDR_Bits = 0x01 << 13
	WDU_1 TSDR_Bits = 0x02 << 13
	WDU_2 TSDR_Bits = 0x04 << 13
	MT    TSDR_Bits = 0x01 << 12 //+
	MU    TSDR_Bits = 0x0F << 8  //+
	MU_0  TSDR_Bits = 0x01 << 8
	MU_1  TSDR_Bits = 0x02 << 8
	MU_2  TSDR_Bits = 0x04 << 8
	MU_3  TSDR_Bits = 0x08 << 8
	DT    TSDR_Bits = 0x03 << 4 //+
	DT_0  TSDR_Bits = 0x01 << 4
	DT_1  TSDR_Bits = 0x02 << 4
	DU    TSDR_Bits = 0x0F << 0 //+
	DU_0  TSDR_Bits = 0x01 << 0
	DU_1  TSDR_Bits = 0x02 << 0
	DU_2  TSDR_Bits = 0x04 << 0
	DU_3  TSDR_Bits = 0x08 << 0
)

const (
	SS TSSSR_Bits = 0xFFFF << 0 //+
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
	TSINSEL      TAFCR_Bits = 0x01 << 17 //+
	TAMPINSEL    TAFCR_Bits = 0x01 << 16 //+
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
	TAMPIE       TAFCR_Bits = 0x01 << 2 //+
	TAMP1TRG     TAFCR_Bits = 0x01 << 1 //+
	TAMP1E       TAFCR_Bits = 0x01 << 0 //+
)

const (
	MASKSS   ALRMSSR_Bits = 0x0F << 24 //+
	MASKSS_0 ALRMSSR_Bits = 0x01 << 24
	MASKSS_1 ALRMSSR_Bits = 0x02 << 24
	MASKSS_2 ALRMSSR_Bits = 0x04 << 24
	MASKSS_3 ALRMSSR_Bits = 0x08 << 24
	SS       ALRMSSR_Bits = 0x7FFF << 0 //+
)
