package main

import (
	"delay"
	"fmt"
	"io"
	"rtos"

	"ws281x"

	"stm32/hal/dma"
	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
	"stm32/hal/usart"
)

var leds, us100 *usart.Driver

func init() {
	system.Setup(8, 72/8, false)
	systick.Setup()

	// GPIO

	gpio.A.EnableClock(true)
	ledport, ledpin := gpio.A, gpio.Pin2

	gpio.B.EnableClock(true)
	usport, ustx, usrx := gpio.B, gpio.Pin10, gpio.Pin11

	// LED USART

	ledport.Setup(ledpin, &gpio.Config{Mode: gpio.Alt})
	d := dma.DMA1
	d.EnableClock(true)
	leds = usart.NewDriver(usart.USART2, nil, d.Channel(7, 0), nil)
	leds.P.EnableClock(true)
	leds.P.SetBaudRate(2250e3)    // 36MHz/16: 444 ns/UARTbit, 1333 ns/WS2811bit.
	leds.P.SetConf(usart.Stop0b5) // F103 has no 7-bit mode: save 0.5 bit only.
	leds.P.Enable()
	leds.EnableTx()
	rtos.IRQ(irq.USART1).Enable()
	rtos.IRQ(irq.DMA1_Channel7).Enable()

	// US-100 USART

	usport.Setup(ustx, &gpio.Config{Mode: gpio.Alt})
	usport.Setup(usrx, &gpio.Config{Mode: gpio.AltIn, Pull: gpio.PullUp})
	d = dma.DMA1
	d.EnableClock(true)
	us100 = usart.NewDriver(
		usart.USART3, d.Channel(3, 0), d.Channel(2, 0), make([]byte, 8),
	)
	us100.P.EnableClock(true)
	us100.P.SetBaudRate(9600)
	us100.P.Enable()
	us100.EnableRx()
	us100.EnableTx()
	rtos.IRQ(irq.USART3).Enable()
	rtos.IRQ(irq.DMA1_Channel2).Enable()
	rtos.IRQ(irq.DMA1_Channel3).Enable()
}

func checkErr(err error) {
	if err != nil {
		fmt.Printf("\nError: %v\n", err)
	}
}

func main() {
	buf := make([]byte, 2)
	ledram := ws281x.MakeS3(50)
	pixel := ws281x.MakeS3(1)

	pixel.EncodeRGB(ws281x.Color(0x888888).Gamma())
	delay.Millisec(200) // Wait for US-100 startup.

	var x int
	for {
		checkErr(us100.WriteByte(0x55))
		_, err := io.ReadFull(us100, buf)
		checkErr(err)
		x = (x*3 + int(buf[0])<<8 + int(buf[1]) + 2) / 4
		fmt.Println(x)

		max := ledram.Len()
		n := max - (x-50)/32
		switch {
		case n < 0:
			n = 0
		case n > max:
			n = max
		}
		for i := 0; i < n; i++ {
			ledram.At(i).Write(pixel)
		}
		ledram.At(n).Clear()
		leds.Write(ledram.Bytes())
		delay.Millisec(25)
	}
}

func ledsISR() {
	leds.ISR()
}

func ledsTxDMAISR() {
	leds.TxDMAISR()
}

func us100ISR() {
	us100.ISR()
}

func us100RxDMAISR() {
	us100.RxDMAISR()
}

func us100TxDMAISR() {
	us100.TxDMAISR()
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.USART2:        ledsISR,
	irq.DMA1_Channel7: ledsTxDMAISR,

	irq.USART3:        us100ISR,
	irq.DMA1_Channel2: us100TxDMAISR,
	irq.DMA1_Channel3: us100RxDMAISR,
}
