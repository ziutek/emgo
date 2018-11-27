package main

import (
	"rtos"

	"nrf5/hal/clock"
	"nrf5/hal/gpio"
	"nrf5/hal/irq"
	"nrf5/hal/rtc"
	"nrf5/hal/spi"
	"nrf5/hal/system"
	"nrf5/hal/system/timer/rtcst"
)

var (
	leds [5]gpio.Pin
	d    *spi.Driver
	cs   gpio.Pin
	dc   gpio.Pin
	rst  gpio.Pin
	busy gpio.Pin
)

func init() {
	system.Setup(clock.XTAL, clock.XTAL, true)
	rtcst.Setup(rtc.RTC0, 1)

	p0 := gpio.P0
	for i := range leds {
		led := p0.Pin(18 + i)
		led.Setup(gpio.ModeOut)
		leds[i] = led
	}
	cs = p0.Pin(25)
	dc = p0.Pin(23)
	rst = p0.Pin(21)
	busy = p0.Pin(22)

	d = spi.NewDriver(spi.SPI0)
	d.P.StorePSEL(spi.SCK, p0.Pin(27))
	d.P.StorePSEL(spi.MOSI, p0.Pin(29))
	d.P.StoreFREQUENCY(spi.Freq125k)
	d.Enable()
	rtos.IRQ(d.P.NVIRQ()).Enable()
}

func main() {

}

func spiISR() {
	d.ISR()
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTC0:      rtcst.ISR,
	irq.SPI0_TWI0: spiISR,
}
