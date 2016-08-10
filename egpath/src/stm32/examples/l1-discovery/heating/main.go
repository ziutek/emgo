// Control of power (water heater, house heating system).
package main

import (
	"delay"
	"rtos"

	"arch/cortexm/bitband"

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

func init() {
	system.Setup32(0)
	systick.Setup()

	// GPIO.

	gpio.A.EnableClock(true)
	btnport, btnpin := gpio.A, gpio.Pin0
	wsport, wspin := gpio.A, gpio.Pin15

	gpio.B.EnableClock(true)
	ledport, bluepin, greenpin := gpio.B, 6, 7
	lcdport, lcdpins := gpio.B, gpio.Pin10|gpio.Pin11

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
	rtos.IRQ(irq.EXTI0).Enable()

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
	waterIRQPrio := rtos.IRQPrioHighest - rtos.IRQPrioStep*rtos.IRQPrioNum/4
	e := rtos.IRQ(irq.TIM3)
	e.SetPrio(waterIRQPrio)
	e.Enable()

	// Water flow sensor.
	wsport.Setup(wspin, &gpio.Config{Mode: gpio.AltIn, Pull: gpio.PullUp})
	wsport.SetAltFunc(wspin, gpio.TIM2)
	rcc.RCC.TIM2EN().Set()
	waterCnt.Init(tim.TIM2)
	e = rtos.IRQ(irq.TIM2)
	e.SetPrio(waterIRQPrio)
	e.Enable()

	// 1-wire bus.
	owport.Setup(owpin, &gpio.Config{Mode: gpio.Alt, Driver: gpio.OpenDrain})
	owport.SetAltFunc(owpin, gpio.USART3)
	tempd.Init(usart.USART3, dma1.Channel(3, 0), dma1.Channel(2, 0))
	rtos.IRQ(irq.USART3).Enable()
	rtos.IRQ(irq.DMA1_Channel3).Enable()
	rtos.IRQ(irq.DMA1_Channel2).Enable()

	// I2C LCD (PCF8574T + HD44780)
	lcdport.Setup(lcdpins, &gpio.Config{Mode: gpio.Alt, Driver: gpio.OpenDrain})
	lcdport.SetAltFunc(lcdpins, gpio.I2C2)
	initI2C(i2c.I2C2, dma1.Channel(5, 0), dma1.Channel(4, 0))
	rtos.IRQ(irq.I2C2_EV).Enable()
	rtos.IRQ(irq.I2C2_ER).Enable()
	rtos.IRQ(irq.DMA1_Channel5).Enable()
	rtos.IRQ(irq.DMA1_Channel4).Enable()

	initLCD(i2cdrv.NewMasterConn(0x27, i2c.ASRD))
}

func main() {
	//go waterTask()
	lcd.WriteString("Blaaa!")
	tempd.Loop()
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

	irq.TIM2: waterCntISR,
	irq.TIM3: waterPWMISR,

	irq.USART3:        tempdUSARTISR,
	irq.DMA1_Channel3: tempdRxDMAISR,
	irq.DMA1_Channel2: tempdTxDMAISR,

	irq.I2C2_EV:       i2cISR,
	irq.I2C2_ER:       i2cISR,
	irq.DMA1_Channel5: i2cRxDMAISR,
	irq.DMA1_Channel4: i2cTxDMAISR,
}
