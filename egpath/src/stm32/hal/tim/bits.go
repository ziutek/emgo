package tim

import (
	"stm32/hal/raw/tim"
)

type (
	CR1   = tim.CR1
	CR2   = tim.CR2
	SMCR  = tim.SMCR
	DIER  = tim.DIER
	SR    = tim.SR
	EGR   = tim.EGR
	CCMR1 = tim.CCMR1
	CCMR2 = tim.CCMR2
	CCER  = tim.CCER
	CNT   = tim.CNT
	PSC   = tim.PSC
	ARR   = tim.ARR
	RCR   = tim.RCR
	CCR1  = tim.CCR1
	CCR2  = tim.CCR2
	CCR3  = tim.CCR3
	CCR4  = tim.CCR4
	BDTR  = tim.BDTR
	DCR   = tim.DCR
	DMAR  = tim.DMAR
)

type (
	RCR1   = tim.RCR1
	RCR2   = tim.RCR2
	RSMCR  = tim.RSMCR
	RDIER  = tim.RDIER
	RSR    = tim.RSR
	REGR   = tim.REGR
	RCCMR1 = tim.RCCMR1
	RCCMR2 = tim.RCCMR2
	RCCER  = tim.RCCER
	RCNT   = tim.RCNT
	RPSC   = tim.RPSC
	RARR   = tim.RARR
	RRCR   = tim.RRCR
	RCCR1  = tim.RCCR1
	RCCR2  = tim.RCCR2
	RCCR3  = tim.RCCR3
	RCCR4  = tim.RCCR4
	RBDTR  = tim.RBDTR
	RDCR   = tim.RDCR
	RDMAR  = tim.RDMAR
)

type (
	RMCR1   = tim.RMCR1
	RMCR2   = tim.RMCR2
	RMSMCR  = tim.RMSMCR
	RMDIER  = tim.RMDIER
	RMSR    = tim.RMSR
	RMEGR   = tim.RMEGR
	RMCCMR1 = tim.RMCCMR1
	RMCCMR2 = tim.RMCCMR2
	RMCCER  = tim.RMCCER
	RMCNT   = tim.RMCNT
	RMPSC   = tim.RMPSC
	RMARR   = tim.RMARR
	RMRCR   = tim.RMRCR
	RMCCR1  = tim.RMCCR1
	RMCCR2  = tim.RMCCR2
	RMCCR3  = tim.RMCCR3
	RMCCR4  = tim.RMCCR4
	RMBDTR  = tim.RMBDTR
	RMDCR   = tim.RMDCR
	RMDMAR  = tim.RMDMAR
)

const (
	CEN  CR1 = 0x01 << 0 //+ Counter enable.
	UDIS CR1 = 0x01 << 1 //+ Update disable.
	URS  CR1 = 0x01 << 2 //+ Update request source.
	OPM  CR1 = 0x01 << 3 //+ One pulse mode.
	DIR  CR1 = 0x01 << 4 //+ Counter direction: 0: up (CNT+), 1: down (CNT-).
	CMS  CR1 = 0x03 << 5 //+ Center-aligned mode selection:
	EAM  CR1 = 0x00 << 5 //  edge-aligned mode (DIR sets direction),
	CAM1 CR1 = 0x01 << 5 //  center-aligned mode 1: CCxF set when counting down,
	CAM2 CR1 = 0x01 << 5 //  center-aligned mode 2: CCxF set when counting up,
	CAM3 CR1 = 0x01 << 5 //  center-aligned mode 3: CCxF set on any direction.
	ARPE CR1 = 0x01 << 7 //+ Auto-reload preload enable.
	CKD  CR1 = 0x03 << 8 //+ Clock division (used by digital filters):
	CKD1 CR1 = 0x00 << 8 //  t_DTS = t_CK_INT,
	CKD2 CR1 = 0x01 << 8 //  t_DTS = 2*t_CK_INT,
	CKD4 CR1 = 0x02 << 8 //  t_DTS = 4*t_CK_INT.
)

const (
	CENn  = 0
	UDISn = 1
	URSn  = 2
	OPMn  = 3
	DIRn  = 4
	CMSn  = 5
	ARPEn = 7
	CKDn  = 8
)

