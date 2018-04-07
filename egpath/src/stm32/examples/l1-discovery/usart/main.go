// This example shows how to use USART as serial console.
//
// You need terminal emulator (eg. screen, minicom, hyperterm) and some USB to
// 3.3V TTL serial adapter (eg. FT232RL, CP2102 based). Warninig! If you use
// USB to RS232 it can destroy your Discovery board.
//
// Connct adapter's GND, Rx and Tx pins respectively to Discovery's GND, PB10,
// PB11 pins.
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

const (
	Blue  = gpio.Pin6
	Green = gpio.Pin7
)

var (
	leds     *gpio.Port
	tts      *usart.Driver
	dmarxbuf [88]byte
)

func init() {
	system.Setup32(0)
	systick.Setup(2e6)

	// GPIO

	gpio.B.EnableClock(true)
	port, tx, rx := gpio.B, gpio.Pin10, gpio.Pin11
	leds = gpio.B

	// LEDs

	cfg := &gpio.Config{Mode: gpio.Out, Speed: gpio.Low}
	leds.Setup(Green|Blue, cfg)

	// USART

	port.Setup(tx, &gpio.Config{Mode: gpio.Alt})
	port.Setup(rx, &gpio.Config{Mode: gpio.AltIn, Pull: gpio.PullUp})
	port.SetAltFunc(tx|rx, gpio.USART3)
	d := dma.DMA1
	d.EnableClock(true) // DMA clock must remain enabled in sleep mode.
	tts = usart.NewDriver(
		usart.USART3, d.Channel(2, 0), d.Channel(3, 0), dmarxbuf[:],
	)
	tts.Periph().EnableClock(true)
	tts.Periph().SetBaudRate(115200)
	tts.Periph().Enable()
	tts.EnableRx()
	tts.EnableTx()
	rtos.IRQ(irq.USART3).Enable()
	rtos.IRQ(irq.DMA1_Channel3).Enable()
	rtos.IRQ(irq.DMA1_Channel2).Enable()
	fmt.DefaultWriter = linewriter.New(
		bufio.NewWriterSize(tts, 88),
		linewriter.CRLF,
	)
}

func blink(c gpio.Pins, dly int) {
	leds.SetPins(c)
	if dly > 0 {
		delay.Millisec(dly)
	} else {
		delay.Loop(-1e4 * dly)
	}
	leds.ClearPins(c)
}

func checkErr(err error) {
	if err != nil {
		fmt.Printf("\nError: %v\n", err)
		blink(Blue, 10)
	}
}

func main() {
	var uts [10]int64

	// Following loop disassembled:
	// 0b:
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
		blink(Green, 10)
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
	irq.USART3:        ttsISR,
	irq.DMA1_Channel3: ttsRxDMAISR,
	irq.DMA1_Channel2: ttsTxDMAISR,
}
