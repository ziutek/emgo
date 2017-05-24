package main

import (
	"delay"
	"fmt"
	"rtos"

	"ble"

	"nrf5/blec"
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
	bctr *blec.Ctrl
	udrv *uart.Driver
)

func init() {
	system.Setup(clock.XTAL, clock.XTAL, true)
	rtcst.Setup(rtc.RTC1, 0)

	p0 := gpio.P0

	for i := range leds {
		led := p0.Pin(18 + i)
		led.Setup(gpio.ModeOut)
		leds[i] = led
	}

	bctr = blec.NewCtrl()
	bctr.InitPhy()
	bctr.LEDs = &leds

	udrv = uart.NewDriver(uart.UART0, make([]byte, 80))
	udrv.P.StorePSEL(uart.SignalRXD, p0.Pin(11))
	udrv.P.StorePSEL(uart.SignalTXD, p0.Pin(9))
	udrv.P.StoreBAUDRATE(uart.Baud115200)
	udrv.P.StoreENABLE(true)
	udrv.EnableTx()
	rtos.IRQ(udrv.P.IRQ()).Enable()
	fmt.DefaultWriter = udrv
}

func main() {
	fmt.Printf("DevAddr: %08x\r\n", uint64(bctr.DevAddr()))

	pdu := ble.MakeAdvPDU(nil)
	pdu.SetType(ble.ScanRsp)
	pdu.SetTxAdd(bctr.DevAddr() < 0)
	pdu.AppendAddr(bctr.DevAddr())
	pdu.AppendString(ble.LocalName, "Emgo & nRF5")
	pdu.AppendBytes(ble.TxPower, 0)
	bctr.Advertise(pdu, 625)
	for i := 0; ; i++ {
		fmt.Printf(
			"CC=%d CNT=%d\r\n",
			rtc.RTC0.LoadCC(0), rtc.RTC0.LoadCOUNTER(),
		)
		delay.Millisec(625)
	}
}

func radioISR() {
	bctr.RadioISR()
}

func uartISR() {
	udrv.ISR()
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTC1: rtcst.ISR,

	irq.RADIO: radioISR,

	irq.UART0: uartISR,
}
