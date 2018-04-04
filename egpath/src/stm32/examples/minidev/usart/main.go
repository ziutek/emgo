package main

import (
	"delay"
	"rtos"

	"stm32/hal/dma"
	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/rtcst"
	"stm32/hal/usart"
)

var (
	led gpio.Pin
	tts *usart.Driver
)

func init() {
	system.SetupPLL(8, 1, 72/8)
	rtcst.Setup(32768)

	// GPIO

	gpio.B.EnableClock(true)
	port, tx, rx := gpio.B, gpio.Pin10, gpio.Pin11

	gpio.C.EnableClock(true)
	led = gpio.C.Pin(13)

	// LEDs

	cfg := gpio.Config{Mode: gpio.Out, Driver: gpio.OpenDrain, Speed: gpio.Low}
	led.Setup(&cfg)

	// USART

	port.Setup(tx, &gpio.Config{Mode: gpio.Alt})
	port.Setup(rx, &gpio.Config{Mode: gpio.AltIn, Pull: gpio.PullUp})
	d := dma.DMA1
	d.EnableClock(true) // DMA clock must remain enabled in sleep mode.
	tts = usart.NewDriver(
		usart.USART3, d.Channel(3, 0), d.Channel(2, 0), make([]byte, 40),
	)
	tts.P.EnableClock(true)
	tts.P.SetBaudRate(115200)
	tts.P.Enable()
	tts.EnableRx()
	tts.EnableTx()

	rtos.IRQ(irq.USART3).Enable()
	rtos.IRQ(irq.DMA1_Channel2).Enable()
	rtos.IRQ(irq.DMA1_Channel3).Enable()
}

func main() {
	tts.WriteString("\r\nEcho:\r\n")
	var buf [40]byte
	for {
		n, err := tts.Read(buf[:])
		if err != nil {
			tts.WriteString(err.Error())
			tts.WriteString(" -> ")
		} else {
			tts.WriteString("ok -> ")
		}
		tts.Write(buf[:n])
		tts.WriteString("\r\n")
		led.Set()
		delay.Millisec(50)
		led.Clear()
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

//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTCAlarm: rtcst.ISR,

	irq.USART3:        ttsISR,
	irq.DMA1_Channel2: ttsTxDMAISR,
	irq.DMA1_Channel3: ttsRxDMAISR,
}
