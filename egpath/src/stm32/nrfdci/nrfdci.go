// Package nrfdci allows to configure set of STM32 peripherals to control
// nRF24L01(+) Data and Control Interface.
package nrfdci

import (
	"rtos"

	"stm32/hal/exti"
	"stm32/hal/gpio"
	"stm32/hal/spi"
	"stm32/hal/system"

	"stm32/hal/raw/tim"
)

// DCI implements nrf24.DCI interface. Additionally, it allows to control CE
// signal and wait for interrupt.
type DCI struct {
	rtos.EventFlag

	spi     *spi.Driver
	cet     *tim.TIM_Periph
	irq     exti.Lines
	csnport *gpio.Port
	csnpin  gpio.Pins // uint16
	ocmn    uint16
}

func NewDCI(spidrv *spi.Driver, csnport *gpio.Port, csnpin gpio.Pins, pclk uint, cet *tim.TIM_Periph, ch int, irqline exti.Lines) *DCI {
	dci := new(DCI)
	dci.spi = spidrv
	spidrv.P.SetConf(
		spi.Master | spi.MSBF | spi.CPOL0 | spi.CPHA0 |
			spidrv.P.BR(10e6) | // 10 MHz max.
			spi.SoftSS | spi.ISSHigh,
	)
	spidrv.P.Enable()

	dci.csnport = csnport
	dci.csnpin = csnpin
	csnport.SetPins(csnpin)

	dci.cet = cet
	if pclk != system.AHB.Clock() {
		pclk *= 2
	}

	// PSC=1 gives shortest possible delay (1/pclk) before CE pulse.
	cet.PSC.U16.Store(1)

	// ARR=(pclk+1e5-1)/1e5 corresponds to the shortest posible pulse but
	// not less than 10 us. CE will be asserted after 1/pclk (CCRn=1) for
	// ARR/pclk seconds.
	cet.ARR.Store(tim.ARR_Bits(uint32(pclk+1e5-1) / 1e5))

	// Reset CNT and transfer PSC, ARR to the corresponding shadow registers.
	cet.EGR.Store(tim.UG)

	var cce tim.CCER_Bits
	switch ch {
	case 1:
		cet.CCR1.Store(1)
		cce = tim.CC1E
		dci.ocmn = tim.OC1Mn
	case 2:
		cet.CCR2.Store(1)
		cce = tim.CC2E
		dci.ocmn = tim.OC2Mn
	case 3:
		cet.CCR3.Store(1)
		cce = tim.CC3E
		dci.ocmn = tim.OC3Mn
	case 4:
		cet.CCR4.Store(1)
		cce = tim.CC4E
		dci.ocmn = tim.OC4Mn
	}
	cet.CCER.Store(cce)

	dci.irq = irqline
	irqline.EnableFallTrig()
	irqline.EnableIRQ()
	return dci
}

func (dci *DCI) SPI() *spi.Driver {
	return dci.spi
}

func (dci *DCI) WriteRead(oi ...[]byte) (n int, err error) {
	dci.csnport.ClearPins(dci.csnpin)
	dci.spi.WriteReadMany(oi...)
	dci.csnport.SetPins(dci.csnpin)
	return n, dci.spi.Err()
}

// SetCE allows to control CE line.. v==0 sets CE low, v==1 sets CE high, v==2
// pulses CE high for 10 µs and leaves it low.
func (dci *DCI) SetCE(v int) error {
	switch v {
	case 0:
		dci.cet.CCMR2.Store(4 << dci.ocmn)
	case 1:
		dci.cet.CCMR2.Store(5 << dci.ocmn)
	case 2:
		dci.cet.CCMR2.Store(7 << dci.ocmn)
		dci.cet.CR1.Store(tim.OPM | tim.CEN)
	}
	return nil
}

func (dci *DCI) IRQ() exti.Lines {
	return dci.irq
}

func (dci *DCI) ISR() {
	dci.EventFlag.Set()
}

func (dci *DCI) Baudrate() uint {
	return dci.spi.P.Baudrate(dci.spi.P.Conf())
}
