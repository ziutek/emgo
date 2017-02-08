// +build f303xe

package adc

import (
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
