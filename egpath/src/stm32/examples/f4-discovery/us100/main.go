// Connct US-100 Tx and Rx pins respectively to Discovery's PA2 (USART2_TX),
// PA3 (USART2_RX) pins (Tx-Tx, Rx-Rx).
package main

import (
	"delay"
	"fmt"
	"io"
	"mmio"
	"rtos"

	"stm32/hal/dma"
	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
	"stm32/hal/usart"

	"stm32/hal/raw/rcc"
	"stm32/hal/raw/tim"
)

const ledmax = 15

var (
	tts                      *usart.Driver
	dmarxbuf                 [8]byte
	green, orange, red, blue *mmio.U32
)

func init() {
	system.Setup168(8)
	systick.Setup(2e6)

	// GPIO

	gpio.A.EnableClock(true)
	port, tx, rx := gpio.A, gpio.Pin2, gpio.Pin3
	gpio.D.EnableClock(false)
	ledport, ledpins := gpio.D, gpio.Pin12|gpio.Pin13|gpio.Pin14|gpio.Pin15

	// LEDS

	ledport.Setup(ledpins, &gpio.Config{Mode: gpio.Alt, Speed: gpio.Low})
	ledport.SetAltFunc(ledpins, gpio.TIM4)
	rcc.RCC.TIM4EN().Set()
	t := tim.TIM4
	const (
		pwmmode = 6    // Mode 1
		pwmfreq = 1000 // Hz
		pwmmax  = 1 << ledmax
	)
	pclk := system.APB1.Clock()
	if pclk != system.AHB.Clock() {
		pclk *= 2
	}
	t.PSC.U16.Store(uint16(pclk/(pwmfreq*pwmmax) - 1))
	t.ARR.Store(pwmmax - 1)
	t.CCMR1.Store(
		pwmmode<<tim.OC1Mn | pwmmode<<tim.OC2Mn | tim.OC1PE | tim.OC2PE,
	)
	t.CCMR2.Store(
		pwmmode<<tim.OC3Mn | pwmmode<<tim.OC4Mn | tim.OC3PE | tim.OC4PE,
	)
	t.CCER.Store(tim.CC1E | tim.CC2E | tim.CC3E | tim.CC4E)
	t.EGR.Store(tim.UG)
	t.CR1.Store(tim.ARPE | tim.CEN)
	green = &t.CCR1.U32
	orange = &t.CCR2.U32
	red = &t.CCR3.U32
	blue = &t.CCR4.U32

	// USART

	port.Setup(tx, &gpio.Config{Mode: gpio.Alt})
	port.Setup(rx, &gpio.Config{Mode: gpio.AltIn, Pull: gpio.PullUp})
	port.SetAltFunc(tx|rx, gpio.USART2)
	d := dma.DMA1
	d.EnableClock(true) // DMA clock must remain enabled in sleep mode.
	tts = usart.NewDriver(
		usart.USART2, d.Channel(6, 4), d.Channel(5, 4), dmarxbuf[:],
	)
	tts.Periph().EnableClock(true)
	tts.Periph().SetBaudRate(9600)
	tts.Periph().Enable()
	tts.EnableRx()
	tts.EnableTx()
	rtos.IRQ(irq.USART2).Enable()
	rtos.IRQ(irq.DMA1_Stream5).Enable()
	rtos.IRQ(irq.DMA1_Stream6).Enable()
}

func checkErr(err error) {
	if err != nil {
		fmt.Printf("\nError: %v\n", err)
	}
}

func split(v, max uint32) (x, rest uint32) {
	if v <= max {
		return v, 0
	}
	return max, v - max
}

func main() {
	buf := make([]byte, 2)
	delay.Millisec(200) // Wait for US-100 startup.
	for {
		checkErr(tts.WriteByte(0x55))
		_, err := io.ReadFull(tts, buf)
		checkErr(err)

		x := int(buf[0])<<8 + int(buf[1])

		var v uint32
		if x > 30 {
			v = uint32(x-30) / 8
		}
		if v < 4*ledmax {
			v = 4*ledmax - v
		} else {
			v = 0
		}
		g, v := split(v, ledmax)
		o, v := split(v, ledmax)
		r, v := split(v, ledmax)
		b, _ := split(v, ledmax)
		green.Store(1 << g)
		orange.Store(1 << o)
		red.Store(1 << r)
		blue.Store(1 << b)

		fmt.Printf("%d mm\n", x)
	}
}

func ttsISR() {
	tts.ISR()
}

func ttsRxDMAISR() {
	tts.RxDMAISR()
}

func ttsTxDMAISR() {
	tts.TxDMAISR()
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.USART2:       ttsISR,
	irq.DMA1_Stream5: ttsRxDMAISR,
	irq.DMA1_Stream6: ttsTxDMAISR,
}
