package main

import (
	"fmt"
	"rtos"

	"bluetooth/ble"

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

	bctr = blec.NewCtrl(ble.MaxDataPay, 3, 3)
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
	fmt.Printf("\r\nDevAddr: %08x\r\n", uint64(bctr.DevAddr()))

	pdu := ble.MakeAdvPDU(ble.MaxDataPay)
	pdu.SetType(ble.ScanRsp)
	pdu.SetTxAdd(bctr.DevAddr() < 0)
	pdu.AppendAddr(bctr.DevAddr())
	pdu.AppendString(ble.LocalName, "Emgo & nRF5")
	pdu.AppendBytes(ble.TxPower, 0)
	bctr.Advertise(pdu, 625)
	for {
		pdu, _ := bctr.Recv()
		i := bctr.Iter
		if pdu.PayLen() > ble.MaxDataPay {
			fmt.Printf("error\r\n")
			continue
		}
		fmt.Printf(
			"%d LLID=%x P=%02x\r\n",
			i, pdu.Header()&ble.LLID, pdu.Payload())
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
