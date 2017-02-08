package main

import (
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
