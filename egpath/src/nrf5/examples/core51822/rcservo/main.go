// Remotely cotrolled servo. This exemple shows how to use ppipwm to controll
// locally connected servo. Additionally it shows how to use RADIO peripheral
// to receive servo settings. stm32/examples/minidev/nrfrc contains code of
// example transmitter.
package main

import (
	"fmt"
	"rtos"
	"sync/fence"
	"unsafe"

	"nrf5/ppipwm"

	"nrf5/hal/clock"
	"nrf5/hal/gpio"
	"nrf5/hal/gpiote"
	"nrf5/hal/irq"
	"nrf5/hal/ppi"
	"nrf5/hal/radio"
	"nrf5/hal/rtc"
	"nrf5/hal/system"
	"nrf5/hal/system/timer/rtcst"
	"nrf5/hal/timer"
	"nrf5/hal/te"
	"nrf5/hal/uart"
)

var (
	pwm        *ppipwm.Toggle
	u          *uart.Driver
	radioEvent rtos.EventFlag
)

func init() {
	system.Setup(clock.XTAL, clock.XTAL, true)
	rtcst.Setup(rtc.RTC0, 1)

	p0 := gpio.P0

	u = uart.NewDriver(uart.UART0, make([]byte, 80))
	u.P.StorePSEL(uart.RXD, p0.Pin(11))
	u.P.StorePSEL(uart.TXD, p0.Pin(9))
	u.P.StoreBAUDRATE(uart.Baud115200)
	u.P.StoreENABLE(true)
	rtos.IRQ(u.P.NVIRQ()).Enable()
	u.EnableTx()
	fmt.DefaultWriter = u

	r := radio.RADIO
	r.StorePCNF0(0)
	r.StorePCNF1(
		radio.MaxLen(2) | radio.StatLen(2) | radio.BALen(2) | radio.MSBFirst,
	)
	r.StoreCRCCNF(2, false)
	r.StoreCRCPOLY(1<<16 | 1<<12 | 1<<5 | 1)
	r.StoreCRCINIT(0xFFFF)
	r.StoreBASE(0, 0xE7E70000) // Reversed 0xE7E70000.
	r.StorePREFIX(0, 0xE7)     // Reversed 0xE7.
	r.StoreRXADDRESSES(1 << 0)
	r.StoreMODE(radio.NRF_250K)
	r.StoreFREQUENCY(radio.Channel(50))
	r.StoreSHORTS(radio.READY_START)
	rtos.IRQ(r.NVIRQ()).Enable()

	pwm = ppipwm.NewToggle(timer.TIMER1)
	pwm.SetFreq(8, 20e3) // Gives resolution of 1250 levels of duty cycle.
	pwm.Setup(0, p0.Pin(22), gpiote.Chan(0), ppi.Chan(0), ppi.Chan(1))
}

func main() {
	// For SG92R servo (PWM: 3.3V, 20 ms).
	var (
		min    = pwm.Max() * 600 / 20e3
		max    = pwm.Max() * 2400 / 20e3
		center = (min + max) / 2
	)

	data := make([]byte, 2)

	r := radio.RADIO
	r.StorePACKETPTR(unsafe.Pointer(&data[0]))

	start := r.Task(radio.RXEN)
	for {
		endev := r.Event(radio.END)
		endev.Clear()
		endev.EnableIRQ()
		radioEvent.Reset(0)
		fence.W()
		start.Trigger()
		radioEvent.Wait(1, 0)
		crcok := r.LoadCRCSTATUS()

		x, y := int(int8(data[0])), int(int8(data[1]))
		fmt.Printf("x=%d y=%d crc=%t\r\n", x, y, crcok)

		start = r.Task(radio.START)

		if !crcok {
			continue
		}
		switch {
		case x < -64:
			x = -64
		case x > 64:
			x = 64
		}
		pwm.Set(0, center+x*(max-min)/128)
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
	irq.RTC0:  rtcst.ISR,
	irq.UART0: uartISR,
	irq.RADIO: radioISR,
}
