package main

import (
	//"debug/semihosting"
	"fmt"
	"rtos"

	"nrf5/hal/clock"
	"nrf5/hal/gpio"
	"nrf5/hal/irq"
	"nrf5/hal/rtc"
	"nrf5/hal/system"
	"nrf5/hal/system/timer/rtcst"
	"nrf5/hal/uart"
)

var (
	leds [5]gpio.Pin
	blep *Periph
	udrv *uart.Driver
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

	blep = NewPeriph("Emgo & nRF5", -4)

	udrv = uart.NewDriver(uart.UART0, make([]byte, 80))
	udrv.P.StorePSEL(uart.SignalRXD, p0.Pin(11))
	udrv.P.StorePSEL(uart.SignalTXD, p0.Pin(9))
	udrv.P.StoreBAUDRATE(uart.Baud115200)
	udrv.P.StoreENABLE(true)
	udrv.EnableTx()
	rtos.IRQ(udrv.P.IRQ()).Enable()
	fmt.DefaultWriter = udrv

	//f, err := semihosting.OpenFile(":tt", semihosting.W)
	//for err != nil {
	//}
	//fmt.DefaultWriter = f
}

func main() {
	fmt.Printf("DevAddr: %08x\r\n", uint64(getDevAddr()))
	blep.StartAdvert()
	for scanReq := range blep.scanReq {
		fmt.Printf("%d: %02x\r\n", rtos.Nanosec(), scanReq)
	}
}

func radioISR() {
	blep.ISR()
}

func uartISR() {
	udrv.ISR()
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTC0:  rtcst.ISR,
	irq.RADIO: radioISR,
	irq.UART0: uartISR,
}
