// Reflection test. It uses USART2 configured as Tx only.
package main

import (
	"bufio"
	"fmt"
	"reflect"
	"rtos"
	"text/linewriter"

	"stm32/hal/dma"
	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
	"stm32/hal/usart"
)

var tts *usart.Driver

func init() {
	system.Setup96(8)
	systick.Setup(2e6)

	gpio.A.EnableClock(true)
	port, tx := gpio.A, gpio.Pin2

	port.Setup(tx, &gpio.Config{Mode: gpio.Alt})
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

type S struct {
	s string
}

func (s S) String() string {
	return s.s
}

type T struct {
	a int
	b byte
	S
}

type Stringer interface {
	String() string
}

var ivals = [...]interface{}{
	1,
	byte(2),
	uintptr(3),
	rtos.Debug(4),
	T{5, 6, S{"foo"}},
	nil,
}

func main() {
	fmt.Println("\nReflection test:\n")
	for _, iv := range ivals {
		v := reflect.ValueOf(iv)
		fmt.Println(v.Type())
	}
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
