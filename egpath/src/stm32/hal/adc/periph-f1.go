// +build f10x_ld f10x_ld_vl f10x_md f10x_md_vl f10x_hd f10x_hd_vl f10x_xl f10x_cl

package adc

import (
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

func (p *Periph) enabled() {
	return p.raw.ADON().Load() != 0
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

func (p *Periph) setSmplTime(ch int, st SmplTime) {
	checkCh(ch)
	if ch < 10 {
		n := uint(ch) * 3
		p.raw.SMPR2.StoreBits(adc.SMPR2_Bits(7)<<n, adc.SMPR2_Bits(st)<<n)
	} else {
		n := uint(ch-10) * 3
		p.raw.SMPR1.StoreBits(adc.SMPR1_Bits(7)<<n, adc.SMPR1_Bits(st)<<n)
	}
}
