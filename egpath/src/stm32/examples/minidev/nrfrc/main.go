// This is example of nRF24L01 based RC transmiter.
// See ../../nucleo-f411re/rcservo that can be used on receiver side.
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
	"stm32/hal/system/timer/rtcst"

	"stm32/hal/raw/rcc"
	"stm32/hal/raw/tim"
)

var (
	led gpio.Pin
	dci *nrfdci.DCI
	enc *mmio.U16
)

func init() {
	system.Setup(8, 1, 72/8)
	rtcst.Setup(32768)
	start := rtos.Nanosec()

	gpio.A.EnableClock(true)
	csn := gpio.A.Pin(4)
	spiport, sck, miso, mosi := gpio.A, gpio.Pin5, gpio.Pin6, gpio.Pin7

	gpio.B.EnableClock(true)
	ctrport, ce, irqn := gpio.B, gpio.Pin0, gpio.Pin1
	encport, encpins := gpio.B, gpio.Pin6|gpio.Pin7

	gpio.C.EnableClock(false)
	led = gpio.C.Pin(13)

	// LED

	led.Setup(&gpio.Config{Mode: gpio.Out, Driver: gpio.OpenDrain, Speed: gpio.Low})
	led.Set()

	// Encoder

	encport.Setup(encpins, &gpio.Config{Mode: gpio.AltIn, Pull: gpio.PullUp})
	rcc.RCC.TIM4EN().Set()
	t := tim.TIM4
	t.SMCR.Store(1 << tim.SMSn)
	t.CCMR1.Store(tim.CC1S | tim.CC2S | 0xf<<tim.IC1Fn | 0xf<<tim.IC2Fn)
	t.CCER.Store(tim.CC1P | tim.CC2P)
	t.CR1.Store(2<<tim.CKDn | tim.CEN)
	enc = &t.CNT.U16

	// nRF24 SPI

	spiport.Setup(sck|mosi, &gpio.Config{Mode: gpio.Alt, Speed: gpio.High})
	spiport.Setup(miso, &gpio.Config{Mode: gpio.AltIn})
	csn.Setup(&gpio.Config{Mode: gpio.Out, Speed: gpio.High})
	d := dma.DMA1
	d.EnableClock(true)
	spid := spi.NewDriver(spi.SPI1, d.Channel(2, 0), d.Channel(3, 0))
	spid.P.EnableClock(true)
	rtos.IRQ(irq.SPI1).Enable()
	rtos.IRQ(irq.DMA1_Channel2).Enable()
	rtos.IRQ(irq.DMA1_Channel3).Enable()

	// nRF24 control lines.

	ctrport.Setup(ce, &gpio.Config{Mode: gpio.Alt, Speed: gpio.High})
	rcc.RCC.TIM3EN().Set()
	ctrport.Setup(irqn, &gpio.Config{Mode: gpio.In, Pull: gpio.PullUp})
	irqline := exti.Lines(irqn)
	irqline.Connect(ctrport)
	rtos.IRQ(irq.EXTI1).Enable()

	dci = nrfdci.NewDCI(spid, csn, system.APB1.Clock(), tim.TIM3, 3, irqline)
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
	nrf := nrf24.NewRadio(dci)
	nrf.Set_RF_CH(50)
	nrf.Set_RF_SETUP(nrf24.RF_DR_LOW | nrf24.RF_PWR(0))
	nrf.Set_EN_AA(0)
	nrf.Set_SETUP_AW(3)
	nrf.FLUSH_TX()
	nrf.Set_CONFIG(nrf24.PWR_UP | nrf24.EN_CRC | nrf24.CRCO)

	// Wait for transition from Power Down to Standby I.
	delay.Millisec(5)

	for {
		cnt := int16(enc.Load())
		switch {
		case cnt < -64:
			cnt = -64
			enc.Store(uint16(cnt))
		case cnt > 64:
			cnt = 64
			enc.Store(uint16(cnt))
		}
		nrf.W_TX_PAYLOAD([]byte{byte(cnt), 0})
		nrf.ClearIRQ(nrf24.TX_DS)
		dci.IRQF().Reset(0)
		dci.SetCE(2)
		dci.IRQF().Wait(1, 0)
		delay.Millisec(40)
	}
}

func exti1ISR() {
	exti.L1.ClearPending()
	dci.ISR()
}

func nrfSPIISR() {
	dci.SPI().ISR()
}

func nrfRxDMAISR() {
	dci.SPI().DMAISR(dci.SPI().RxDMA)
}

func nrfTxDMAISR() {
	dci.SPI().DMAISR(dci.SPI().TxDMA)
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTCAlarm: rtcst.ISR,

	irq.EXTI1:         exti1ISR,
	irq.SPI1:          nrfSPIISR,
	irq.DMA1_Channel2: nrfRxDMAISR,
	irq.DMA1_Channel3: nrfTxDMAISR,
}
