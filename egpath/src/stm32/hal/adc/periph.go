package adc

import (
	"bits"
	"rtos"
	"unsafe"

	"stm32/hal/internal"
	"stm32/hal/system"

	"stm32/hal/raw/adc"
)

type Periph struct {
	raw adc.ADC_Periph
}

// Bus returns a bus to which p is connected to.
func (p *Periph) Bus() system.Bus {
	return internal.Bus(unsafe.Pointer(p))
}

// EnableClock enables clock for p. Lp determines whether the clock remains on
// in low power (sleep) mode. In some MCUs clock cannot be enabled for only one
// ADC, then EnableClock can affects a pair of peripherals (eg.
// ADC1.EnableClock() can affect ADC1 and ADC2 simultaneously).
func (p *Periph) EnableClock(lp bool) {
	p.enableClock(lp)
}

// DisableClock disables clock for p. In some MCUs clock cannot be disabled for
// only one ADC, then DisableClock can affects a pair of peripherals (eg.
// ADC1.DisableClock() can affect ADC1 and ADC2 simultaneously).
func (p *Periph) DisableClock() {
	p.disableClock()
}

const advregen = 1 << adc.ADVREGENn

// EnableVoltage enables p's internal voltage regulator. Must wait for setup
// time (typically 10 µs) before calibrate or enable p.
func (p *Periph) EnableVoltage() {
	raw := &p.raw
	raw.CR.Store(0)
	raw.CR.Store(advregen)
}

// DisableVoltage disables p's internal voltage regulator.
func (p *Periph) DisableVoltage() {
	raw := &p.raw
	raw.CR.Store(0)
	raw.CR.Store(2 << adc.ADVREGENn)
}

type ClockMode byte

const (
	CKADC ClockMode = 0 // Asynchronous clock
	HCLK1 ClockMode = 1 // AHB clock / 1.
	HCLK2 ClockMode = 2 // AHB clock / 2.
	HCLK4 ClockMode = 3 // AHB clock / 4.
)

func (p *Periph) SetClockMode(ckmode ClockMode) {
	p.common().CKMODE().Store(adc.CCR_Bits(ckmode) << adc.CKMODEn)
}

// Calibrate calibrates p. Wait at least 5 ADC clock cycles after calibration
// before enable p.
func (p *Periph) Callibrate() {
	raw := &p.raw
	raw.CR.Store(adc.ADCAL | advregen)
	for raw.ADCAL().Load() != 0 {
		rtos.SchedYield()
	}
}

//
func (p *Periph) Enable() {
	raw := &p.raw
	raw.CR.Store(adc.ADEN | advregen)
}

// Enabled reports whether p is enabled.
func (p *Periph) Enabled() bool {
	return p.raw.ADEN().Load() != 0
}

// Disable disables p. Use Enabled to check that p is effectively disabled.
func (p *Periph) Disable() {
	p.raw.CR.Store(adc.ADDIS | advregen)
}

type Resolution byte

const (
	Res12 Resolution = 0
	Res10 Resolution = 1
	Res8  Resolution = 2
	Res6  Resolution = 3
)

func (p *Periph) SetResolution(res Resolution) {
	p.raw.RES().Store(adc.CFGR_Bits(res) << adc.RESn)
}

func (p *Periph) SetRegularSeq(ch ...int) {
	sqr1 := adc.SQR1_Bits(len(ch)-1) << adc.Ln
	sq := ch
	ch = nil
	if len(sq) > 4 {
		ch = sq[4:]
		sq = sq[:4]
	}
	for i, c := range sq {
		sqr1 |= adc.SQR1_Bits(c) << (uint(i+1) * 6)
	}
	raw := &p.raw
	raw.SQR1.Store(sqr1)
	sq = ch
	ch = nil
	if len(sq) > 5 {
		ch = sq[5:]
		sq = sq[:5]
	}
	var sqr2 adc.SQR2_Bits
	for i, c := range sq {
		sqr2 |= adc.SQR2_Bits(c) << (uint(i) * 6)
	}
	raw.SQR2.Store(sqr2)
	sq = ch
	ch = nil
	if len(sq) > 5 {
		ch = sq[5:]
		sq = sq[:5]
	}
	var sqr3 adc.SQR3_Bits
	for i, c := range sq {
		sqr3 |= adc.SQR3_Bits(c) << (uint(i) * 6)
	}
	raw.SQR3.Store(sqr3)
	if len(ch) > 2 {
		ch = ch[:2]
	}
	var sqr4 adc.SQR4_Bits
	for i, c := range ch {
		sqr4 |= adc.SQR4_Bits(c) << (uint(i) * 6)
	}
	raw.SQR4.Store(sqr4)
}

