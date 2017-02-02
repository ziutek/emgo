package main

import (
	"delay"
	"rtos"
	"sync/fence"
	"unsafe"

	"stm32/hal/dma"

	"stm32/hal/raw/adc"
	"stm32/hal/raw/tim"
)

type ADCDriver struct {
	ADC  *adc.ADC_Periph
	DMA  *dma.Channel
	TIM  *tim.TIM_Periph
	done rtos.EventFlag
}

func NewADCDriver(adc *adc.ADC_Periph, dma *dma.Channel, tim *tim.TIM_Periph) *ADCDriver {
	d := new(ADCDriver)
	d.ADC = adc
	d.DMA = dma
	d.TIM = tim
	return d
}

const advregen = 1 << adc.ADVREGENn

func (d *ADCDriver) EnableVReg() {
	d.ADC.CR.Store(0)
	d.ADC.CR.Store(advregen)
	delay.Millisec(1)
}

func (d *ADCDriver) Callibrate(log2ahbdiv int) {
	ckmode := adc.CCR_Bits(log2ahbdiv+1) << adc.CKMODEn
	a := d.ADC
	if a == adc.ADC1 || a == adc.ADC2 {
		adc.ADC1_2.CKMODE().Store(ckmode)
	} else {
		adc.ADC3_4.CKMODE().Store(ckmode)
	}
	a.CR.Store(adc.ADCAL | advregen)
	for a.ADCAL().Load() != 0 {
		rtos.SchedYield()
	}
	// ADEN can be set not sooner than 4 ADC clock cycles after ADCAL == 0.
	delay.Loop(5 << uint(log2ahbdiv))
}

func (d *ADCDriver) Enable() {
	a := d.ADC
	a.CR.Store(adc.ADEN | advregen)
	for a.ADRDY().Load() == 0 {
		rtos.SchedYield()
	}

}

type ADCRes byte

const (
	Res12 ADCRes = 0
	Res10 ADCRes = 1
	Res8  ADCRes = 2
	Res6  ADCRes = 3
)

func (d *ADCDriver) SetResolution(res ADCRes) {
	d.ADC.RES().Store(adc.CFGR_Bits(res) << adc.RESn)
}

func (d *ADCDriver) SetRegularSeq(ch ...int) {
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
	d.ADC.SQR1.Store(sqr1)
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
	d.ADC.SQR2.Store(sqr2)
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
	d.ADC.SQR3.Store(sqr3)
	if len(ch) > 2 {
		ch = ch[:2]
	}
	var sqr4 adc.SQR4_Bits
	for i, c := range ch {
		sqr4 |= adc.SQR4_Bits(c) << (uint(i) * 6)
	}
	d.ADC.SQR4.Store(sqr4)
}

type ADCExtTrigSrc byte

const (
	ADC12_TIM1_CC1    ADCExtTrigSrc = 0
	ADC12_TIM1_CC2    ADCExtTrigSrc = 1
	ADC12_TIM1_CC3    ADCExtTrigSrc = 2
	ADC12_TIM20_TRGO  ADCExtTrigSrc = 2
	ADC12_TIM2_CC2    ADCExtTrigSrc = 3
	ADC12_TIM20_TRGO2 ADCExtTrigSrc = 3
	ADC12_TIM3_TRGO   ADCExtTrigSrc = 4
	ADC12_TIM4_CC4    ADCExtTrigSrc = 5
	ADC12_TIM20_CC1   ADCExtTrigSrc = 5
	ADC12_EXTI11      ADCExtTrigSrc = 6
	ADC12_TIM8_TRGO   ADCExtTrigSrc = 7
	ADC12_TIM8_TRGO2  ADCExtTrigSrc = 8
	ADC12_TIM1_TRGO   ADCExtTrigSrc = 9
	ADC12_TIM1_TRGO2  ADCExtTrigSrc = 10
	ADC12_TIM2_TRGO   ADCExtTrigSrc = 11
	ADC12_TIM4_TRGO   ADCExtTrigSrc = 12
	ADC12_TIM6_TRGO   ADCExtTrigSrc = 13
	ADC12_TIM20_CC2   ADCExtTrigSrc = 13
	ADC12_TIM15_TRGO  ADCExtTrigSrc = 14
	ADC12_TIM3_CC4    ADCExtTrigSrc = 15
	ADC12_TIM20_CC3   ADCExtTrigSrc = 15
)

func (d *ADCDriver) SetExtTrigSrc(src ADCExtTrigSrc) {
	d.ADC.EXTSEL().Store(adc.CFGR_Bits(src) << adc.EXTSELn)
}

type ADCExtTrigEdge byte

const (
	EdgeNone    ADCExtTrigEdge = 0
	EdgeRising  ADCExtTrigEdge = 1
	EdgeFalling ADCExtTrigEdge = 2
)

func (d *ADCDriver) SetExtTrigEdge(edge ADCExtTrigEdge) {
	d.ADC.EXTEN().Store(adc.CFGR_Bits(edge) << adc.EXTENn)
}

// div1 > 0, div2 > 1.
func (d *ADCDriver) SetTimer(div1, div2 int) {
	t := d.TIM
	t.CR2.Store(2 << tim.MMSn)
	t.ARR.Store(tim.ARR_Bits(div2 - 1))
	t.PSC.Store(tim.PSC_Bits(div1 - 1))
}

func (d *ADCDriver) StartADC() {
	a := d.ADC
	a.DMAEN().Set()
	d.TIM.CR1.Store(tim.ARPE | tim.CEN)
}

func (d *ADCDriver) setupDMA(wordSize uintptr) {
	ch := d.DMA
	ch.Setup(dma.PTM | dma.IncM | dma.FIFO_1_4)
	ch.SetWordSize(wordSize, wordSize)
	ch.SetAddrP(unsafe.Pointer(d.ADC.DR.U32.Addr()))
}

func (d *ADCDriver) readDMA(addr uintptr, n int) {
	d.done.Reset(0)
	ch := d.DMA
	ch.SetAddrM(unsafe.Pointer(addr))
	ch.SetLen(n)
	ch.Clear(dma.EvAll, dma.ErrAll)
	ch.EnableIRQ(dma.Complete, dma.ErrAll&^dma.ErrFIFO)
	fence.W() // This orders writes to normal and I/O memory.
	ch.Enable()
	d.ADC.CR.Store(adc.ADSTART | advregen)
	d.done.Wait(1, 0)
	ch.Disable() // Required by F1
}

func (d *ADCDriver) DMAISR() {
	d.DMA.DisableIRQ(dma.EvAll, dma.ErrAll)
	d.done.Signal(1)
}

func (d *ADCDriver) Read(buf []byte) (int, error) {
	if len(buf) == 0 {
		return 0, nil
	}
	d.setupDMA(1)
	d.readDMA(uintptr(unsafe.Pointer(&buf[0])), len(buf))
	return len(buf), nil
}
