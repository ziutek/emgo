// Peripheral: RTC_Periph  Real-Time Clock.
// Instances:
//  RTC  mmap.RTC_BASE
// Registers:
//  0x00 32  TR       Time register.
//  0x04 32  DR       Date register.
//  0x08 32  CR       Control register.
//  0x0C 32  ISR      Initialization and status register.
//  0x10 32  PRER     Prescaler register.
//  0x14 32  WUTR     Wakeup timer register.
//  0x18 32  CALIBR   Calibration register.
//  0x1C 32  ALRMAR   Alarm A register.
//  0x20 32  ALRMBR   Alarm B register.
//  0x24 32  WPR      Write protection register.
//  0x28 32  SSR      Sub second register.
//  0x2C 32  SHIFTR   Shift control register.
//  0x30 32  TSTR     Time stamp time register.
//  0x34 32  TSDR     Time stamp date register.
//  0x38 32  TSSSR    Time-stamp sub second register.
//  0x3C 32  CALR     Calibration register.
//  0x40 32  TAFCR    Tamper and alternate function configuration register.
//  0x44 32  ALRMASSR
//  0x48 32  ALRMBSSR
//  0x50 32  BKP0R    Backup register 1.
//  0x54 32  BKP1R    Backup register 1.
//  0x58 32  BKP2R    Backup register 2.
//  0x5C 32  BKP3R    Backup register 3.
//  0x60 32  BKP4R    Backup register 4.
//  0x64 32  BKP5R    Backup register 5.
//  0x68 32  BKP6R    Backup register 6.
//  0x6C 32  BKP7R    Backup register 7.
//  0x70 32  BKP8R    Backup register 8.
//  0x74 32  BKP9R    Backup register 9.
//  0x78 32  BKP10R   Backup register 10.
//  0x7C 32  BKP11R   Backup register 11.
//  0x80 32  BKP12R   Backup register 12.
//  0x84 32  BKP13R   Backup register 13.
//  0x88 32  BKP14R   Backup register 14.
//  0x8C 32  BKP15R   Backup register 15.
//  0x90 32  BKP16R   Backup register 16.
//  0x94 32  BKP17R   Backup register 17.
//  0x98 32  BKP18R   Backup register 18.
//  0x9C 32  BKP19R   Backup register 19.
// Import:
//  stm32/o/f411xe/mmap
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
	MSK4  ALRMAR_Bits = 0x01 << 31 //+
	WDSEL ALRMAR_Bits = 0x01 << 30 //+
	DT    ALRMAR_Bits = 0x03 << 28 //+
	DT_0  ALRMAR_Bits = 0x01 << 28
	DT_1  ALRMAR_Bits = 0x02 << 28
	DU    ALRMAR_Bits = 0x0F << 24 //+
	DU_0  ALRMAR_Bits = 0x01 << 24
	DU_1  ALRMAR_Bits = 0x02 << 24
	DU_2  ALRMAR_Bits = 0x04 << 24
	DU_3  ALRMAR_Bits = 0x08 << 24
	MSK3  ALRMAR_Bits = 0x01 << 23 //+
	PM    ALRMAR_Bits = 0x01 << 22 //+
	HT    ALRMAR_Bits = 0x03 << 20 //+
	HT_0  ALRMAR_Bits = 0x01 << 20
	HT_1  ALRMAR_Bits = 0x02 << 20
	HU    ALRMAR_Bits = 0x0F << 16 //+
	HU_0  ALRMAR_Bits = 0x01 << 16
	HU_1  ALRMAR_Bits = 0x02 << 16
	HU_2  ALRMAR_Bits = 0x04 << 16
	HU_3  ALRMAR_Bits = 0x08 << 16
	MSK2  ALRMAR_Bits = 0x01 << 15 //+
	MNT   ALRMAR_Bits = 0x07 << 12 //+
	MNT_0 ALRMAR_Bits = 0x01 << 12
	MNT_1 ALRMAR_Bits = 0x02 << 12
	MNT_2 ALRMAR_Bits = 0x04 << 12
	MNU   ALRMAR_Bits = 0x0F << 8 //+
	MNU_0 ALRMAR_Bits = 0x01 << 8
	MNU_1 ALRMAR_Bits = 0x02 << 8
	MNU_2 ALRMAR_Bits = 0x04 << 8
	MNU_3 ALRMAR_Bits = 0x08 << 8
	MSK1  ALRMAR_Bits = 0x01 << 7 //+
	ST    ALRMAR_Bits = 0x07 << 4 //+
	ST_0  ALRMAR_Bits = 0x01 << 4
	ST_1  ALRMAR_Bits = 0x02 << 4
	ST_2  ALRMAR_Bits = 0x04 << 4
	SU    ALRMAR_Bits = 0x0F << 0 //+
	SU_0  ALRMAR_Bits = 0x01 << 0
	SU_1  ALRMAR_Bits = 0x02 << 0
	SU_2  ALRMAR_Bits = 0x04 << 0
	SU_3  ALRMAR_Bits = 0x08 << 0
)

