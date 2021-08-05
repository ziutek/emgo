// Control of power (water heater, house heating system).
package main

import (
	"rtos"

	"arch/cortexm/nvic"

	"stm32/hal/dma"
	"stm32/hal/exti"
	"stm32/hal/gpio"
	"stm32/hal/i2c"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
	"stm32/hal/usart"

	"stm32/hal/raw/rcc"
	"stm32/hal/raw/tim"
)

// Onboard LED (green is diconected to use PB7 for room heater).
var ledBlue gpio.Pin

// prio16 must be in the range [0;15]. Do not use prio16 <= 8 (SVCall prio) for
// realtime ISRs.
func irqen(irq nvic.IRQ, prio16 rtos.IRQPrio) {
	e := rtos.IRQ(irq)
	e.SetPrio(rtos.IRQPrioLowest + prio16*rtos.IRQPrioStep*(rtos.IRQPrioNum/16))
	e.Enable()
}

func init() {
	system.Setup32(0) // use 8MHz HSI as clock source
	systick.Setup(2e6)

	rcc.RCC.HSITRIM().Store((16 + 1) << rcc.HSITRIMn) // trim HSI clock around 1*40 kHz up

	// GPIO.

	gpio.A.EnableClock(true)
	btn := gpio.A.Pin(0)
	encport, encpins := gpio.A, gpio.Pin1|gpio.Pin5
	encbtn := gpio.A.Pin(4)

	gpio.B.EnableClock(true)
	ledBlue = gpio.B.Pin(6)
	rhport, rhpins := gpio.B, gpio.Pin7|gpio.Pin8|gpio.Pin9
	lcdport, lcdpins := gpio.B, gpio.Pin10|gpio.Pin11
	wsport, wspin := gpio.B, gpio.Pin13

	gpio.C.EnableClock(true)
	whport, whpins := gpio.C, gpio.Pin6|gpio.Pin7|gpio.Pin8
	owport, owpin := gpio.C, gpio.Pin10

	// DMA
	dma1 := dma.DMA1
	dma1.EnableClock(true)

	// Button.
	btn.Setup(&gpio.Config{Mode: gpio.In})
	line := exti.LineIndex(btn.Index())
	line.Connect(btn.Port())
	line.EnableRiseTrig()
	line.EnableIRQ()
	irqen(irq.EXTI0, 1)

	slowOutCfg := gpio.Config{Mode: gpio.Out, Speed: gpio.Low}

	// LED.
	ledBlue.Setup(&slowOutCfg)

	// Room heating PWM.
	rhport.Setup(rhpins, &gpio.Config{Mode: gpio.Alt, Speed: gpio.Low})
	rhport.SetAltFunc(rhpins, gpio.TIM4)
	rcc.RCC.TIM4EN().Set()

	// Water heating PWM.
	whport.Setup(whpins, &gpio.Config{Mode: gpio.Alt, Speed: gpio.Low})
	whport.SetAltFunc(whpins, gpio.TIM3)
	rcc.RCC.TIM3EN().Set()
	irqen(irq.TIM3, 13) // Prio must be the same as for water flow sensor.

	// Water flow sensor.
	wsport.Setup(wspin, &gpio.Config{Mode: gpio.AltIn, Pull: gpio.PullUp})
	wsport.SetAltFunc(wspin, gpio.TIM9)
	rcc.RCC.TIM9EN().Set()
	irqen(irq.TIM9, 13) // Prio must be the same as for PWM IRQ prio.

	water.Init(tim.TIM3, tim.TIM9, system.APB1.Clock())

	// 1-wire bus.
	owport.Setup(owpin, &gpio.Config{Mode: gpio.Alt, Driver: gpio.OpenDrain})
	owport.SetAltFunc(owpin, gpio.USART3)
	owd.Start(usart.USART3, dma1.Channel(2, 0), dma1.Channel(3, 0))
	irqen(irq.USART3, 12)
	irqen(irq.DMA1_Channel3, 12)
	irqen(irq.DMA1_Channel2, 12)

	// I2C LCD (PCF8574T + HD44780)
	lcdport.Setup(lcdpins, &gpio.Config{Mode: gpio.Alt, Driver: gpio.OpenDrain})
	lcdport.SetAltFunc(lcdpins, gpio.I2C2)
	initI2C(i2c.I2C2, dma1.Channel(5, 0), dma1.Channel(4, 0))
	irqen(irq.I2C2_EV, 11)
	irqen(irq.I2C2_ER, 11)
	irqen(irq.DMA1_Channel5, 11)
	irqen(irq.DMA1_Channel4, 11)

	// Encoder.
	encport.Setup(encpins, &gpio.Config{Mode: gpio.AltIn, Pull: gpio.PullUp})
	encport.SetAltFunc(encpins, gpio.TIM2)
	rcc.RCC.TIM2EN().Set()
	encbtn.Setup(&gpio.Config{Mode: gpio.In, Pull: gpio.PullUp})
	encoder.Init(tim.TIM2, encbtn)
	irqen(irq.TIM2, 10)
	irqen(irq.EXTI4, 10)

	rcc.RCC.TIM6EN().Set()
	menu.Setup(tim.TIM6, system.APB1.Clock())
	irqen(irq.TIM6, 9)

	//initRTC()

	// room.Start must be after owd.Start.
	room.Start(tim.TIM4, system.APB1.Clock())

	// startLCD must be last to allow work without LCD.
	startLCD(i2cdrv.NewMasterConn(0x27, i2c.ASRD))
}

func main() {
	menu.Loop()
}

func exti0ISR() {
	exti.L0.ClearPending()
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