const (
	CCPC  CR2 = 0x01 << 0  //+ Capture/Compare Preloaded Control.
	CCUS  CR2 = 0x01 << 2  //+ Capture/Compare Control Update Selection.
	CCDS  CR2 = 0x01 << 3  //+ Capture/Compare DMA Selection.
	MMS   CR2 = 0x07 << 4  //+ Master Mode Selection (what is used as TRGO):
	MMUG  CR2 = 0x00 << 4  //  EGR.UG bit,
	MMEN  CR2 = 0x01 << 4  //  CNT_EN = CR1.CEN | TRGI,
	MMUPD CR2 = 0x02 << 4  //  update event,
	MMCC1 CR2 = 0x03 << 4  //  SR.CC1IF (enven if not cleared),
	MMOC1 CR2 = 0x04 << 4  //  OC1REF,
	MMOC2 CR2 = 0x05 << 4  //  OC2REF,
	MMOC3 CR2 = 0x06 << 4  //  OC3REF,
	MMOC4 CR2 = 0x07 << 4  //  OC4REF.
	TI1S  CR2 = 0x01 << 7  //+ TI1 Selection.
	OIS1  CR2 = 0x01 << 8  //+ Output Idle state 1 (OC1 output).
	OIS1N CR2 = 0x01 << 9  //+ Output Idle state 1 (OC1N output).
	OIS2  CR2 = 0x01 << 10 //+ Output Idle state 2 (OC2 output).
	OIS2N CR2 = 0x01 << 11 //+ Output Idle state 2 (OC2N output).
	OIS3  CR2 = 0x01 << 12 //+ Output Idle state 3 (OC3 output).
	OIS3N CR2 = 0x01 << 13 //+ Output Idle state 3 (OC3N output).
	OIS4  CR2 = 0x01 << 14 //+ Output Idle state 4 (OC4 output).
)

const (
	CCPCn  = 0
	CCUSn  = 2
	CCDSn  = 3
	MMSn   = 4
	TI1Sn  = 7
	OIS1n  = 8
	OIS1Nn = 9
	OIS2n  = 10
	OIS2Nn = 11
	OIS3n  = 12
	OIS3Nn = 13
	OIS4n  = 14
)

// Input filter configuration. Use for SMCR.ETF, CCMRx.ICxF.
const (
	FDIS   = 0x00 //  No filter, sampling at f_DTS.
	FCN2   = 0x01 //  f_SAMPLING = f_CK_INT, N = 2.
	FCN4   = 0x02 //  f_SAMPLING = f_CK_INT, N = 4.
	FCN8   = 0x03 //  f_SAMPLING = f_CK_INT, N = 8.
	FD2N6  = 0x04 //  f_SAMPLING = f_DTS/2, N = 6.
	FD2N8  = 0x05 //  f_SAMPLING = f_DTS/2, N = 8.
	FD4N6  = 0x06 //  f_SAMPLING = f_DTS/4, N = 6.
	FD4N8  = 0x07 //  f_SAMPLING = f_DTS/4, N = 8.
	FD8N6  = 0x08 //  f_SAMPLING = f_DTS/8, N = 6.
	FD8N8  = 0x09 //  f_SAMPLING = f_DTS/8, N = 8.
	FD16N5 = 0x0A //  f_SAMPLING = f_DTS/16, N = 5.
	FD16N6 = 0x0B //  f_SAMPLING = f_DTS/16, N = 6.
	FD16N8 = 0x0C //  f_SAMPLING = f_DTS/16, N = 8.
	FD32N5 = 0x0D //  f_SAMPLING = f_DTS/32, N = 5.
	FD32N6 = 0x0E //  f_SAMPLING = f_DTS/32, N = 6.
	FD32N8 = 0x0F //  f_SAMPLING = f_DTS/32, N = 8.
)

