// This example shows how to use encoder package to handle rotary encoder. It
// uses semihosting to print encoder events. Use Black Magic Probe
// (../debug-bmp.sh) or OpenOCD (../semihosting.sh) to see program output.
package main

import (
	"bufio"
	"bytes"
	"debug/semihosting"
	"fmt"
	"rtos"

	"nrf5/encoder"

	"nrf5/hal/clock"
	"nrf5/hal/gpio"
	"nrf5/hal/gpiote"
	"nrf5/hal/irq"
	"nrf5/hal/rtc"
	"nrf5/hal/system"
	"nrf5/hal/system/timer/rtcst"
)

var enc *encoder.Driver

func init() {
	// Initialize system and runtime.
	system.Setup(clock.XTAL, clock.XTAL, true)
	rtcst.Setup(rtc.RTC0, 1)

	// Allocate pins (always do it in one place to avoid conflicts).

	p0 := gpio.P0
	btn := p0.Pin(4)
	a := p0.Pin(5)
	b := p0.Pin(7)

	// Configure peripherals.

	enc = encoder.New(a, b, btn, gpiote.Chan(0), true)

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
	for ev := range enc.Events() {
		fmt.Printf("offset=%d button=%t\n", ev.Offset(), ev.Button())
	}
}

func qdecISR() {
	enc.QDECISR()
}

func gpioteISR() {
	enc.GPIOTEISR()
}

//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTC0:   rtcst.ISR,
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
