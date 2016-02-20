// Control of power (water heater, house heating system).
package main

import (
	"delay"
	"mmio"
	"rtos"

	"arch/cortexm/bitband"

	"stm32/hal/exti"
	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"

	"stm32/hal/raw/rcc"
	"stm32/hal/raw/tim"
)

// Interrupt inputs.
const (
	// Port A.
	button = gpio.Pin0
	// Port B.
	waterFlow = gpio.Pin9
	// Port C.
	oneWire = gpio.Pin10
)

const PWMmax = 1e4

// Outputs.
var (
	// Solid state relays.
	SSR *gpio.Port

	// Water heater PWM.
	Wpwm [3]*mmio.U32
)

// Bitband outputs.
var (
	// Onboard LEDs.
	Blue  bitband.Bit
	Green bitband.Bit

	RoomHeater [3]bitband.Bit // Room heating solid state relays.
)

func init() {
	system.Setup32(0)
	systick.Setup()

	// GPIO.

	gpio.A.EnableClock(true)
	btnport := gpio.A
	gpio.B.EnableClock(true)
	ledport := gpio.B
	wfport := gpio.B
	gpio.C.EnableClock(true)
	SSR = gpio.C

	// Inputs

	// Button.
	btnport.Setup(button, &gpio.Config{Mode: gpio.In})
	line := exti.Lines(button)
	line.Connect(btnport)
	line.EnableRiseTrig()
	line.EnableInt()
	rtos.IRQ(irq.EXTI0).Enable()

	// Water flow sensor.
	wfport.Setup(waterFlow, &gpio.Config{Mode: gpio.In, Pull: gpio.PullUp})
	line = exti.Lines(waterFlow)
	line.Connect(wfport)
	line.EnableFallTrig()
	line.EnableInt()
	rtos.IRQ(irq.EXTI9_5).Enable()

	// Outputs

	slowOut := gpio.Config{Mode: gpio.Out, Speed: gpio.Low}

	// LEDs.
	bbpins := ledport.OutPins()
	ledport.SetupPin(6, &slowOut)
	ledport.SetupPin(7, &slowOut)
	Blue = bbpins.Bit(6)
	Green = bbpins.Bit(7)

	// Room heating.
	for _, pin := range []int{0, 1, 2} {
		SSR.SetupPin(pin, &slowOut)
	}
	bbpins = SSR.OutPins()
	RoomHeater[0] = bbpins.Bit(0)
	RoomHeater[1] = bbpins.Bit(1)
	RoomHeater[2] = bbpins.Bit(2)

	// Water heating, PWM.
	const (
		pwmfreq = 2 // Hz
		pwmpins = gpio.Pin6 | gpio.Pin7 | gpio.Pin8
		pwmmode = 6 //  Mode 1
	)
	SSR.Setup(pwmpins, &gpio.Config{Mode: gpio.Alt, Speed: gpio.Low})
	SSR.SetAltFunc(pwmpins, gpio.TIM3)
	rcc.RCC.TIM3EN().Set()
	t := tim.TIM3
	t.PSC.U16.Store(uint16(system.APB1.Clock()/(PWMmax*pwmfreq) - 1))
	t.ARR.Store(PWMmax - 1)
	t.OC1M().Store(pwmmode << tim.OC1Mn)
	t.OC2M().Store(pwmmode << tim.OC2Mn)
	t.OC3M().Store(pwmmode << tim.OC3Mn)
	t.OC1PE().Set()
	t.OC2PE().Set()
	t.OC3PE().Set()
	t.CCER.SetBits(tim.CC1E | tim.CC2E | tim.CC3E)
	t.ARPE().Set()
	t.UG().Set()
	t.CEN().Set()

	Wpwm[0] = &t.CCR1.U32
	Wpwm[1] = &t.CCR2.U32
	Wpwm[2] = &t.CCR3.U32
}

// ISRs

func exti0() {
	exti.L0.ClearPending()
	Green.Set()
	delay.Loop(1e5)
	Green.Clear()
	//buttonISR()
}

func exti9_5() {
	p := exti.Pending()
	(exti.L9 | exti.L8 | exti.L7 | exti.L6 | exti.L5).ClearPending()
	if p&exti.Lines(waterFlow) != 0 {
		waterISR()
	}
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.EXTI0:   exti0,
	irq.EXTI9_5: exti9_5,
}

func main() {
	Wpwm[0].Store(PWMmax * 1 / 4)
	Wpwm[1].Store(PWMmax * 2 / 4)
	Wpwm[2].Store(PWMmax * 3 / 4)

	for {
		delay.Millisec(100)
	}
	//waterTask()
}
