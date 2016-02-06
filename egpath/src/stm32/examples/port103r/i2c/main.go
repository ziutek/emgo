package main

import (
	"delay"
	"fmt"

	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/rtc"

	"stm32/hal/i2c"
)

var (
	leds *gpio.Port
	twi  = &i2c.Driver{Periph: i2c.I2C2}
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
}

func twiISR() {
	twi.ISR()
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTCAlarm: rtc.ISR,
	irq.I2C2_EV:  twiISR,
	irq.I2C2_ER:  twiISR,
}

func main() {
	delay.Millisec(5)

	leds.SetPins(LED1)

	twi.EnableClock(true)
configure:
	twi.Reset() // Mandatory!
	twi.Setup(&i2c.Config{Speed: 4400})

	c := twi.MasterConn(0x27)

loop:
	_, err := c.Write([]byte{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0xff,
	})
	if err != nil {
		if err.(i2c.Error)&i2c.SoftTimeout != 0 {
			fmt.Printf("SoftTimeout\n")
			goto configure
		} else {
			fmt.Printf("0x%02x\n", err)
			twi.SoftReset()
			goto loop
		}
	}
	c.Stop()
	goto loop
}