const (
	SMS    SMCR = 0x07 << 0  //+ Slave mode selection:
	SMDIS  SMCR = 0x01 << 0  //  slave mode disabled,
	SME1   SMCR = 0x01 << 0  //  encoder mode 1: counts on TI2FP2 edge,
	SME2   SMCR = 0x02 << 0  //  encoder mode 2: counts on TI1FP1 edge,
	SME3   SMCR = 0x03 << 0  //  encoder mode 3: counts on TI1FP1, TI2FP2 edge.
	SMRST  SMCR = 0x04 << 0  //  reset mode: TRGI reinitializes the counter,
	SMGTD  SMCR = 0x05 << 0  //  gated mode: TRGI high enables counter clock,
	SMTRG  SMCR = 0x06 << 0  //  trigger mode: TRGI starts the counter,
	SMEXT  SMCR = 0x07 << 0  //  rxternal clock mode: TRGI clocks the counter.
	TS     SMCR = 0x07 << 4  //+ Trigger selection (what is used as TRGI):
	ITR0   SMCR = 0x00 << 4  //  internal trigger 0,
	ITR1   SMCR = 0x01 << 4  //  internal trigger 1,
	ITR2   SMCR = 0x02 << 4  //  internal trigger 2,
	ITR3   SMCR = 0x03 << 4  //  internal trigger 3,
	TI1FED SMCR = 0x04 << 4  //  TI1 edge detector,
	TI1FP1 SMCR = 0x05 << 4  //  filtered timer input 1,
	TI1FP2 SMCR = 0x06 << 4  //  filtered timer input 2,
	ETRF   SMCR = 0x07 << 4  //  external trigger input.
	MSM    SMCR = 0x01 << 7  //+ Master/slave mode.
	ETF    SMCR = 0x0F << 8  //+ External trigg. filter, need N valid samples:
	ETPS   SMCR = 0x03 << 12 //+ External trigger prescaler:
	ETPS1  SMCR = 0x00 << 12 //  prescaller off,
	ETPS2  SMCR = 0x01 << 12 //  ETRP /= 2,
	ETPS4  SMCR = 0x02 << 12 //  ETRP /= 4,
	ETPS8  SMCR = 0x03 << 12 //  ETRP /= 8.
	ECE    SMCR = 0x01 << 14 //+ External clock enable.
	ETP    SMCR = 0x01 << 15 //+ External trigger polarity.
)

const (
	SMSn  = 0
	TSn   = 4
	MSMn  = 7
	ETFn  = 8
	ETPSn = 12
	ECEn  = 14
	ETPn  = 15
)

const (
	UIE   DIER = 0x01 << 0  //+ Update interrupt enable.
	CC1IE DIER = 0x01 << 1  //+ Capture/Compare 1 interrupt enable.
	CC2IE DIER = 0x01 << 2  //+ Capture/Compare 2 interrupt enable.
	CC3IE DIER = 0x01 << 3  //+ Capture/Compare 3 interrupt enable.
	CC4IE DIER = 0x01 << 4  //+ Capture/Compare 4 interrupt enable.
	COMIE DIER = 0x01 << 5  //+ COM interrupt enable.
	TIE   DIER = 0x01 << 6  //+ Trigger interrupt enable.
	BIE   DIER = 0x01 << 7  //+ Break interrupt enable.
	UDE   DIER = 0x01 << 8  //+ Update DMA request enable.
	CC1DE DIER = 0x01 << 9  //+ Capture/Compare 1 DMA request enable.
	CC2DE DIER = 0x01 << 10 //+ Capture/Compare 2 DMA request enable.
	CC3DE DIER = 0x01 << 11 //+ Capture/Compare 3 DMA request enable.
	CC4DE DIER = 0x01 << 12 //+ Capture/Compare 4 DMA request enable.
	COMDE DIER = 0x01 << 13 //+ COM DMA request enable.
	TDE   DIER = 0x01 << 14 //+ Trigger DMA request enable.
)

const (
	UIEn   = 0
	CC1IEn = 1
	CC2IEn = 2
	CC3IEn = 3
	CC4IEn = 4
	COMIEn = 5
	TIEn   = 6
	BIEn   = 7
	UDEn   = 8
	CC1DEn = 9
	CC2DEn = 10
	CC3DEn = 11
	CC4DEn = 12
	COMDEn = 13
	TDEn   = 14
)

const (
	UIF   SR = 0x01 << 0  //+ Update interrupt Flag.
	CC1IF SR = 0x01 << 1  //+ Capture/Compare 1 interrupt Flag.
	CC2IF SR = 0x01 << 2  //+ Capture/Compare 2 interrupt Flag.
	CC3IF SR = 0x01 << 3  //+ Capture/Compare 3 interrupt Flag.
	CC4IF SR = 0x01 << 4  //+ Capture/Compare 4 interrupt Flag.
	COMIF SR = 0x01 << 5  //+ COM interrupt Flag.
	TIF   SR = 0x01 << 6  //+ Trigger interrupt Flag.
	BIF   SR = 0x01 << 7  //+ Break interrupt Flag.
	CC1OF SR = 0x01 << 9  //+ Capture/Compare 1 Overcapture Flag.
	CC2OF SR = 0x01 << 10 //+ Capture/Compare 2 Overcapture Flag.
	CC3OF SR = 0x01 << 11 //+ Capture/Compare 3 Overcapture Flag.
	CC4OF SR = 0x01 << 12 //+ Capture/Compare 4 Overcapture Flag.
)

