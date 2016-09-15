// Control of power (water heater, house heating system).
package main

import (
	"delay"
	"fmt"
	"rtos"

	"arch/cortexm/bitband"
	"arch/cortexm/nvic"

	"stm32/hal/dma"
	"stm32/hal/exti"
	"stm32/hal/gpio"
	"stm32/hal/i2c"
	"stm32/hal/irq"
	"stm32/hal/raw/rcc"
	"stm32/hal/raw/tim"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
	"stm32/hal/usart"
)

var (
	// Onboard LEDs.
	ledBlue  bitband.Bit
	ledGreen bitband.Bit

	RoomHeater [3]bitband.Bit // Room heating solid state relays.
)

// prio16 must be in the range [0;15]. Do not use prio16 <= 4 for realtime ISRs.
func irqen(irq nvic.IRQ, prio16 rtos.IRQPrio) {
	e := rtos.IRQ(irq)
	e.SetPrio(rtos.IRQPrioLowest + prio16*rtos.IRQPrioStep*(rtos.IRQPrioNum/16))
	e.Enable()
}

func init() {
	system.Setup32(0)
	systick.Setup()

	// GPIO.

	gpio.A.EnableClock(true)
	btnport, btnpin := gpio.A, gpio.Pin0
	encport, encpins := gpio.A, gpio.Pin1|gpio.Pin5
	ebtnport, ebtnpin := gpio.A, 4

	gpio.B.EnableClock(true)
	ledport, bluepin, greenpin := gpio.B, 6, 7
	lcdport, lcdpins := gpio.B, gpio.Pin10|gpio.Pin11
	wsport, wspin := gpio.B, gpio.Pin13

	gpio.C.EnableClock(true)
	rhport, rhpins := gpio.C, []int{0, 1, 2}
	whport, whpins := gpio.C, gpio.Pin6|gpio.Pin7|gpio.Pin8
	owport, owpin := gpio.C, gpio.Pin10

	// DMA
	dma1 := dma.DMA1
	dma1.EnableClock(true)

	// Button.
	btnport.Setup(btnpin, &gpio.Config{Mode: gpio.In})
	line := exti.Lines(btnpin)
	line.Connect(btnport)
	line.EnableRiseTrig()
	line.EnableInt()
	irqen(irq.EXTI0, 1)

	slowOutCfg := &gpio.Config{Mode: gpio.Out, Speed: gpio.Low}

	// LEDs.
	ledport.SetupPin(bluepin, slowOutCfg)
	ledport.SetupPin(greenpin, slowOutCfg)
	bb := ledport.OutPins()
	ledBlue = bb.Bit(bluepin)
	ledGreen = bb.Bit(greenpin)

	// Room heating.
	bb = rhport.OutPins()
	for i, pin := range rhpins {
		rhport.SetupPin(pin, slowOutCfg)
		RoomHeater[i] = bb.Bit(pin)
	}

	// Water PWM.
	whport.Setup(whpins, &gpio.Config{Mode: gpio.Alt, Speed: gpio.Low})
	whport.SetAltFunc(whpins, gpio.TIM3)
	rcc.RCC.TIM3EN().Set()
	t := tim.TIM3
	setupPulsePWM(t, system.APB1.Clock(), 500, 9999)
	waterPWM.Init(t)
	irqen(irq.TIM3, 12) // Prio must be the same as for water flow sensor.

	// Water flow sensor.
	wsport.Setup(wspin, &gpio.Config{Mode: gpio.AltIn, Pull: gpio.PullUp})
	wsport.SetAltFunc(wspin, gpio.TIM9)
	rcc.RCC.TIM9EN().Set()
	waterCnt.Init(tim.TIM9)
	irqen(irq.TIM9, 12) // Prio must be the same as for PWM IRQ prio.

	// 1-wire bus.
	owport.Setup(owpin, &gpio.Config{Mode: gpio.Alt, Driver: gpio.OpenDrain})
	owport.SetAltFunc(owpin, gpio.USART3)
	owd.Start(usart.USART3, dma1.Channel(3, 0), dma1.Channel(2, 0))
	irqen(irq.USART3, 11)
	irqen(irq.DMA1_Channel3, 11)
	irqen(irq.DMA1_Channel2, 11)

	// I2C LCD (PCF8574T + HD44780)
	lcdport.Setup(lcdpins, &gpio.Config{Mode: gpio.Alt, Driver: gpio.OpenDrain})
	lcdport.SetAltFunc(lcdpins, gpio.I2C2)
	initI2C(i2c.I2C2, dma1.Channel(5, 0), dma1.Channel(4, 0))
	irqen(irq.I2C2_EV, 10)
	irqen(irq.I2C2_ER, 10)
	irqen(irq.DMA1_Channel5, 10)
	irqen(irq.DMA1_Channel4, 10)

	// Encoder.
	encport.Setup(encpins, &gpio.Config{Mode: gpio.AltIn, Pull: gpio.PullUp})
	encport.SetAltFunc(encpins, gpio.TIM2)
	rcc.RCC.TIM2EN().Set()
	ebtnport.SetupPin(ebtnpin, &gpio.Config{Mode: gpio.In, Pull: gpio.PullUp})
	line = exti.LineN(ebtnpin)
	line.Connect(ebtnport)
	encoder.Init(tim.TIM2, ebtnport.InPins().Bit(ebtnpin), line)
	irqen(irq.TIM2, 9)
	irqen(irq.EXTI4, 9)

	rcc.RCC.TIM6EN().Set()
	menu.Setup(tim.TIM6, system.APB1.Clock())
	irqen(irq.TIM6, 5)

	startLCD(i2cdrv.NewMasterConn(0x27, i2c.ASRD))

	initRTC()
	if !checkRTC() {
		fmt.Printf("RTC not set. Setting...\n")
		t := makeDateTime(2016, 9, 11, 22, 51, 20, Sunday)
		setRTC(t)
		fmt.Printf("Done.\n")
	}
}

func main() {
	//go waterTask()
	menu.Loop()
}

func exti0ISR() {
	exti.L0.ClearPending()
	ledGreen.Set()
	delay.Loop(1e5)
	ledGreen.Clear()
	//buttonISR()
	//waterISR()
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.EXTI0: exti0ISR,

	irq.TIM3: waterPWMISR,
	irq.TIM9: waterCntISR,

	irq.USART3:        owdUSARTISR,
	irq.DMA1_Channel3: owdRxDMAISR,
	irq.DMA1_Channel2: owdTxDMAISR,

	irq.I2C2_EV:       i2cISR,
	irq.I2C2_ER:       i2cISR,
	irq.DMA1_Channel5: i2cRxDMAISR,
	irq.DMA1_Channel4: i2cTxDMAISR,

	irq.TIM2:  encoderISR,
	irq.EXTI4: encoderISR,

	irq.TIM6: menuISR,
}
