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
	
	// Set UART baudrate to 2250 kb/s (36 MHz / 16 = 2.25 MHz). This gives
	// 444 ns/UARTbit and 1333 ns/WS2811bit. It would be best to use the 7-bit
	// mode but F103 does not support it. Use 0.5 stop bit to slightly speed-up
	// transmission.
	
	// Edit: It seems that 1333 ns/WS2811bit is wrong. According to datasheet
	// WS2811 bit takes 2500±300 ns. However, this timing works but it is WS2812
	// timing (WS2812 bit takes 1390 ns).
	
	tts = usart.NewDriver(usart.USART2, d.Channel(7, 0), nil, nil)
	tts.Periph().EnableClock(true)
	tts.Periph().SetBaudRate(2250e3)    
	tts.Periph().SetConf2(usart.Stop0b5)
	tts.Periph().Enable()
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
