// +build f303xe

package adc

import (
	"bits"
	"delay"
	"rtos"
	"unsafe"

	"stm32/hal/raw/adc"
	"stm32/hal/raw/mmap"
	"stm32/hal/raw/rcc"
)

//emgo:const
var (
	ADC1 = (*Periph)(unsafe.Pointer(mmap.ADC1_BASE))
	ADC2 = (*Periph)(unsafe.Pointer(mmap.ADC2_BASE))
	ADC3 = (*Periph)(unsafe.Pointer(mmap.ADC3_BASE))
	ADC4 = (*Periph)(unsafe.Pointer(mmap.ADC4_BASE))
)

func (p *Periph) common() *adc.ADC_Common_Periph {
	addr := uintptr(unsafe.Pointer(p))&^0x100 + 0x300
	return (*adc.ADC_Common_Periph)(unsafe.Pointer(addr))
}

func (p *Periph) enableClock(_ bool) {
	switch p.common() {
	case adc.ADC1_2:
		rcc.RCC.ADC12EN().Set()
	case adc.ADC3_4:
		rcc.RCC.ADC34EN().Set()
	}
}

func (p *Periph) disableClock() {
	switch p.common() {
	case adc.ADC1_2:
		rcc.RCC.ADC12EN().Clear()
	case adc.ADC3_4:
		rcc.RCC.ADC34EN().Clear()
	}
}

const advregen = 1 << adc.ADVREGENn

// EnableVoltage enables p's internal voltage regulator.
func (p *Periph) EnableVoltage() {
	raw := &p.raw
	raw.CR.Store(0)
	raw.CR.Store(advregen)
	delay.Millisec(1) // TODO:  Do not wait so long (setup time <= 10 µs).
}

// DisableVoltage disables p's internal voltage regulator.
func (p *Periph) DisableVoltage() {
	raw := &p.raw
	raw.CR.Store(0)
	raw.CR.Store(2 << adc.ADVREGENn)
}

func (p *Periph) calibrate() {
	raw := &p.raw
	raw.CR.Store(adc.ADCAL | advregen)
	for raw.ADCAL().Load() != 0 {
		rtos.SchedYield()
	}
}

func (p *Periph) enable() {
	p.raw.CR.Store(adc.ADEN | advregen)
}

func (p *Periph) enabled() bool {
	return p.raw.ADEN().Load() != 0
}