const (
	UIFn   = 0
	CC1IFn = 1
	CC2IFn = 2
	CC3IFn = 3
	CC4IFn = 4
	COMIFn = 5
	TIFn   = 6
	BIFn   = 7
	CC1OFn = 9
	CC2OFn = 10
	CC3OFn = 11
	CC4OFn = 12
)

const (
	UG   EGR = 0x01 << 0 //+ Update Generation.
	CC1G EGR = 0x01 << 1 //+ Capture/Compare 1 Generation.
	CC2G EGR = 0x01 << 2 //+ Capture/Compare 2 Generation.
	CC3G EGR = 0x01 << 3 //+ Capture/Compare 3 Generation.
	CC4G EGR = 0x01 << 4 //+ Capture/Compare 4 Generation.
	COMG EGR = 0x01 << 5 //+ Capture/Compare Control Update Generation.
	TG   EGR = 0x01 << 6 //+ Trigger Generation.
	BG   EGR = 0x01 << 7 //+ Break Generation.
)

const (
	UGn   = 0
	CC1Gn = 1
	CC2Gn = 2
	CC3Gn = 3
	CC4Gn = 4
	COMGn = 5
	TGn   = 6
	BGn   = 7
)

// Output Compare configuration. Use for CCMRx.OCxM.
const (
	OCFR   = 0x00 // OCxREF is frozen.
	OCHM   = 0x01 // OCxREF set to high on match.
	OCLM   = 0x02 // OCxREF set to low on match.
	OCTG   = 0x03 // OCxREF toggles on match.
	OCLO   = 0x04 // OCxREF forced low.
	OCHI   = 0x05 // OCxREF forced high.
	OCPWM1 = 0x06 // PWM mode 1: OCxREF = CNT+ < CCR1, CNT- ≤ CCR1.
	OCPWM2 = 0x07 // PWM mode 2: OCxREF = CNT+ ≥ CCR1, CNT- > CCR1.
)

// Input Capture prescaler configuration. Use for CCMRx.ICxPSC.
const (
	ICPSC1 = 0x00 // No prescaler.
	ICPSC2 = 0x01 // Capture once every 2 events.
	ICPSC4 = 0x02 // Capture once every 4 events.
	ICPSC8 = 0x03 // Capture once every 8 events.
)

const (
	CC1S   CCMR1 = 0x03 << 0 //+ Capture/Compare 1 Selection:
	CC1OUT CCMR1 = 0x00 << 0 //  as output,
	CC1TI1 CCMR1 = 0x01 << 0 //  input from TI1,
	CC1TI2 CCMR1 = 0x02 << 0 //  input from TI2,
	CC1TRC CCMR1 = 0x03 << 0 //  input from TRC.

	OC1FE  CCMR1 = 0x01 << 2 //+ Output Compare 1 Fast enable.
	OC1PE  CCMR1 = 0x01 << 3 //+ Output Compare 1 Preload enable.
	OC1M   CCMR1 = 0x07 << 4 //+ Output Compare 1 Mode.
	OC1CE  CCMR1 = 0x01 << 7 //+ Output Compare 1 Clear Enable.
	IC1PSC CCMR1 = 0x03 << 2 //+ Input Capture 1 Prescaler:
	IC1F   CCMR1 = 0x0F << 4 //+ Input Capture 1 Filter.

	CC2S   CCMR1 = 0x03 << 8 //+ Capture/Compare 2 Selection:
	CC2OUT CCMR1 = 0x00 << 8 //  as output,
	CC2TI2 CCMR1 = 0x01 << 8 //  input from TI2,
	CC2TI1 CCMR1 = 0x02 << 8 //  input from TI1,
	CC2TRC CCMR1 = 0x03 << 8 //  input from TRC.

	OC2FE  CCMR1 = 0x01 << 10 //+ Output Compare 2 Fast enable.
	OC2PE  CCMR1 = 0x01 << 11 //+ Output Compare 2 Preload enable.
	OC2M   CCMR1 = 0x07 << 12 //+ Output Compare 2 Mode.
	OC2CE  CCMR1 = 0x01 << 15 //+ Output Compare 2 Clear Enable.
	IC2PSC CCMR1 = 0x03 << 10 //+ Input Capture 2 Prescaler:
	IC2F   CCMR1 = 0x0F << 12 //+ Input Capture 2 Filter.
)

const (
	CC1Sn   = 0
	OC1FEn  = 2
	OC1PEn  = 3
	OC1Mn   = 4
	OC1CEn  = 7
	CC2Sn   = 8
	OC2FEn  = 10
	OC2PEn  = 11
	OC2Mn   = 12
	OC2CEn  = 15
	IC1PSCn = 2
	IC1Fn   = 4
	IC2PSCn = 10
	IC2Fn   = 12
)

