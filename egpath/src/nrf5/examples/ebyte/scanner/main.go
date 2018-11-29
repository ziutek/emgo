package main

import (
	"delay"
	"fmt"
	"rtos"
	"sync/fence"
	"unsafe"

	"nrf5/hal/clock"
	"nrf5/hal/gpio"
	"nrf5/hal/irq"
	"nrf5/hal/radio"
	"nrf5/hal/rtc"
	"nrf5/hal/system"
	"nrf5/hal/system/timer/rtcst"
	"nrf5/hal/te"
	"nrf5/hal/uart"
)

var (
	u          *uart.Driver
	radioEvent rtos.EventFlag
)

func init() {
	system.Setup(clock.XTAL, clock.XTAL, true)
	rtcst.Setup(rtc.RTC1, 0)

	p0 := gpio.P0

	u = uart.NewDriver(uart.UART0, make([]byte, 80))
	u.P.StorePSEL(uart.RXD, p0.Pin(11))
	u.P.StorePSEL(uart.TXD, p0.Pin(12))
	u.P.StoreBAUDRATE(uart.Baud115200)
	u.P.StoreENABLE(true)
	rtos.IRQ(u.P.NVIRQ()).Enable()
	u.EnableTx()
	fmt.DefaultWriter = u

	rtos.IRQ(radio.RADIO.NVIRQ()).Enable()
}

func main_() {
	data := make([]byte, 16)
	addr := uint32(0x12345678)

	r := radio.RADIO
	r.StorePCNF1(
		radio.MaxLen(16) | radio.StatLen(16) | radio.BALen(2) | radio.MSBFirst,
	)
	r.StoreBASE(0, addr<<8)
	r.StorePREFIX(0, addr>>24)
	r.StoreTXADDRESS(0)
	r.StorePACKETPTR(unsafe.Pointer(&data[0]))
	r.StoreMODE(radio.NRF_1M)
	r.StoreTXPOWER(0)
	//r.StoreTEST(true, true)

	for ch := 0; ; ch = (ch + 1) & 0x7F {
		r.StoreFREQUENCY(radio.Channel(ch))
		r.StoreSHORTS(radio.READY_START | radio.END_START)
		r.Task(radio.TXEN).Trigger()

		fmt.Printf("Freq: %d\r\n", 2400+ch)
		delay.Millisec(1000)

		disabled := r.Event(radio.DISABLED)
		disabled.Clear()
		disabled.EnableIRQ()
		radioEvent.Reset(0)
		fence.W()
		r.StoreSHORTS(radio.READY_START | radio.END_DISABLE)
		radioEvent.Wait(1, 0)
	}
}

func main() {
	data := make([]byte, 16)

	r := radio.RADIO
	r.StorePCNF1(
		radio.MaxLen(16) | radio.StatLen(16) | radio.BALen(2) | radio.MSBFirst,
	)
	r.StorePACKETPTR(unsafe.Pointer(&data[0]))
	r.StoreMODE(radio.NRF_1M)
	r.StoreSHORTS(radio.READY_START | radio.DISABLED_RSSISTOP)

	scan := make([]int, 0x7F+1)
	const numScan = 16

	for {
		for i := 0; i < numScan; i++ {
			for ch := range scan {
				r.StoreFREQUENCY(radio.Channel(ch))

				ev := r.Event(radio.READY)
				ev.Clear()
				ev.EnableIRQ()
				radioEvent.Reset(0)
				fence.W()

				r.Task(radio.RXEN).Trigger()
				radioEvent.Wait(1, 0)

				ev = r.Event(radio.RSSIEND)
				ev.Clear()
				ev.EnableIRQ()
				radioEvent.Reset(0)
				fence.W()

				r.Task(radio.RSSISTART).Trigger()
				radioEvent.Wait(1, 0)
				scan[ch] += r.LoadRSSISAMPLE()

				ev = r.Event(radio.DISABLED)
				ev.Clear()
				ev.EnableIRQ()
				radioEvent.Reset(0)
				fence.W()

				r.Task(radio.DISABLE).Trigger()
				radioEvent.Wait(1, 0)
			}
		}
		for ch, rssi := range scan {
			n := (103*numScan + rssi) / (2 * numScan)
			if n < 0 {
				n = 0
			} else if n > 30 {
				n = 30
			}
			scan[ch] = 1 << uint(n)
		}
		u.WriteString("\r\n\r\n\r\n\r\n\r\n")
		for i := 28; i >= 0; i -= 2 {
			for _, rssi := range scan {
				b := byte(' ')
				switch {
				case rssi>>uint(i+1) > 0:
					b = ':'
				case rssi>>uint(i) > 0:
					b = '.'
				}
				u.WriteByte(b)
			}
			u.WriteString("\r\n")
		}
		for ch := range scan {
			scan[ch] = 0
		}
	}
}

func uartISR() {
	u.ISR()
}

func radioISR() {
	radio.RADIO.DisableIRQ(te.EvAll)
	radioEvent.Signal(1)
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTC1:  rtcst.ISR,
	irq.UART0: uartISR,
	irq.RADIO: radioISR,
}
