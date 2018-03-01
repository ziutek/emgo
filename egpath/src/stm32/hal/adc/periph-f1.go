// +build f10x_ld f10x_ld_vl f10x_md f10x_md_vl f10x_hd f10x_hd_vl f10x_xl f10x_cl

package adc

import (
	"rtos"
	"unsafe"

	"stm32/hal/raw/adc"
	"stm32/hal/raw/mmap"
)

//emgo:const
var (
	ADC1 = (*Periph)(unsafe.Pointer(mmap.ADC1_BASE))
	ADC2 = (*Periph)(unsafe.Pointer(mmap.ADC2_BASE))
	ADC3 = (*Periph)(unsafe.Pointer(mmap.ADC3_BASE))
)

// External trigger sources for ADC1 and ADC2 regular channels.
const (
	ADC12_TIM1_CC1  TrigSrc = 0
	ADC12_TIM1_CC2  TrigSrc = 1
	ADC12_TIM1_CC3  TrigSrc = 2
	ADC12_TIM2_CC2  TrigSrc = 3
	ADC12_TIM3_TRGO TrigSrc = 4
	ADC12_TIM4_CC4  TrigSrc = 5
	ADC12_EXTI11    TrigSrc = 6
	ADC12_TIM8_TRGO TrigSrc = 6
)

// External trigger sources for ADC3 regular channels.
const (
	ADC3_TIM3_CC1  TrigSrc = 0
	ADC3_TIM2_CC3  TrigSrc = 1
	ADC3_TIM1_CC3  TrigSrc = 2
	ADC3_TIM8_CC1  TrigSrc = 3
	ADC3_TIM8_TRGO TrigSrc = 4
	ADC3_TIM5_CC1  TrigSrc = 5
	ADC3_TIM5_CC3  TrigSrc = 6
)

const (
	Watchdog   = Event(adc.AWD)   // Analog watchdog event occurred.
	ConvEnd    = Event(adc.EOC)   // Regular channel conversion complete.
	InjConvEnd = Event(adc.JEOC)  // Injected channel conversion complete.
	InjStart   = Event(adc.JSTRT) // Injected channel conversion has started
	Start      = Event(adc.STRT)  // Regular channel conversion has started

	evAll = Watchdog | ConvEnd | InjConvEnd | InjStart | Start
)

const errAll Error = 0

func (e *Error) error() string {
	return ""
}

func (p *Periph) calibrate() {
	raw := &p.raw
	raw.CAL().Set()
	for raw.CAL().Load() != 0 {
		rtos.SchedYield()
	}
}

func (p *Periph) enable() {
	raw := &p.raw
	if raw.ADON().Load() == 0 {
		raw.ADON().Set()
	}
}

func (p *Periph) enabled() bool {
	return p.raw.ADON().Load() != 0
}

func (p *Periph) disable() {
	p.raw.ADON().Clear()
}

//emgo:const
var halfCycles = [8]uint16{
	1.5 * 2,
	7.5 * 2,
	13.5 * 2,
	28.5 * 2,
	41.5 * 2,
	55.5 * 2,
	71.5 * 2,
	239.5 * 2,
}

func checkCh(ch int) {
	if ch < 0 || ch > 17 {
		panicCN()
	}
}

func (p *Periph) setSamplTime(ch int, st SamplTime) {
	checkCh(ch)
	if ch < 10 {
		n := uint(ch) * 3
		p.raw.SMPR2.StoreBits(adc.SMPR2(7)<<n, adc.SMPR2(st)<<n)
	} else {
		n := uint(ch-10) * 3
		p.raw.SMPR1.StoreBits(adc.SMPR1(7)<<n, adc.SMPR1(st)<<n)
	}
}

func (p *Periph) setSequence(ch []int) {
	if len(ch) > 17 {
		panicSeq()
	}
	sqr1 := adc.SQR1(len(ch)-1) << adc.Ln
	raw := &p.raw
	var sqr3 adc.SQR3
	sq := ch
	ch = nil
	if len(sq) > 6 {
		ch = sq[6:]
		sq = sq[:6]
	}
	for i, c := range sq {
		checkCh(c)
		sqr3 |= adc.SQR3(c) << (uint(i) * 5)
	}
	raw.SQR3.Store(sqr3)
	var sqr2 adc.SQR2
	sq = ch
	ch = nil
	if len(sq) > 6 {
		ch = sq[6:]
		sq = sq[:6]
	}
	for i, c := range sq {
		checkCh(c)
		sqr2 |= adc.SQR2(c) << (uint(i) * 5)
	}
	raw.SQR2.Store(sqr2)
	for i, c := range ch {
		checkCh(c)
		sqr1 |= adc.SQR1(c) << (uint(i) * 5)
	}
	raw.SQR1.Store(sqr1)
}

func (p *Periph) setTrigSrc(src TrigSrc) {
	p.raw.EXTSEL().Store(adc.CR2(src) << adc.EXTSELn)
}

func (p *Periph) setTrigEdge(edge TrigEdge) {
	p.raw.EXTTRIG().Store(adc.CR2(edge) << adc.EXTTRIGn)
}

func (p *Periph) status() (Event, Error) {
	return Event(p.raw.SR.Load()), 0
}

func (p *Periph) clear(ev Event, _ Error) {
	p.raw.SR.Store(adc.SR(EvAll &^ ev))
}

func (p *Periph) enableIRQ(ev Event, _ Error) {
	cr1 := ev&ConvEnd<<(adc.EOCIEn-adc.EOCn) |
		ev&InjConvEnd<<(adc.JEOCIEn-adc.JEOCn) |
		ev&Watchdog<<(adc.AWDIEn-adc.AWDn)
	p.raw.CR1.SetBits(adc.CR1(cr1))
}

func (p *Periph) disableIRQ(ev Event, _ Error) {
	cr1 := ev&ConvEnd<<(adc.EOCIEn-adc.EOCn) |
		ev&InjConvEnd<<(adc.JEOCIEn-adc.JEOCn) |
		ev&Watchdog<<(adc.AWDIEn-adc.AWDn)
	p.raw.CR1.ClearBits(adc.CR1(cr1))
}

func (p *Periph) enableDMA(_ bool) {
	p.raw.DMA().Set()
}

func (p *Periph) disableDMA() {
	p.raw.DMA().Clear()
}

func (p *Periph) start() {
	p.raw.ADON().Set()
}
