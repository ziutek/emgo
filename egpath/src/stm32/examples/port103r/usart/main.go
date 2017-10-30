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
	tts      *usart.Driver
	dmarxbuf [88]byte

	led1 = gpio.B.Pin(7)
	led2 = gpio.B.Pin(6)
	led3 = gpio.B.Pin(5)
	led4 = gpio.D.Pin(2)
)

func init() {
	system.SetupPLL(8, 1, 72/8)
	rtcst.Setup(32768)

	// GPIO

	gpio.A.EnableClock(true)
	port, tx, rx := gpio.A, gpio.Pin9, gpio.Pin10
	gpio.B.EnableClock(false)
	gpio.D.EnableClock(false)

	// LEDs

	cfg := gpio.Config{Mode: gpio.Out, Speed: gpio.Low}
	led1.Setup(&cfg)
	led2.Setup(&cfg)
	led3.Setup(&cfg)
	led4.Setup(&cfg)

	// USART

	port.Setup(tx, &gpio.Config{Mode: gpio.Alt})
	port.Setup(rx, &gpio.Config{Mode: gpio.AltIn, Pull: gpio.PullUp})
	d := dma.DMA1
	d.EnableClock(true) // DMA clock must remain enabled in sleep.
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

func main() {
	/*
		tts.WriteString("\r\nAFIO.MAPR: ")
		strconv.WriteUint32(tts, afio.AFIO.MAPR.U32.Load(), -16, 0)
		tts.WriteString("\r\n")
	*/

	if ok, set := rtcst.Status(); ok && !set {
		rtcst.SetTime(time.Date(2016, 1, 24, 22, 58, 30, 0, time.UTC))
	}
	for {
		led4.Set()
		delay.Millisec(500)
		led4.Clear()
		delay.Millisec(500)

		t := time.Now()
		y, mo, d := t.Date()
		h, mi, s := t.Clock()
		ns := t.Nanosecond()

		// Is there an easy way to print formated date without fmt package? Yes,
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
