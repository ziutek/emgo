package main

import (
	"delay"

	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/rtc"

	"stm32/hal/raw/i2c"
	"stm32/hal/raw/rcc"
)

var (
	leds *gpio.Port
	twi  *i2c.I2C_Periph
)

const (
	LED1 = gpio.Pin7
	LED2 = gpio.Pin6
)

func init() {
	system.Setup(8, 72/8, false)
	rtc.Setup(32768)

	gpio.B.EnableClock(true)
	leds = gpio.B
	port, pins := gpio.B, gpio.Pin10|gpio.Pin11

	cfg := gpio.Config{Mode: gpio.Out, Speed: gpio.Low}
	leds.Setup(LED1|LED2, &cfg)

	cfg = gpio.Config{
		Mode:   gpio.Alt,
		Driver: gpio.OpenDrain,
	}
	port.Setup(pins, &cfg)

	rcc.RCC.I2C2EN().Set()
	// Mandatory reset.
	rcc.RCC.I2C2RST().Set()
	rcc.RCC.I2C2RST().Clear()
	twi = i2c.I2C2
	freq := system.APB1.Clock() / 1e6 // MHz
	twi.FREQ().Store(i2c.CR2_Bits(freq))
	twi.PE().Clear()
	twi.TRISE.Store(i2c.TRISE_Bits(freq + 1))
	speed := 4400 // Hz
	ccr := system.APB1.Clock() / uint(speed*2)
	if ccr < 4 {
		ccr = 4
	}
	twi.CCR.Store(i2c.CCR_Bits(ccr))
	twi.PE().Set()
}

func main() {
	leds.SetPins(LED1)
	addr := i2c.DR_Bits(0x4e)
	twi.START().Set()
	for twi.SB().Load() == 0 {
	}
	twi.DR.Store(addr)
	for twi.ADDR().Load() == 0 {
	}
	twi.SR2.Load()
	n := 0
	for {
		twi.DR.Store(i2c.DR_Bits(n << 4))
		for twi.BTF().Load() == 0 {
		}
		delay.Millisec(100)
		if n++; n == 16 {
			n = 0
		}
	}
	twi.STOP().Set()
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTCAlarm: rtc.ISR,
}
