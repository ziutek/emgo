package main

import (
	"delay"
	"rtos"
	"strconv"
	"time"

	"stm32/serial"

	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/osclk/rtc"
	"stm32/hal/system"
	"stm32/hal/usart"
)

var (
	leds *gpio.Port
	con  *serial.Dev
)

const (
	LED1 = gpio.Pin7
	LED2 = gpio.Pin6
)

func init() {
	system.Setup(8, 72/8, false)
	rtc.Setup(32768)

	// GPIO

	gpio.A.EnableClock(true)
	port, tx, rx := gpio.A, gpio.Pin9, gpio.Pin10
	gpio.B.EnableClock(false)
	leds = gpio.B

	// LEDs

	cfg := &gpio.Config{Mode: gpio.Out, Speed: gpio.Low}
	leds.Setup(LED1|LED2, cfg)

	// USART

	port.Setup(tx, &gpio.Config{Mode: gpio.Alt})
	port.Setup(rx, &gpio.Config{Mode: gpio.AltIn})

	s := usart.USART1

	s.EnableClock(true)
	s.SetBaudRate(115200)
	s.SetConf(usart.RxEna | usart.TxEna)
	s.EnableIRQs(usart.RxNotEmptyIRQ)
	s.Enable()

	con = serial.New(s, 80, 8)
	con.SetUnix(true)

	rtos.IRQ(irq.USART1).Enable()
}

func conISR() {
	con.IRQ()
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.USART1:   conISR,
	irq.RTCAlarm: rtc.ISR,
}

func blink(led gpio.Pins, dly int) {
	for {
		leds.SetPins(led)
		delay.Millisec(dly)
		leds.ClearPins(led)
		delay.Millisec(dly)
		t := time.Now()
		y, mo, d := t.Date()
		h, mi, s := t.Clock()
		ns := t.Nanosecond()

		// Is ther a easy way to print formated date without fmt package?
		// Yes. Bonus: whole program fits into 40 KB SRAM.

		strconv.WriteInt(con, y, -10, 4)
		con.WriteByte('-')
		strconv.WriteInt(con, int(mo), -10, 2)
		con.WriteByte('-')
		strconv.WriteInt(con, d, -10, 2)
		con.WriteByte(' ')
		strconv.WriteInt(con, h, -10, 2)
		con.WriteByte(':')
		strconv.WriteInt(con, mi, -10, 2)
		con.WriteByte(':')
		strconv.WriteInt(con, s, -10, 2)
		con.WriteByte('.')
		strconv.WriteInt(con, ns, -10, 9)
		con.WriteByte('\n')
	}
}

func main() {
	if ok, set := rtc.Status(); ok && !set {
		rtc.SetTime(time.Date(2016, 1, 24, 22, 58, 30, 0, time.UTC))
	}
	blink(LED2, 500)
}
