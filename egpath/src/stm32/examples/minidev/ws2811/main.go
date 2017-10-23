package main

import (
	"delay"
	"fmt"
	"math/rand"
	"rtos"

	"ws281x"

	"stm32/hal/dma"
	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
	"stm32/hal/usart"
)

var tts *usart.Driver

func init() {
	system.SetupPLL(8, 1, 72/8)
	systick.Setup(2e6)

	// GPIO

	gpio.A.EnableClock(true)
	port, tx := gpio.A, gpio.Pin2

	// USART

	port.Setup(tx, &gpio.Config{Mode: gpio.Alt})
	d := dma.DMA1
	d.EnableClock(true) // DMA clock must remain enabled in s
	tts = usart.NewDriver(usart.USART2, nil, d.Channel(7, 0), nil)
	tts.P.EnableClock(true)
	tts.P.SetBaudRate(2250e3)    // 36MHz/16: 444 ns/UARTbit, 1333 ns/WS2811bit.
	tts.P.SetConf(usart.Stop0b5) // F103 has no 7-bit mode: save 0.5 bit only.
	tts.P.Enable()
	tts.EnableTx()
	rtos.IRQ(irq.USART1).Enable()
	rtos.IRQ(irq.DMA1_Channel7).Enable()
}

func main() {
	delay.Millisec(250)

	ledram := ws281x.MakeFBU(50)
	pixel := ws281x.MakeFBU(1)
	var rnd rand.XorShift64
	rnd.Seed(rtos.Nanosec())
	for {
		c := ws281x.Color(rnd.Uint32())
		fmt.Printf("%3d %3d %3d\n", c.Red(), c.Green(), c.Blue())
		pixel.EncodeRGB(c.Gamma())
		for i := 0; i < ledram.Len(); i++ {
			ledram.Clear()
			ledram.At(i).Write(pixel)
			tts.Write(ledram.Bytes())
			delay.Millisec(20)
		}
		for i := 0; i < ledram.Len(); i++ {
			ledram.At(i).Write(pixel)
		}
		tts.Write(ledram.Bytes())
		delay.Millisec(450)
	}
}

func ttsISR() {
	tts.ISR()
}

func ttsTxDMAISR() {
	tts.TxDMAISR()
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.USART2:        ttsISR,
	irq.DMA1_Channel7: ttsTxDMAISR,
}