func (p *Periph) disable() {
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

//emgo:const
var halfCycles = [8]uint16{
	1.5 * 2,
	2.5 * 2,
	4.5 * 2,
	7.5 * 2,
	19.5 * 2,
	61.5 * 2,
	181.5 * 2,
	601.5 * 2,
}

func checkCh(ch int) {
	if ch < 1 || ch > 18 {
		panicCN()
	}
}

func (p *Periph) setSmplTime(ch int, st SmplTime) {
	checkCh(ch)
	if ch < 10 {
		n := uint(ch) * 3
		p.raw.SMPR1.StoreBits(adc.SMPR1_Bits(7)<<n, adc.SMPR1_Bits(st)<<n)
	} else {
		n := uint(ch-10) * 3
		p.raw.SMPR2.StoreBits(adc.SMPR2_Bits(7)<<n, adc.SMPR2_Bits(st)<<n)
	}
}

func (p *Periph) setRegularSeq(ch []int) {
	sqr1 := adc.SQR1_Bits(len(ch)-1) << adc.Ln
	sq := ch
	ch = nil
	if len(sq) > 4 {
		ch = sq[4:]
		sq = sq[:4]
	}
	for i, c := range sq {
		checkCh(c)
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
		checkCh(c)
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
		checkCh(c)
		sqr3 |= adc.SQR3_Bits(c) << (uint(i) * 6)
	}
	raw.SQR3.Store(sqr3)
	if len(ch) > 2 {
		ch = ch[:2]
	}
	var sqr4 adc.SQR4_Bits
	for i, c := range ch {
		checkCh(c)
		sqr4 |= adc.SQR4_Bits(c) << (uint(i) * 6)
	}
	raw.SQR4.Store(sqr4)
}

// External trigger for ADC1 and ADC2 regular channels.
const (
	ADC12_TIM1_CC1    TrigSrc = 0
	ADC12_TIM1_CC2    TrigSrc = 1
	ADC12_TIM1_CC3    TrigSrc = 2
	ADC12_TIM20_TRGO  TrigSrc = 2
	ADC12_TIM2_CC2    TrigSrc = 3
	ADC12_TIM20_TRGO2 TrigSrc = 3
	ADC12_TIM3_TRGO   TrigSrc = 4
	ADC12_TIM4_CC4    TrigSrc = 5
	ADC12_TIM20_CC1   TrigSrc = 5
	ADC12_EXTI11      TrigSrc = 6
	ADC12_TIM8_TRGO   TrigSrc = 7
	ADC12_TIM8_TRGO2  TrigSrc = 8
	ADC12_TIM1_TRGO   TrigSrc = 9
	ADC12_TIM1_TRGO2  TrigSrc = 10
	ADC12_TIM2_TRGO   TrigSrc = 11
	ADC12_TIM4_TRGO   TrigSrc = 12
	ADC12_TIM6_TRGO   TrigSrc = 13
	ADC12_TIM20_CC2   TrigSrc = 13
	ADC12_TIM15_TRGO  TrigSrc = 14
	ADC12_TIM3_CC4    TrigSrc = 15
	ADC12_TIM20_CC3   TrigSrc = 15
)

// External trigger for ADC2 and ADC3 regular channels.
const (
	ADC23_TIM3_CC1    TrigSrc = 0
	ADC23_TIM2_CC3    TrigSrc = 1
	ADC23_TIM1_CC3    TrigSrc = 2
	ADC23_TIM8_CC1    TrigSrc = 3
	ADC23_TIM8_TRGO   TrigSrc = 4
	ADC23_EXTI2       TrigSrc = 5
	ADC23_TIM20_TRGO  TrigSrc = 5
	ADC23_TIM4_CC1    TrigSrc = 6
	ADC23_TIM20_TRGO2 TrigSrc = 6
	ADC23_TIM2_TRGO   TrigSrc = 7
	ADC23_TIM8_TRGO2  TrigSrc = 8
	ADC23_TIM1_TRGO   TrigSrc = 9
	ADC23_TIM1_TRGO2  TrigSrc = 10
	ADC23_TIM3_TRGO   TrigSrc = 11
	ADC23_TIM4_TRGO   TrigSrc = 12
	ADC23_TIM7_TRGO   TrigSrc = 13
	ADC23_TIM15_TRGO  TrigSrc = 14
	ADC23_TIM2_CC1    TrigSrc = 15
	ADC23_TIM20_CC1   TrigSrc = 15
)

func (p *Periph) setTrigSrc(src TrigSrc) {
	p.raw.EXTSEL().Store(adc.CFGR_Bits(src) << adc.EXTSELn)
}

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

	evAll = Ready | SmplEnd | ConvEnd | SeqEnd | Overrun | InjConvEnd |
		InjSeqEnd | Watchdog1 | Watchdog2 | Watchdog3 | InjOverflow
)

func (p *Periph) event() Event {
	return Event(p.raw.ISR.Load())
}

func (p *Periph) clear(events Event) {
	p.raw.ISR.Store(adc.ISR_Bits(events))
}

func (p *Periph) enableIRQ(events Event) {
	p.raw.IER.SetBits(adc.IER_Bits(events))
}

func (p *Periph) disableIRQ(events Event) {
	if events == EvAll {
		p.raw.IER.Store(0)
	} else {
		p.raw.IER.ClearBits(adc.IER_Bits(events))
	}
}

func (p *Periph) enableDMA(circural bool) {
	p.raw.CFGR.StoreBits(
		adc.DMAEN|adc.DMACFG,
		adc.DMAEN|adc.CFGR_Bits(bits.One(circural)<<adc.DMACFGn),
	)
}

func (p *Periph) disableDMA() {
	p.raw.DMAEN().Clear()
}

func (p *Periph) start() {
	p.raw.CR.Store(adc.ADSTART | advregen)
}
