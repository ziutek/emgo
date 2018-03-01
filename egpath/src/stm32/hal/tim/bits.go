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
	DIR  CR1 = 0x01 << 4 //+ Counter direction.
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
	CCPC  CR2 = 0x01 << 0  //+ Capture/Compare Preloaded Control.
	CCUS  CR2 = 0x01 << 2  //+ Capture/Compare Control Update Selection.
	CCDS  CR2 = 0x01 << 3  //+ Capture/Compare DMA Selection.
	MMS   CR2 = 0x07 << 4  //+ Master Mode Selection (what is used as TRGO):
	MMUG  CR2 = 0x00 << 4  //  EGR.UG bit,
	MMEN  CR2 = 0x01 << 4  //  CNT_EN = CR1.CEN | TRGI,
	MMUPD CR2 = 0x02 << 4  //  update event,
	MMCC1 CR2 = 0x03 << 4  //  CC1IF (enven if not cleared),
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
	SMS      SMCR = 0x07 << 0  //+ Slave mode selection:
	SMDIS    SMCR = 0x01 << 0  //  slave mode disabled,
	SME1     SMCR = 0x01 << 0  //  encoder mode 1: counts on TI2FP2 edge,
	SME2     SMCR = 0x02 << 0  //  encoder mode 2: counts on TI1FP1 edge,
	SME3     SMCR = 0x03 << 0  //  encoder mode 3: counts on TI1FP1, TI2FP2 edge.
	SMRST    SMCR = 0x04 << 0  //  reset mode: TRGI reinitializes the counter,
	SMGTD    SMCR = 0x05 << 0  //  gated mode: TRGI high enables counter clock,
	SMTRG    SMCR = 0x06 << 0  //  trigger mode: TRGI starts the counter,
	SMEXT    SMCR = 0x07 << 0  //  rxternal clock mode: TRGI clocks the counter.
	TS       SMCR = 0x07 << 4  //+ Trigger selection (what is used as TRGI):
	ITR0     SMCR = 0x00 << 4  //  internal trigger 0,
	ITR1     SMCR = 0x01 << 4  //  internal trigger 1,
	ITR2     SMCR = 0x02 << 4  //  internal trigger 2,
	ITR3     SMCR = 0x03 << 4  //  internal trigger 3,
	TI1FED   SMCR = 0x04 << 4  //  TI1 edge detector,
	TI1FP1   SMCR = 0x05 << 4  //  filtered timer input 1,
	TI1FP2   SMCR = 0x06 << 4  //  filtered timer input 2,
	ETRF     SMCR = 0x07 << 4  //  external trigger input.
	MSM      SMCR = 0x01 << 7  //+ Master/slave mode.
	ETF      SMCR = 0x0F << 8  //+ External trigger filter (N valid samples need):
	ETFDIS   SMCR = 0x00 << 8  //  no filter, sampling at f_DTS,
	ETFCN2   SMCR = 0x01 << 8  //  f_SAMPLING=f_CK_INT, N=2,
	ETFCN4   SMCR = 0x02 << 8  //  f_SAMPLING=f_CK_INT, N=4,
	ETFCN8   SMCR = 0x03 << 8  //  f_SAMPLING=f_CK_INT, N=8,
	ETFD2N6  SMCR = 0x04 << 8  //  f_SAMPLING=f_DTS/2, N=6,
	ETFD2N8  SMCR = 0x05 << 8  //  f_SAMPLING=f_DTS/2, N=8,
	ETFD4N6  SMCR = 0x06 << 8  //  f_SAMPLING=f_DTS/4, N=6,
	ETFD4N8  SMCR = 0x07 << 8  //  f_SAMPLING=f_DTS/4, N=8,
	ETFD8N6  SMCR = 0x08 << 8  //  f_SAMPLING=f_DTS/8, N=6,
	ETFD8N6  SMCR = 0x09 << 8  //  f_SAMPLING=f_DTS/8, N=8,
	ETFD16N5 SMCR = 0x0A << 8  //  f_SAMPLING=f_DTS/16, N=5,
	ETFD16N6 SMCR = 0x0B << 8  //  f_SAMPLING=f_DTS/16, N=6,
	ETFD16N8 SMCR = 0x0C << 8  //  f_SAMPLING=f_DTS/16, N=8,
	ETFD32N5 SMCR = 0x0D << 8  //  f_SAMPLING=f_DTS/32, N=5,
	ETFD32N6 SMCR = 0x0E << 8  //  f_SAMPLING=f_DTS/32, N=6,
	ETFD32N8 SMCR = 0x0F << 8  //  f_SAMPLING=f_DTS/32, N=8.
	ETPS     SMCR = 0x03 << 12 //+ External trigger prescaler:
	ETPS1    SMCR = 0x00 << 12 //  prescaller off,
	ETPS2    SMCR = 0x01 << 12 //  ETRP /= 2,
	ETPS4    SMCR = 0x02 << 12 //  ETRP /= 4,
	ETPS8    SMCR = 0x03 << 12 //  ETRP /= 8.
	ECE      SMCR = 0x01 << 14 //+ External clock enable.
	ETP      SMCR = 0x01 << 15 //+ External trigger polarity.
)
