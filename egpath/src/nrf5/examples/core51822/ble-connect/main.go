package main

import (
	"fmt"
	"rtos"

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

type pduLogger struct {
	c *blec.Ctrl
}

func (p pduLogger) Recv() (ble.DataPDU, error) {
	pdu, err := p.c.Recv()
	i := p.c.Iter
	if err == nil {
		fmt.Printf(
			"R LL PDU %d LLID=%x P=%02x\r\n",
			i, pdu.Header()&ble.LLID, pdu.Payload(),
		)
	}
	return pdu, err
}

func (p pduLogger) Send(pdu ble.DataPDU) error {
	return p.c.Send(pdu)
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

	far := new(l2cap.BLEFAR)
	far.SetHCI(pduLogger{bctr})
	buf := make([]byte, 128)
	for {
		n, cid, err := far.ReadHeader()
		if err != nil {
			fmt.Printf("ReadHeader: %v\r\n", err)
			continue
		}
		fmt.Printf("R L2CAP header: len=%d cid=%d\r\n", n, cid)
		for {
			m, err := far.Read(buf)
			if err != nil {
				fmt.Printf("Read: %v\r\n", err)
				break
			}
			if m == 0 {
				break
			}
			fmt.Printf("R L2CAP payload: %02x\r\n", buf[:m])
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
	irq.RTC1: rtcst.ISR,

	irq.RADIO: radioISR,

	irq.UART0: uartISR,
}
