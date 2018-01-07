// This example shows how to use input/encoder and input/button packages to
// handle rotary encoder. It uses semihosting to print encoder events. Use
// Black Magic Probe (../debug-bmp.sh) or OpenOCD (../semihosting.sh) to see
// program output.
package main

import (
	"bufio"
	"bytes"
	"debug/semihosting"
	"fmt"
	"rtos"

	"nrf5/input"
	"nrf5/input/button"
	"nrf5/input/encoder"

	"nrf5/hal/clock"
	"nrf5/hal/gpio"
	"nrf5/hal/gpiote"
	"nrf5/hal/irq"
	"nrf5/hal/rtc"
	"nrf5/hal/system"
	"nrf5/hal/system/timer/rtcst"
)

const (
	Enc = iota
	Btn
)

var (
	inpCh = make(chan input.Event, 4)
	enc   *encoder.Driver
	btn   *button.Driver
)

func init() {
	// Initialize system and runtime.
	system.Setup(clock.XTAL, clock.XTAL, true)
	rtcst.Setup(rtc.RTC1, 0)

	// Allocate pins (always do it in one place to avoid conflicts).

	p0 := gpio.P0
	bt := p0.Pin(4)
	a := p0.Pin(5)
	b := p0.Pin(7)

	// Configure peripherals.

	enc = encoder.New(a, b, true, true, inpCh, Enc)
	btn = button.New(bt, gpiote.Chan(0), true, rtc.RTC1, 1, inpCh, Btn)

	// Configure interrupts.

	rtos.IRQ(irq.QDEC).Enable()
	rtos.IRQ(irq.GPIOTE).Enable()

	// Semihosting console.

	f, err := semihosting.OpenFile(":tt", semihosting.W)
	for err != nil {
	}
	fmt.DefaultWriter = lineWriter{bufio.NewWriterSize(f, 40)}
}

func main() {
	for ev := range inpCh {
		fmt.Printf("src=%d val=%d\n", ev.Src(), ev.Val())
	}
}

func qdecISR() {
	enc.ISR()
}

func gpioteISR() {
	btn.ISR()
}

func rtcISR() {
	rtcst.ISR()
	btn.RTCISR()
}

//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTC1:   rtcISR,
	irq.QDEC:   qdecISR,
	irq.GPIOTE: gpioteISR,
}

type lineWriter struct {
	w *bufio.Writer
}

func (b lineWriter) Write(s []byte) (int, error) {
	n, err := b.w.Write(s)
	if err != nil {
		return n, err
	}
	if bytes.IndexByte(s, '\n') >= 0 {
		err = b.w.Flush()
	}
	return n, err
}
