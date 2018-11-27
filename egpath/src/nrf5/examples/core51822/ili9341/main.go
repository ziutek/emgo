package main

import (
	"bufio"
	"delay"
	"fmt"
	"rtos"
	"text/linewriter"

	"display/ili9341"
	"display/ili9341/ili9341test"

	"nrf5/ilidci"

	"nrf5/hal/clock"
	"nrf5/hal/gpio"
	"nrf5/hal/irq"
	"nrf5/hal/rtc"
	"nrf5/hal/spi"
	"nrf5/hal/system"
	"nrf5/hal/system/timer/rtcst"
	"nrf5/hal/uart"
)

var (
	leds   [5]gpio.Pin
	u      *uart.Driver
	lcdspi *spi.Driver
	lcd    *ili9341.Display
)

func init() {
	system.Setup(clock.XTAL, clock.XTAL, true)
	rtcst.Setup(rtc.RTC0, 1)

	// GPIO

	p0 := gpio.P0

	ilireset := p0.Pin(0)
	ilidc := p0.Pin(1)
	ilimosi := p0.Pin(2)
	ilisck := p0.Pin(3)
	ilimiso := p0.Pin(4)

	utx := p0.Pin(9)
	urx := p0.Pin(11)

	for i := range leds {
		leds[i] = p0.Pin(18 + i)
	}

	ilicsn := p0.Pin(30)

	// LEDs

	for _, led := range leds {
		led.Setup(gpio.ModeOut)
	}

	// UART

	u = uart.NewDriver(uart.UART0, make([]byte, 80))
	u.P.StorePSEL(uart.RXD, urx)
	u.P.StorePSEL(uart.TXD, utx)
	u.P.StoreBAUDRATE(uart.Baud115200)
	u.Enable()
	//u.EnableRx()
	u.EnableTx()
	rtos.IRQ(u.P.NVIRQ()).Enable()
	fmt.DefaultWriter = linewriter.New(
		bufio.NewWriterSize(u, 80),
		linewriter.CRLF,
	)

	// LCD SPI

	lcdspi = spi.NewDriver(spi.SPI0)
	lcdspi.P.StorePSEL(spi.SCK, ilisck)
	lcdspi.P.StorePSEL(spi.MISO, ilimiso)
	lcdspi.P.StorePSEL(spi.MOSI, ilimosi)
	rtos.IRQ(lcdspi.P.NVIRQ()).Enable()

	// LCD controll

	p0.Setup(ilicsn.Mask()|ilidc.Mask()|ilireset.Mask(), gpio.ModeOut)
	ilicsn.Set()
	delay.Millisec(1) // Reset pulse.
	ilireset.Set()
	delay.Millisec(5) // Wait for reset.
	ilicsn.Clear()

	lcd = ili9341.NewDisplay(ilidci.New(lcdspi, ilidc, spi.Freq8M), 240, 320)
	lcd.DCI().Setup()
}

func main() {
	ili9341test.Run(lcd, 4, true)
}

func spiISR() {
	lcdspi.ISR()
}

func uartISR() {
	u.ISR()
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTC0:      rtcst.ISR,
	irq.SPI0_TWI0: spiISR,
	irq.UART0:     uartISR,
}
