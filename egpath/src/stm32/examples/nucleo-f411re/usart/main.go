// This example shows how to use USART as serial console.
// It uses PA2, PA3 pins as Tx, Rx that are by default connected to ST-LINK
// Virtual Com Port.
package main

import (
	"bufio"
	"fmt"
	"rtos"
	"text/linewriter"

	"stm32/hal/dma"
	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
	"stm32/hal/usart"
)

var (
	tts      *usart.Driver
	dmarxbuf [88]byte
)

func init() {
	system.Setup96(8)
	systick.Setup(2e6)

	gpio.A.EnableClock(true)
	port, tx, rx := gpio.A, gpio.Pin2, gpio.Pin3

	port.Setup(tx, &gpio.Config{Mode: gpio.Alt})
	port.Setup(rx, &gpio.Config{Mode: gpio.AltIn, Pull: gpio.PullUp})
	port.SetAltFunc(tx|rx, gpio.USART2)
	d := dma.DMA1
	d.EnableClock(true) // DMA clock must remain enabled in sleep mode.
	tts = usart.NewDriver(
		usart.USART2, d.Channel(5, 4), d.Channel(6, 4), dmarxbuf[:],
	)
	tts.P.EnableClock(true)
	tts.P.SetBaudRate(115200)
	tts.P.Enable()
	tts.EnableRx()
	tts.EnableTx()
	rtos.IRQ(irq.USART2).Enable()
	rtos.IRQ(irq.DMA1_Stream5).Enable()
	rtos.IRQ(irq.DMA1_Stream6).Enable()
	fmt.DefaultWriter = linewriter.New(
		bufio.NewWriterSize(tts, 88),
		linewriter.CRLF,
	)
}

func main() {
	fmt.Printf("Enter/paste some text that contains newlines:\n")
	var buf [40]byte
	for {
		n, err := tts.Read(buf[:])
		fmt.DefaultWriter.Write(buf[:n])
		if err != nil {
			fmt.Printf("\n\n*** Error: %v ***\n\n", err)
		}
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
	irq.USART2:       ttsISR,
	irq.DMA1_Stream5: ttsRxDMAISR,
	irq.DMA1_Stream6: ttsTxDMAISR,
}
