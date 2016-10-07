package main

import (
	"bufio"
	"fmt"
	"rtos"
	"text/linewriter"

	"stm32/hal/dma"
	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/usart"
)

var tts *usart.Driver

func initConsole() {
	gpio.A.EnableClock(true)
	port, tx := gpio.A, gpio.Pin2

	port.Setup(tx, gpio.Config{Mode: gpio.Alt})
	port.SetAltFunc(tx, gpio.USART2)
	d := dma.DMA1
	d.EnableClock(true) // DMA clock must remain enabled in sleep mode

	tts = usart.NewDriver(usart.USART2, nil, d.Channel(6, 4), nil)
	tts.P.EnableClock(true)
	tts.P.SetBaudRate(115200)
	tts.P.Enable()
	tts.EnableTx()
	rtos.IRQ(irq.USART2).Enable()
	rtos.IRQ(irq.DMA1_Stream6).Enable()
	fmt.DefaultWriter = linewriter.New(
		bufio.NewWriterSize(tts, 88),
		linewriter.CRLF,
	)
}

func ttsISR() {
	tts.ISR()
}

func ttsTxDMAISR() {
	tts.TxDMAISR()
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.USART2:       ttsISR,
	irq.DMA1_Stream6: ttsTxDMAISR,
}