const (
	MSK4  ALRMBR_Bits = 0x01 << 31 //+
	WDSEL ALRMBR_Bits = 0x01 << 30 //+
	DT    ALRMBR_Bits = 0x03 << 28 //+
	DT_0  ALRMBR_Bits = 0x01 << 28
	DT_1  ALRMBR_Bits = 0x02 << 28
	DU    ALRMBR_Bits = 0x0F << 24 //+
	DU_0  ALRMBR_Bits = 0x01 << 24
	DU_1  ALRMBR_Bits = 0x02 << 24
	DU_2  ALRMBR_Bits = 0x04 << 24
	DU_3  ALRMBR_Bits = 0x08 << 24
	MSK3  ALRMBR_Bits = 0x01 << 23 //+
	PM    ALRMBR_Bits = 0x01 << 22 //+
	HT    ALRMBR_Bits = 0x03 << 20 //+
	HT_0  ALRMBR_Bits = 0x01 << 20
	HT_1  ALRMBR_Bits = 0x02 << 20
	HU    ALRMBR_Bits = 0x0F << 16 //+
	HU_0  ALRMBR_Bits = 0x01 << 16
	HU_1  ALRMBR_Bits = 0x02 << 16
	HU_2  ALRMBR_Bits = 0x04 << 16
	HU_3  ALRMBR_Bits = 0x08 << 16
	MSK2  ALRMBR_Bits = 0x01 << 15 //+
	MNT   ALRMBR_Bits = 0x07 << 12 //+
	MNT_0 ALRMBR_Bits = 0x01 << 12
	MNT_1 ALRMBR_Bits = 0x02 << 12
	MNT_2 ALRMBR_Bits = 0x04 << 12
	MNU   ALRMBR_Bits = 0x0F << 8 //+
	MNU_0 ALRMBR_Bits = 0x01 << 8
	MNU_1 ALRMBR_Bits = 0x02 << 8
	MNU_2 ALRMBR_Bits = 0x04 << 8
	MNU_3 ALRMBR_Bits = 0x08 << 8
	MSK1  ALRMBR_Bits = 0x01 << 7 //+
	ST    ALRMBR_Bits = 0x07 << 4 //+
	ST_0  ALRMBR_Bits = 0x01 << 4
	ST_1  ALRMBR_Bits = 0x02 << 4
	ST_2  ALRMBR_Bits = 0x04 << 4
	SU    ALRMBR_Bits = 0x0F << 0 //+
	SU_0  ALRMBR_Bits = 0x01 << 0
	SU_1  ALRMBR_Bits = 0x02 << 0
	SU_2  ALRMBR_Bits = 0x04 << 0
	SU_3  ALRMBR_Bits = 0x08 << 0
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
	MASKSS   ALRMASSR_Bits = 0x0F << 24 //+
	MASKSS_0 ALRMASSR_Bits = 0x01 << 24
	MASKSS_1 ALRMASSR_Bits = 0x02 << 24
	MASKSS_2 ALRMASSR_Bits = 0x04 << 24
	MASKSS_3 ALRMASSR_Bits = 0x08 << 24
	SS       ALRMASSR_Bits = 0x7FFF << 0 //+
)

const (
	MASKSS   ALRMBSSR_Bits = 0x0F << 24 //+
	MASKSS_0 ALRMBSSR_Bits = 0x01 << 24
	MASKSS_1 ALRMBSSR_Bits = 0x02 << 24
	MASKSS_2 ALRMBSSR_Bits = 0x04 << 24
	MASKSS_3 ALRMBSSR_Bits = 0x08 << 24
	SS       ALRMBSSR_Bits = 0x7FFF << 0 //+
)
