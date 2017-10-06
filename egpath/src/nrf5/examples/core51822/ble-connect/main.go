package main

import (
	"fmt"
	"rtos"

	"bluetooth/att"
	"bluetooth/ble"
	"bluetooth/l2cap"

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
	bctr *blec.Controller
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

	bctr = blec.NewController(ble.MaxDataPay, 3, 3)
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

type pduLogger struct {
	*blec.Controller
}

func (p pduLogger) Recv() (ble.DataPDU, error) {
	pdu, err := p.Controller.Recv()
	if err == nil {
		fmt.Printf(
			"\r\nR LL#%d LLID=%x P=%02x\r\n",
			p.ConnEventCnt(), pdu.Header()&ble.LLID, pdu.Payload(),
		)
	}
	return pdu, err
}

func main() {
	fmt.Printf("\r\nDevAddr: %08x\r\n", uint64(bctr.DevAddr()))

	pdu := ble.MakeAdvPDU(ble.MaxAdvPay)
	pdu.SetType(ble.ScanRsp)
	pdu.SetTxAdd(bctr.DevAddr() < 0)
	pdu.AppendAddr(bctr.DevAddr())
	pdu.AppendString(ble.LocalName, "Emgo BLE")
	pdu.AppendUUIDs(ble.Services, srvNordicUART)
	pdu.AppendBytes(ble.TxPower, 0)
	bctr.Advertise(pdu, 625)

	far := l2cap.NewBLEFAR(pduLogger{bctr})
	srv := att.NewServer(23)
	srv.SetHandler(gattSrv)
	for {
		cid, err := far.ReadHeader()
		if err != nil {
			fmt.Printf("ReadHeader: %v\r\n", err)
			continue
		}
		fmt.Printf("R L2CAP header: len=%d cid=%d\r\n", far.Len(), cid)
		switch cid {
		case 4: // ATT
			srv.HandleTransaction(far, cid)
		}
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
	irq.RTC1:  rtcst.ISR,
	irq.RADIO: radioISR,
	irq.UART0: uartISR,
}
