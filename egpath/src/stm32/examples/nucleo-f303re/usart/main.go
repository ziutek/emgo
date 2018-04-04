package main

import (
	"delay"
	"rtos"

	"stm32/hal/dma"
	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
	"stm32/hal/usart"
)

var (
	led gpio.Pin
	tts *usart.Driver
)

func init() {
	system.SetupPLL(8, 1, 72/8)
	systick.Setup(2e6)

	// GPIO

	gpio.A.EnableClock(true)
	port, tx, rx := gpio.A, gpio.Pin2, gpio.Pin3
	led = gpio.A.Pin(5)

	// LEDs

	cfg := gpio.Config{Mode: gpio.Out, Speed: gpio.Low}
	led.Setup(&cfg)

	// USART

	port.Setup(tx, &gpio.Config{Mode: gpio.Alt})
	port.Setup(rx, &gpio.Config{Mode: gpio.AltIn, Pull: gpio.PullUp})
	port.SetAltFunc(tx|rx, gpio.USART2)
	d := dma.DMA1
	d.EnableClock(true) // DMA clock must remain enabled in sleep mode.
	tts = usart.NewDriver(
		usart.USART2, d.Channel(6, 0), d.Channel(7, 0), make([]byte, 40),
	)
	tts.P.EnableClock(true)
	tts.P.SetBaudRate(115200)
	tts.P.Enable()
	tts.EnableRx()
	tts.EnableTx()

	rtos.IRQ(irq.USART2).Enable()
	rtos.IRQ(irq.DMA1_Channel6).Enable()
	rtos.IRQ(irq.DMA1_Channel7).Enable()
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

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.USART2:        ttsISR,
	irq.DMA1_Channel6: ttsRxDMAISR,
	irq.DMA1_Channel7: ttsTxDMAISR,
}
