// Simple WS2811 example. This example can work with F4-Discovery, but VDD=3V
// instead of 3.3V can be a problem.
package main

import (
	"delay"
	"fmt"
	"rtos"

	"stm32/hal/dma"
	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/spi"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
)

var wspi *spi.Driver

func init() {
	// For SPI clock closer to required 3.2 MHz use 102 MHz SysClk.
	//system.Setup(8, 102, 2)

	// This gives 3 MHz SPI clock (slower but seems to work more realiable).
	system.Setup96(8)

	systick.Setup()

	// GPIO

	gpio.C.EnableClock(true)
	spiport, mosi := gpio.C, gpio.Pin12

	// SPI.
	cfg := gpio.Config{Mode: gpio.Alt, Speed: gpio.High}
	spiport.Setup(mosi, &cfg)
	spiport.SetAltFunc(mosi, gpio.SPI3)
	d := dma.DMA1
	d.EnableClock(true)
	wspi = spi.NewDriver(spi.SPI3, nil, d.Channel(7, 0))
	wspi.P.EnableClock(true)
	rtos.IRQ(irq.SPI3).Enable()
	rtos.IRQ(irq.DMA1_Stream7).Enable()
}

type color struct {
	r, g, b byte
}

func encodeRGB(buf []byte, c color) {
	r := gamma[c.r]
	g := gamma[c.g]
	b := gamma[c.b]
	for n := uint(0); n < 4; n++ {
		buf[3-n] = 0x88 + r>>(2*n+1)&1<<6 + r>>(2*n)&1<<4
		buf[7-n] = 0x88 + g>>(2*n+1)&1<<6 + g>>(2*n)&1<<4
		buf[11-n] = 0x88 + b>>(2*n+1)&1<<6 + b>>(2*n)&1<<4
	}
}

func reverse(buf []byte) {
	for i, b := range buf {
		buf[i] = ^b
	}
	//buf[len(buf)-1] = 0xff
}

func main() {
	wspi.P.SetConf(spi.Master | wspi.P.BR(3200e3) | spi.SoftSS | spi.ISSHigh)
	wspi.P.Enable()
	delay.Millisec(250) // For SWO handling in ST-Link.

	fmt.Printf("\nSPI speed: %d Hz\n", wspi.P.Baudrate(wspi.P.Conf()))

	N := 50
	pixels := make([]byte, N*12+1)
	colors := []color{
		{250, 220, 0},
		{135, 0, 135},
		{123, 222, 245},
		{200, 200, 200},
		{200, 0, 0},
		{0, 200, 0},
		{0, 0, 200},
		{200, 0, 200},
		{200, 200, 0},
		{0, 200, 200},
	}

	for _, c := range colors {
		for i := 0; i < 50; i++ {
			for n := 0; n < N; n++ {
				if n == i {
					encodeRGB(pixels[n*12:], c)
				} else {
					encodeRGB(pixels[n*12:], color{})
				}
			}
			wspi.WriteRead(pixels, nil)
			//fmt.Printf(".")
			delay.Millisec(20)
		}
		for n := 0; n < N; n++ {
			encodeRGB(pixels[n*12:], c)
		}
		wspi.WriteRead(pixels, nil)
		delay.Millisec(500)
	}
	for i := 255; i >= 0; i-- {
		for n := 0; n < N; n++ {
			encodeRGB(pixels[n*12:], color{byte(i), byte(i), byte(i)})
		}
		wspi.WriteRead(pixels, nil)
		delay.Millisec(50)
	}
	fmt.Printf("End.\n")
}

func spiISR() {
	wspi.ISR()
}

func spiTxDMAISR() {
	wspi.DMAISR(wspi.TxDMA)
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.SPI3:         spiISR,
	irq.DMA1_Stream7: spiTxDMAISR,
}

// Gamma table according to CIE 1976.
//
//	const max = 255
//
//	for i := 0; i <= max; i++ {
//		var x float64
//		y := 100 * float64(i) / max
//		if y > 8 {
//			x = math.Pow((y+16)/116, 3)
//		} else {
//			x = y / 903.3
//		}
//		fmt.Printf("%d,", int(max*x+0.5))
//	}
//
//emgo.const
var gamma = [256]byte{
	0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 2, 3, 3,
	3, 3, 3, 3, 3, 3, 4, 4, 4, 4, 4, 4, 5, 5, 5, 5, 5, 6, 6, 6, 6, 6, 7, 7, 7,
	7, 8, 8, 8, 8, 9, 9, 9, 10, 10, 10, 10, 11, 11, 11, 12, 12, 12, 13, 13, 13,
	14, 14, 15, 15, 15, 16, 16, 17, 17, 17, 18, 18, 19, 19, 20, 20, 21, 21, 22,
	22, 23, 23, 24, 24, 25, 25, 26, 26, 27, 28, 28, 29, 29, 30, 31, 31, 32, 32,
	33, 34, 34, 35, 36, 37, 37, 38, 39, 39, 40, 41, 42, 43, 43, 44, 45, 46, 47,
	47, 48, 49, 50, 51, 52, 53, 54, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64,
	65, 66, 67, 68, 70, 71, 72, 73, 74, 75, 76, 77, 79, 80, 81, 82, 83, 85, 86,
	87, 88, 90, 91, 92, 94, 95, 96, 98, 99, 100, 102, 103, 105, 106, 108, 109,
	110, 112, 113, 115, 116, 118, 120, 121, 123, 124, 126, 128, 129, 131, 132,
	134, 136, 138, 139, 141, 143, 145, 146, 148, 150, 152, 154, 155, 157, 159,
	161, 163, 165, 167, 169, 171, 173, 175, 177, 179, 181, 183, 185, 187, 189,
	191, 193, 196, 198, 200, 202, 204, 207, 209, 211, 214, 216, 218, 220, 223,
	225, 228, 230, 232, 235, 237, 240, 242, 245, 247, 250, 252, 255,
}
