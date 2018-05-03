// Simple WS2812 example. Notice, that logical high state on PB15 is 3V, that is
// 0.3V bellow spec if the LED is powered 4.7V (5V pins on Discovery).
package main

import (
	"delay"
	"fmt"
	"math/rand"
	"rtos"

	"led"
	"led/ws281x/wsspi"

	"stm32/hal/dma"
	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/spi"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
)

var ws *spi.Driver

func init() {
	system.Setup168(8)
	systick.Setup(2e6)

	// GPIO

	gpio.B.EnableClock(true)
	mosi := gpio.B.Pin(15)

	// SPI (Use SPI2. SPI1 is used by onboard accelerometer).

	mosi.Setup(&gpio.Config{Mode: gpio.Alt, Speed: gpio.High})
	mosi.SetAltFunc(gpio.SPI2)
	d := dma.DMA1
	d.EnableClock(true)
	ws = spi.NewDriver(spi.SPI2, d.Channel(4, 0), nil)
	ws.Periph().EnableClock(true)
	br := ws.Periph().BR(3200e3)
	ws.Periph().SetConf(spi.Master | br | spi.SoftSS | spi.ISSHigh)
	ws.Periph().Enable()
	rtos.IRQ(irq.SPI2).Enable()
	rtos.IRQ(irq.DMA1_Stream4).Enable()
}

func main() {
	delay.Millisec(250) // For SWO handling in ST-Link.

	p := ws.Periph()
	fmt.Printf("\nSPI speed: %d Hz\n", p.Baudrate(p.Conf()))

	var rnd rand.XorShift64
	rnd.Seed(1)
	rgb := wsspi.GRB
	strip := wsspi.Make(24)
	black := rgb.Pixel(0)
	for {
		c := led.Color(rnd.Uint32()).Scale(127)
		pixel := rgb.Pixel(c)
		for i := range strip {
			strip[i] = pixel
			ws.WriteRead(strip.Bytes(), nil)
			ws.WriteReadByte(0) // STM32 leaves MOSI set to the last bit sent.
			delay.Millisec(40)
		}
		for i := range strip {
			strip[i] = black
			ws.WriteRead(strip.Bytes(), nil)
			ws.WriteReadByte(0) // STM32 leaves MOSI set to the last bit sent.
			delay.Millisec(20)
		}
	}
}

func spiISR() {
	ws.ISR()
}

func spiTxDMAISR() {
	ws.DMAISR(ws.TxDMA())
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.SPI2:         spiISR,
	irq.DMA1_Stream4: spiTxDMAISR,
}
