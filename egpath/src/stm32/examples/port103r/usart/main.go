package main

import (
	"delay"
	"rtos"
	"strconv"
	"time"

	"stm32/hal/dma"
	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/rtcst"
	"stm32/hal/usart"
)

var (
	leds     *gpio.Port
	tts      *usart.Driver
	dmarxbuf [88]byte
)

const (
	LED1 = gpio.Pin7
	LED2 = gpio.Pin6
)

func init() {
	system.Setup(8, 1, 72/8)
	rtcst.Setup(32768)

	// GPIO

	gpio.A.EnableClock(true)
	port, tx, rx := gpio.A, gpio.Pin9, gpio.Pin10
	gpio.B.EnableClock(false)
	leds = gpio.B

	// LEDs

	cfg := gpio.Config{Mode: gpio.Out, Speed: gpio.Low}
	leds.Setup(LED1|LED2, &cfg)

	// USART

	port.Setup(tx, &gpio.Config{Mode: gpio.Alt})
	port.Setup(rx, &gpio.Config{Mode: gpio.AltIn, Pull: gpio.PullUp})
	d := dma.DMA1
	d.EnableClock(true) // DMA clock must remain enabled in s
	tts = usart.NewDriver(
		usart.USART1, d.Channel(5, 0), d.Channel(4, 0), dmarxbuf[:],
	)
	tts.P.EnableClock(true)
	tts.P.SetBaudRate(115200)
	tts.P.Enable()
	tts.EnableRx()
	tts.EnableTx()
	rtos.IRQ(irq.USART1).Enable()
	rtos.IRQ(irq.DMA1_Channel5).Enable()
	rtos.IRQ(irq.DMA1_Channel4).Enable()
}

func printDate(led gpio.Pins, dly int) {
	for {
		leds.SetPins(led)
		delay.Millisec(dly)
		leds.ClearPins(led)
		delay.Millisec(dly)
		t := time.Now()
		y, mo, d := t.Date()
		h, mi, s := t.Clock()
		ns := t.Nanosecond()

		// Is there a easy way to print formated date without fmt package? Yes,
		// and by avoiding fmt package the whole program fits into 48 KB SRAM.

		strconv.WriteInt(tts, y, -10, 4)
		tts.WriteByte('-')
		strconv.WriteInt(tts, int(mo), -10, 2)
		tts.WriteByte('-')
		strconv.WriteInt(tts, d, -10, 2)
		tts.WriteByte(' ')
		strconv.WriteInt(tts, h, -10, 2)
		tts.WriteByte(':')
		strconv.WriteInt(tts, mi, -10, 2)
		tts.WriteByte(':')
		strconv.WriteInt(tts, s, -10, 2)
		tts.WriteByte('.')
		strconv.WriteInt(tts, ns, -10, 9)
		tts.WriteString("\r\n")
	}
}

func main() {
	if ok, set := rtcst.Status(); ok && !set {
		rtcst.SetTime(time.Date(2016, 1, 24, 22, 58, 30, 0, time.UTC))
	}
	printDate(LED2, 500)
}

func ttsISR() {
	tts.ISR()
}

func ttsRxDMAISR() {
	tts.RxDMAISR()
}

func ttsTxDMAISR() {
	tts.TxDMAISR()
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTCAlarm:      rtcst.ISR,
	irq.USART1:        ttsISR,
	irq.DMA1_Channel5: ttsRxDMAISR,
	irq.DMA1_Channel4: ttsTxDMAISR,
}
