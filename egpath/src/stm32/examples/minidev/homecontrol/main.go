package main

import (
	"delay"
	"fmt"

	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/rtcst"
	"stm32/hal/tim"

	"stm32/hal/raw/afio"
	"stm32/hal/raw/rcc"
)

var (
	led    gpio.Pin
	relays [4]gpio.Pin
	pwm    tim.PWM
)

func init() {
	system.SetupPLL(8, 1, 72/8)
	rtcst.Setup(32768)

	// Allocate pins.

	gpio.B.EnableClock(false)
	relays[3] = gpio.B.Pin(4)
	relays[2] = gpio.B.Pin(5)
	relays[1] = gpio.B.Pin(6)
	relays[0] = gpio.B.Pin(7)
	ssr := gpio.B.Pin(8) // TIM4.CC3

	gpio.C.EnableClock(false)
	led = gpio.C.Pin(13)

	// Configure pins and peripherals..

	// Release JTDI and NJTRST (PA15 and PB4) to use as GPIO pins.
	rcc.RCC.AFIOEN().Set()
	afio.AFIO.SWJ_CFG().Store(afio.SWJ_CFG_JTAGDISABLE)
	rcc.RCC.AFIOEN().Clear()

	cfg := &gpio.Config{Mode: gpio.Out, Driver: gpio.OpenDrain, Speed: gpio.Low}
	led.Setup(cfg)
	for _, pin := range relays {
		pin.Set()
		pin.Setup(cfg)
	}
	ssr.Setup(&gpio.Config{
		Mode:   gpio.Alt, // TIM4.CC3
		Driver: gpio.PushPull,
		Speed:  gpio.Low,
	})
	pwm = tim.PWM{tim.TIM4}
	pwm.P.EnableClock(true)
	pwm.SetMode(tim.OCPWM1, tim.OCPWM1, tim.OCPWM1, tim.OCPWM1)
	pwm.SetPolarity(0, 0, 1, 0)
	pwm.SetFreq(5e5, 1e4)
	pwm.Enable(1)
}

func main() {
	fmt.Printf("HCLK=%d PCLK=%d\n", system.AHB.Clock(), pwm.P.Bus().Clock())
	pwm.Ch(tim.CC3).Store(1e3)
	for _, relay := range relays {
		led.Clear()
		relay.Clear()
		delay.Millisec(50)
		led.Set()
		delay.Millisec(950)
	}
}

//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTCAlarm: rtcst.ISR,
}
