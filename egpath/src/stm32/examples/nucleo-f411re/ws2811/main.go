// Simple WS2811 example. This example can work with F4-Discovery, but VDD=3V
// instead of 3.3V can be a problem.
package main

import (
	"delay"
	"fmt"
	"rtos"

	"ws281x"

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

	systick.Setup(2e6)

	// GPIO

	gpio.C.EnableClock(true)
	spiport, mosi := gpio.C, gpio.Pin12

	// SPI.
	cfg := gpio.Config{Mode: gpio.Alt, Speed: gpio.High}
	spiport.Setup(mosi, &cfg)
	spiport.SetAltFunc(mosi, gpio.SPI3)
	d := dma.DMA1
	d.EnableClock(true)
	wspi = spi.NewDriver(spi.SPI3, d.Channel(7, 0), nil)
	wspi.P.EnableClock(true)
	rtos.IRQ(irq.SPI3).Enable()
	rtos.IRQ(irq.DMA1_Stream7).Enable()
}

func main() {
	wspi.P.SetConf(spi.Master | wspi.P.BR(3200e3) | spi.SoftSS | spi.ISSHigh)
	wspi.P.Enable()
	delay.Millisec(250) // For SWO handling in ST-Link.

	fmt.Printf("\nSPI speed: %d Hz\n", wspi.P.Baudrate(wspi.P.Conf()))

	ledram := ws281x.MakeSPIFB(50)
	pixel := ws281x.MakeSPIFB(1)
	colors := []ws281x.Color{
		0x99aadd,
		0xddaa99,
	}

	for _, c := range colors {
		fmt.Printf("%3d %3d %3d\n", c.Red(), c.Green(), c.Blue())
		pixel.EncodeRGB(c.Gamma())
		for i := 0; i < ledram.Len(); i++ {
			ledram.Clear()
			ledram.At(i).Write(pixel)
			wspi.WriteRead(ledram.Bytes(), nil)
			delay.Millisec(20)
		}
		for i := 0; i < ledram.Len(); i++ {
			ledram.At(i).Write(pixel)
		}
		wspi.WriteRead(ledram.Bytes(), nil)
		delay.Millisec(450)
	}
	delay.Millisec(2000)
	ledram.Clear()
	wspi.WriteRead(ledram.Bytes(), nil)
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
