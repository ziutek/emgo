package main

import (
	"io"
	"rtos"

	"stm32/hal/dma"
	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
	"stm32/hal/usart"
)

var tts *usart.Driver

func init() {
	system.SetupPLL(8, 1, 48/8)
	systick.Setup(2e6)

	// GPIO

	gpio.A.EnableClock(true)
	port, tx := gpio.A, gpio.Pin9

	// USART

	port.Setup(tx, &gpio.Config{Mode: gpio.Alt})
	port.SetAltFunc(tx, gpio.USART1_AF1)
	d := dma.DMA1
	d.EnableClock(true) // DMA clock must remain enabled in sleep mode.
	tts = usart.NewDriver(usart.USART1, nil, d.Channel(2, 0), nil)
	tts.P.EnableClock(true) // USART clock must remain enabled in sleep mode.
	tts.P.SetBaudRate(115200)
	tts.P.Enable()
	tts.EnableTx()

	rtos.IRQ(irq.USART1).Enable()
	rtos.IRQ(irq.DMA1_Channel2_3).Enable()
}

func main() {
	io.WriteString(tts, "\r\nHello, World!\r\n")
}

func ttsISR() {
	tts.ISR()
}

func ttsDMAISR() {
	tts.TxDMAISR()
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.USART1:          ttsISR,
	irq.DMA1_Channel2_3: ttsDMAISR,
}
