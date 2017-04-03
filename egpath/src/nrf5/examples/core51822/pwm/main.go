package main

import (
	"delay"
	"fmt"
	"rtos"

	"nrf5/hal/clock"
	"nrf5/hal/gpio"
	"nrf5/hal/gpiote"
	"nrf5/hal/irq"
	"nrf5/hal/ppi"
	"nrf5/hal/rtc"
	"nrf5/hal/system"
	"nrf5/hal/system/timer/rtcst"
	"nrf5/hal/timer"
	"nrf5/hal/uart"
)

type PWM struct {
	pin gpio.Pin
	gc  gpiote.Chan
	t   *timer.Periph
	max int
}

func (pwm *PWM) Init(p gpio.Pin, gc gpiote.Chan, t *timer.Periph, pc0, pc1 ppi.Chan) {
	pwm.pin = p
	pwm.t = t
	pwm.gc = gc
	p.Clear()
	p.Setup(gpio.ModeOut)
	t.Task(timer.STOP).Trigger()
	t.StoreSHORTS(timer.COMPARE1_CLEAR)
	pc0.SetEEP(t.Event(timer.COMPARE0))
	pc0.SetTEP(gc.OUT())
	pc0.Enable()
	pc1.SetEEP(t.Event(timer.COMPARE1))
	pc1.SetTEP(gc.OUT())
	pc1.Enable()
}

// SetFreq sets prescaler to 2^pre and period (microseconds).
func (pwm *PWM) SetFreq(pre, period int) {
	if uint(pre) > 9 {
		panic("PWM: bad prescaler")
	}
	if period < 10 {
		panic("PWM: bad period")
	}
	t := pwm.t
	t.StorePRESCALER(pre)
	div := uint32(1) << uint(pre)
	max := 16*uint32(period)/div - 1
	if max > 0xFFFF {
		panic("PWM: bad pre and/or period for 16-bit timer")
	}
	t.StoreCC(1, max)
	pwm.max = int(max)
}

//  Max returns value that represents 100% duty cycle.
func (pwm *PWM) Max() int {
	return pwm.max
}

func (pwm *PWM) SetDutyCycle(dc int) {
	pin := pwm.pin
	gc := pwm.gc
	t := pwm.t
	if dc >= pwm.max {
		pin.Set()
		gc.Setup(pin, 0)
		return
	}
	pin.Clear()
	gc.Setup(pin, 0)
	if dc == 0 {
		return
	}
	t.Task(timer.STOP).Trigger()
	t.Task(timer.CLEAR).Trigger()
	t.StoreCC(0, uint32(dc))
	gc.Setup(pin, gpiote.ModeTask|gpiote.PolarityToggle|gpiote.OutInitHigh)
	t.Task(timer.START).Trigger()
}

var (
	pwm PWM
	u   *uart.Driver
)

func init() {
	system.Setup(clock.XTAL, clock.XTAL, true)
	rtcst.Setup(rtc.RTC0, 1)

	p0 := gpio.P0

	u = uart.NewDriver(uart.UART0, make([]byte, 80))
	u.P.StorePSEL(uart.SignalRXD, p0.Pin(11))
	u.P.StorePSEL(uart.SignalTXD, p0.Pin(9))
	u.P.StoreBAUDRATE(uart.Baud115200)
	u.P.StoreENABLE(true)
	rtos.IRQ(irq.UART0).Enable()
	u.EnableTx()
	fmt.DefaultWriter = u

	pwm.Init(
		p0.Pin(22), gpiote.Chan(0),
		timer.TIMER1,
		ppi.Chan(0), ppi.Chan(1),
	)
	pwm.SetFreq(1, 4000)
}

func main() {
	max := pwm.Max()
	v, a, b, c := 0, 4, 5, 0
	for {
		fmt.Printf("%d/%d\r\n", v, max)
		pwm.SetDutyCycle(v)
		switch {
		case v == 0:
			c = a
			a, b = b, a
		case v > (max*5+4)/4:
			c = 0
			a, b = b, a
		}
		v = (v*a + c) / b
		delay.Millisec(200)
	}
}

func uartISR() {
	u.ISR()
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTC0:  rtcst.ISR,
	irq.UART0: uartISR,
}
