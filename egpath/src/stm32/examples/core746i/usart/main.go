// This example shows how to use USART as serial console.
//
// You need terminal emulator (eg. screen, minicom, hyperterm) and some USB to
// 3.3V TTL serial adapter (eg. FT232RL, CP2102 based). Warninig! If you use
// USB to RS232 it can destroy your Core746I board.
//
// Connct adapter's GND, Rx and Tx pins respectively to Discovery's GND, PA2,
// PA3 pins.
package main

import (
	"bufio"
	"delay"
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
	leds     *gpio.Port
	tts      *usart.Driver
	dmarxbuf [88]byte
)

func init() {
	system.Setup192(8)
	systick.Setup(2e6)

	// GPIO

	gpio.A.EnableClock(true)
	port, tx, rx := gpio.A, gpio.Pin2, gpio.Pin3

	// USART

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

func checkErr(err error) {
	if err != nil {
		fmt.Printf("\nError: %v\n", err)
	}
}

func main() {
	var uts [10]int64

	// Following loop disassembled:
	// 0:
	//   svc	5
	//   strd	r0, r1, [r3, #8]!
	//	 ldr	r2, [sp, #52]
	//   cmp	r3, r2
	//   bne.n  0b
	for i := range uts {
		uts[i] = rtos.Nanosec()
	}

	fmt.Println("\nrtos.Nanosec() in loop:")
	for i, ut := range uts {
		fmt.Print(ut, " ns")
		if i > 0 {
			fmt.Printf(" (dt = %d ns)", ut-uts[i-1])
		}
		fmt.Println()
	}

	fmt.Println("Echo:")

	var buf [40]byte
	for {
		n, err := tts.Read(buf[:])
		checkErr(err)
		ns := rtos.Nanosec()
		fmt.Printf(" %d ns '%s'\n", ns, buf[:n])
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
