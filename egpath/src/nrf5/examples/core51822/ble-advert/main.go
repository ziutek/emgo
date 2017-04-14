// BLE advertising example.
//
// Based on:
// https://github.com/pauloborges/blessed/tree/devel/examples/radio-broadcaster
//
// Install nRF Connect application and run scanner. Your phone/tablet should
// find BLE device with local name "Emgo & nRF5" (tested on LG G2).
package main

import (
	//"debug/semihosting"
	"delay"
	"rtos"
	"sync/fence"
	"unsafe"

	"nrf5/hal/clock"
	"nrf5/hal/ficr"
	"nrf5/hal/gpio"
	"nrf5/hal/irq"
	"nrf5/hal/radio"
	"nrf5/hal/rtc"
	"nrf5/hal/system"
	"nrf5/hal/system/timer/rtcst"
)

var (
	leds       [5]gpio.Pin
	radioEvent rtos.EventFlag
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

	r := radio.RADIO
	if f := ficr.FICR; f.BLE_1MBIT_OK().Load() == 0 {
		r.StoreOVERRIDE(0, f.BLE_1MBIT[0].U32.Load())
		r.StoreOVERRIDE(1, f.BLE_1MBIT[1].U32.Load())
		r.StoreOVERRIDE(2, f.BLE_1MBIT[2].U32.Load())
		r.StoreOVERRIDE(3, f.BLE_1MBIT[3].U32.Load())
		r.StoreOVERRIDE(4, f.BLE_1MBIT[4].U32.Load()|1<<31)
	}
	r.StoreMODE(radio.BLE_1M)
	r.StoreTIFS(150)
	r.StorePCNF0(radio.S0_8b | radio.LenBitN(8))
	r.StorePCNF1(radio.WhiteEna | radio.MaxLen(39-2) | radio.BALen(3))
	//r.StoreRXADDRESSES(1)
	r.StoreTXADDRESS(0)
	r.StoreCRCCNF(3, true)
	r.StoreCRCPOLY(0x100065B)
	r.StoreSHORTS(radio.READY_START | radio.END_DISABLE)
	rtos.IRQ(r.IRQ()).Enable()

	//f, err := semihosting.OpenFile(":tt", semihosting.W)
	//for err != nil {
	//}
	//fmt.DefaultWriter = f
}

func main() {
	channels := [3]radio.Freq{
		radio.Channel(2),  // BLE 37, 2402 MHz
		radio.Channel(26), // BLE 38, 2426 MHz
		radio.Channel(80), // BLE 39, 2480 MHz
	}
	data := []byte{
		0x42,                               // S0
		19,                                 // Length
		0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, // AdvA
		12,   // AD Length
		0x09, // AD Type
		'E', 'm', 'g', 'o',
		' ', '&', ' ',
		'n', 'R', 'F', '5',
	}
	aa := uint32(0x8E89BED6)

	r := radio.RADIO
	r.StorePACKETPTR(unsafe.Pointer(&data[0]))
	r.StoreBASE(0, aa<<8)
	r.StorePREFIX(0, aa>>24)
	r.StoreCRCINIT(0x555555)

	for {
		for n, c := range channels {
			r.StoreFREQUENCY(c)
			r.StoreDATAWHITEIV(uint32(37+n) & 0x3F)

			disev := r.Event(radio.DISABLED)
			disev.Clear()
			disev.EnableIRQ()
			radioEvent.Reset(0)
			fence.W()
			r.Task(radio.TXEN).Trigger()
			radioEvent.Wait(1, 0)
			leds[4].Store(n)
			delay.Millisec(25)
		}
	}
}

func radioISR() {
	radio.RADIO.DisableIRQ(0xFFFFFFFF)
	radioEvent.Signal(1)
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTC0:  rtcst.ISR,
	irq.RADIO: radioISR,
}
