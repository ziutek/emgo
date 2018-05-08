// This is example of nRF24L01 based remote servo.
// See ../../minidev/nrfrc for transmiter application that can controll it.
package main

import (
	"delay"
	"fmt"
	"mmio"
	"rtos"

	"nrf24"

	"stm32/nrfdci"

	"stm32/hal/dma"
	"stm32/hal/exti"
	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/spi"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"

	"stm32/hal/raw/rcc"
	"stm32/hal/raw/tim"
)

var (
	dci *nrfdci.DCI
	pwm *mmio.U32
	led gpio.Pin
)

func init() {
	system.Setup96(8)
	systick.Setup(2e6)
	start := rtos.Nanosec()

	// GPIO

	gpio.A.EnableClock(true)
	led = gpio.A.Pin(1)
	spiport, sck, miso, mosi := gpio.A, gpio.Pin5, gpio.Pin6, gpio.Pin7
	pwmport, pwmpin := gpio.A, gpio.Pin0

	gpio.B.EnableClock(true)
	csn := gpio.B.Pin(6)
	ctrport, irqn, ce := gpio.B, gpio.Pin8, gpio.Pin9

	// LED

	led.Setup(&gpio.Config{Mode: gpio.Out, Speed: gpio.Low})

	// PWM

	pwmport.Setup(pwmpin, &gpio.Config{Mode: gpio.Alt, Speed: gpio.Low})
	pwmport.SetAltFunc(pwmpin, gpio.TIM2)
	rcc.RCC.TIM2EN().Set()
	t := tim.TIM2
	const (
		pwmmode   = 6    // Mode 1
		pwmperiod = 20   // ms
		pwmmax    = 20e3 // (for 20 ms period it gives 1 µs resolution).
	)
	pclk := system.APB1.Clock()
	if pclk != system.AHB.Clock() {
		pclk *= 2
	}
	t.PSC.Store(tim.PSC(pclk/1000*pwmperiod/pwmmax - 1))
	t.ARR.Store(pwmmax - 1)
	t.CCMR1.Store(pwmmode<<tim.OC1Mn | tim.OC1PE)
	t.CCER.Store(tim.CC1E)
	t.EGR.Store(tim.UG)
	t.CR1.Store(tim.ARPE | tim.CEN)
	pwm = &t.CCR1.U32

	// nRF24 SPI.

	spiport.Setup(sck|mosi, &gpio.Config{Mode: gpio.Alt, Speed: gpio.High})
	spiport.Setup(miso, &gpio.Config{Mode: gpio.AltIn})
	spiport.SetAltFunc(sck|miso|mosi, gpio.SPI1)
	d := dma.DMA2
	d.EnableClock(true)
	spid := spi.NewDriver(spi.SPI1, d.Channel(3, 3), d.Channel(2, 3))
	spid.Periph().EnableClock(true)
	rtos.IRQ(irq.SPI1).Enable()
	rtos.IRQ(irq.DMA2_Stream2).Enable()
	rtos.IRQ(irq.DMA2_Stream3).Enable()

	// nRF24 control lines.

	csn.Setup(&gpio.Config{Mode: gpio.Out, Speed: gpio.High})
	ctrport.Setup(ce, &gpio.Config{Mode: gpio.Alt, Speed: gpio.High})
	ctrport.SetAltFunc(ce, gpio.TIM4)
	rcc.RCC.TIM4EN().Set()
	ctrport.Setup(irqn, &gpio.Config{Mode: gpio.In, Pull: gpio.PullUp})
	irqline := exti.Lines(irqn)
	irqline.Connect(ctrport)
	rtos.IRQ(irq.EXTI9_5).Enable()

	dci = nrfdci.NewDCI(
		spid, csn, system.APB1.Clock(), tim.TIM4, 4, irqline,
	)

	// nRF24 requires wait at least 100 ms from start before use it.
	rtos.SleepUntil(start + 100e6)
}

func checkErr(err error) {
	if err == nil {
		return
	}
	fmt.Printf("Error: %v\n", err)
	for {
	}
}

func main() {
	// For SG92R servo (PWM: 3.3V, 20 ms).
	const (
		min    = 600  // µs
		max    = 2400 // µs
		center = (min + max) / 2
	)

	buf := make([]byte, 2)

	nrf := nrf24.NewRadio(dci)
	nrf.Set_RF_CH(50)
	nrf.Set_RF_SETUP(nrf24.RF_DR_LOW)
	nrf.Set_EN_AA(0)
	nrf.Set_EN_RXADDR(nrf24.P0)
	nrf.Set_SETUP_AW(3)
	nrf.Set_RX_PW(0, len(buf))
	nrf.Set_CONFIG(nrf24.PWR_UP | nrf24.EN_CRC | nrf24.CRCO | nrf24.PRIM_RX)

	// Wait for transition from Power Down to Standby I.
	delay.Millisec(5)

	nrf.FLUSH_RX()
	dci.IRQF().Reset(0)
	nrf.ClearIRQ(nrf24.RX_DR)
	dci.SetCE(1)

	var fs nrf24.FIFO_STATUS
	for {
		if fs&nrf24.RX_EMPTY != 0 {
			led.Clear()
			dci.IRQF().Wait(1, 0)
		}
		led.Set()
		nrf.R_RX_PAYLOAD(buf)
		dci.IRQF().Reset(0)
		nrf.ClearIRQ(nrf24.RX_DR)

		x, y := int(int8(buf[0])), int(int8(buf[1]))
		_ = y
		switch {
		case x < -64:
			x = -64
		case x > 64:
			x = 64
		}
		pwm.Store(uint32(center + x*(max-min)/128))
		fs, _ = nrf.FIFO_STATUS()
	}
}

func exti9_5ISR() {
	p := exti.Pending() & (exti.L9 | exti.L8 | exti.L7 | exti.L6 | exti.L5)
	p.ClearPending()
	if p&dci.IRQL() != 0 {
		dci.ISR()
	}
}

func nrfSPIISR() {
	dci.SPI().ISR()
}

func nrfRxDMAISR() {
	dci.SPI().DMAISR(dci.SPI().RxDMA())
}

func nrfTxDMAISR() {
	dci.SPI().DMAISR(dci.SPI().TxDMA())
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.EXTI9_5:      exti9_5ISR,
	irq.SPI1:         nrfSPIISR,
	irq.DMA2_Stream2: nrfRxDMAISR,
	irq.DMA2_Stream3: nrfTxDMAISR,
}