type ExtTrigSrc byte

const (
	ADC12_TIM1_CC1    ExtTrigSrc = 0
	ADC12_TIM1_CC2    ExtTrigSrc = 1
	ADC12_TIM1_CC3    ExtTrigSrc = 2
	ADC12_TIM20_TRGO  ExtTrigSrc = 2
	ADC12_TIM2_CC2    ExtTrigSrc = 3
	ADC12_TIM20_TRGO2 ExtTrigSrc = 3
	ADC12_TIM3_TRGO   ExtTrigSrc = 4
	ADC12_TIM4_CC4    ExtTrigSrc = 5
	ADC12_TIM20_CC1   ExtTrigSrc = 5
	ADC12_EXTI11      ExtTrigSrc = 6
	ADC12_TIM8_TRGO   ExtTrigSrc = 7
	ADC12_TIM8_TRGO2  ExtTrigSrc = 8
	ADC12_TIM1_TRGO   ExtTrigSrc = 9
	ADC12_TIM1_TRGO2  ExtTrigSrc = 10
	ADC12_TIM2_TRGO   ExtTrigSrc = 11
	ADC12_TIM4_TRGO   ExtTrigSrc = 12
	ADC12_TIM6_TRGO   ExtTrigSrc = 13
	ADC12_TIM20_CC2   ExtTrigSrc = 13
	ADC12_TIM15_TRGO  ExtTrigSrc = 14
	ADC12_TIM3_CC4    ExtTrigSrc = 15
	ADC12_TIM20_CC3   ExtTrigSrc = 15
)

func (p *Periph) SetExtTrigSrc(src ExtTrigSrc) {
	p.raw.EXTSEL().Store(adc.CFGR_Bits(src) << adc.EXTSELn)
}

type ExtTrigEdge byte

const (
	EdgeNone    ExtTrigEdge = 0
	EdgeRising  ExtTrigEdge = 1
	EdgeFalling ExtTrigEdge = 2
)

func (p *Periph) SetExtTrigEdge(edge ExtTrigEdge) {
	p.raw.EXTEN().Store(adc.CFGR_Bits(edge) << adc.EXTENn)
}

type Event uint16

const (
	Ready       = Event(adc.ADRDY) // ADC ready to accept conversion requests.
	SmplEnd     = Event(adc.EOSMP) // End of sampling phase reached.
	ConvEnd     = Event(adc.EOC)   // Regular channel conversion complete.
	SeqEnd      = Event(adc.EOS)   // Regular conversions sequence complete.
	Overrun     = Event(adc.OVR)   // Overrun occurred.
	InjConvEnd  = Event(adc.JEOC)  // Injected channel conversion complete.
	InjSeqEnd   = Event(adc.JEOS)  // Injected conversions sequence complete.
	Watchdog1   = Event(adc.AWD1)  // Analog watchdog 1 event occurred.
	Watchdog2   = Event(adc.AWD2)  // Analog watchdog 2 event occurred.
	Watchdog3   = Event(adc.AWD3)  // Analog watchdog 3 event occurred.
	InjOverflow = Event(adc.JQOVF) // Inj. context queue overflow occurred.

	EvAll = Ready | SmplEnd | ConvEnd | SeqEnd | Overrun | InjConvEnd |
		InjSeqEnd | Watchdog1 | Watchdog2 | Watchdog3 | InjOverflow
)

func (p *Periph) Event() Event {
	return Event(p.raw.ISR.Load())
}

func (p *Periph) Clear(events Event) {
	p.raw.ISR.Store(adc.ISR_Bits(events))
}

func (p *Periph) EnableIRQ(events Event) {
	p.raw.IER.SetBits(adc.IER_Bits(events))
}

func (p *Periph) DisableIRQ(events Event) {
	if events == EvAll {
		p.raw.IER.Store(0)
	} else {
		p.raw.IER.ClearBits(adc.IER_Bits(events))
	}
}

func (p *Periph) EnableDMA(circural bool) {
	p.raw.CFGR.StoreBits(
		adc.DMAEN|adc.DMACFG,
		adc.DMAEN|adc.CFGR_Bits(bits.One(circural)<<adc.DMACFGn),
	)
}

func (p *Periph) DisableDMA() {
	p.raw.DMAEN().Clear()
}

func (p *Periph) Start() {
	p.raw.CR.Store(adc.ADSTART | advregen)
}
