// This example shows how to use button.PollDrv to handle input from multiple
// buttons. Additionally it demonstrates differences between gpio.Pin and
// gpio.Pins.
//
// LED1 and LED2 shows the state of KEY1 and KEY2. At the same time semihosting
// standard output is used to print input events. Use ../semihosting.sh or
//../debug-bmp.sh or to see them.
//
// LineWriter demonstrates the adventages of Go interfaces. It "fixes"
// bufio.Writer to speed up the semihosting output.
package main

import (
	"bits"
	"bufio"
	"bytes"
	"debug/semihosting"
	"fmt"

	"nrf5/input"
	"nrf5/input/button"

	"nrf5/hal/clock"
	"nrf5/hal/gpio"
	"nrf5/hal/irq"
	"nrf5/hal/rtc"
	"nrf5/hal/system"
	"nrf5/hal/system/timer/rtcst"
)

var (
	inputCh    = make(chan input.Event, 3)
	btdrv      *button.PollDrv
	key1, key2 gpio.Pins
	led1, led2 gpio.Pin
)

func init() {
	// Initialize system and runtime.
	system.Setup(clock.XTAL, clock.XTAL, true)
	rtcst.Setup(rtc.RTC1, 0)

	// Allocate pins (always do it in one place to avoid conflicts).

	p0 := gpio.P0
	key1 = gpio.Pin16
	key2 = gpio.Pin17
	led1 = p0.Pin(19)
	led2 = p0.Pin(20)

	// Configure peripherals.

	led1.Setup(gpio.ModeOut)
	led2.Setup(gpio.ModeOut)
	btdrv = button.NewPollDrv(p0, key1|key2, true, inputCh, 0)
	btdrv.UseRTC(rtc.RTC1, 1, 20)

	// Semihosting console.

	f, err := semihosting.OpenFile(":tt", semihosting.W)
	for err != nil {
	}
	fmt.DefaultWriter = lineWriter{bufio.NewWriterSize(f, 40)}
}

func main() {
	for ev := range inputCh {
		pins := ev.Pins()
		led1.Store(bits.One(pins&key1 == 0))
		led2.Store(bits.One(pins&key2 == 0))
		fmt.Printf("src=%d val=%08x\n", ev.Src(), pins)
	}
}

func rtcISR() {
	rtcst.ISR()
	btdrv.RTCISR()
}

//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTC1: rtcISR,
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