const (
	CC3S   CCMR2 = 0x03 << 0 //+ Capture/Compare 3 Selection:
	CC3OUT CCMR2 = 0x00 << 0 //  as output,
	CC3TI3 CCMR2 = 0x01 << 0 //  input from TI3,
	CC3TI4 CCMR2 = 0x02 << 0 //  input from TI4,
	CC3TRC CCMR2 = 0x03 << 0 //  input from TRC.

	OC3FE  CCMR2 = 0x01 << 2 //+ Output Compare 3 Fast enable.
	OC3PE  CCMR2 = 0x01 << 3 //+ Output Compare 3 Preload enable.
	OC3M   CCMR2 = 0x07 << 4 //+ Output Compare 3 Mode.
	OC3CE  CCMR2 = 0x01 << 7 //+ Output Compare 3 Clear Enable.
	IC3PSC CCMR2 = 0x03 << 2 //+ Input Capture 3 Prescaler.
	IC3F   CCMR2 = 0x0F << 4 //+ Input Capture 3 Filter.

	CC4S   CCMR2 = 0x03 << 8 //+ Capture/Compare 4 Selection:
	CC4OUT CCMR2 = 0x00 << 0 //  as output,
	CC4TI4 CCMR2 = 0x01 << 0 //  input from TI4,
	CC4TI3 CCMR2 = 0x02 << 0 //  input from TI3,
	CC4TRC CCMR2 = 0x03 << 0 //  input from TRC.

	OC4FE  CCMR2 = 0x01 << 10 //+ Output Compare 4 Fast enable.
	OC4PE  CCMR2 = 0x01 << 11 //+ Output Compare 4 Preload enable.
	OC4M   CCMR2 = 0x07 << 12 //+ Output Compare 4 Mode.
	OC4CE  CCMR2 = 0x01 << 15 //+ Output Compare 4 Clear Enable.
	IC4PSC CCMR2 = 0x03 << 10 //+ Input Capture 4 Prescaler.
	IC4F   CCMR2 = 0x0F << 12 //+ Input Capture 4 Filter.
)

const (
	CC3Sn   = 0
	OC3FEn  = 2
	OC3PEn  = 3
	OC3Mn   = 4
	OC3CEn  = 7
	CC4Sn   = 8
	OC4FEn  = 10
	OC4PEn  = 11
	OC4Mn   = 12
	OC4CEn  = 15
	IC3PSCn = 2
	IC3Fn   = 4
	IC4PSCn = 10
	IC4Fn   = 12
)

const (
	CC1E  CCER = 0x01 << 0  //+ Capture/Compare 1 output enable.
	CC1P  CCER = 0x01 << 1  //+ Capture/Compare 1 output Polarity.
	CC1NE CCER = 0x01 << 2  //+ Capture/Compare 1 Complementary output enable.
	CC1NP CCER = 0x01 << 3  //+ Capture/Compare 1 Complementary output Polarity.
	CC2E  CCER = 0x01 << 4  //+ Capture/Compare 2 output enable.
	CC2P  CCER = 0x01 << 5  //+ Capture/Compare 2 output Polarity.
	CC2NE CCER = 0x01 << 6  //+ Capture/Compare 2 Complementary output enable.
	CC2NP CCER = 0x01 << 7  //+ Capture/Compare 2 Complementary output Polarity.
	CC3E  CCER = 0x01 << 8  //+ Capture/Compare 3 output enable.
	CC3P  CCER = 0x01 << 9  //+ Capture/Compare 3 output Polarity.
	CC3NE CCER = 0x01 << 10 //+ Capture/Compare 3 Complementary output enable.
	CC3NP CCER = 0x01 << 11 //+ Capture/Compare 3 Complementary output Polarity.
	CC4E  CCER = 0x01 << 12 //+ Capture/Compare 4 output enable.
	CC4P  CCER = 0x01 << 13 //+ Capture/Compare 4 output Polarity.
	CC4NP CCER = 0x01 << 15 //+ Capture/Compare 4 Complementary output Polarity.
)

const (
	CC1En  = 0
	CC1Pn  = 1
	CC1NEn = 2
	CC1NPn = 3
	CC2En  = 4
	CC2Pn  = 5
	CC2NEn = 6
	CC2NPn = 7
	CC3En  = 8
	CC3Pn  = 9
	CC3NEn = 10
	CC3NPn = 11
	CC4En  = 12
	CC4Pn  = 13
	CC4NPn = 15
)
