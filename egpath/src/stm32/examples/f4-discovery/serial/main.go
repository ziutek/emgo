// This example shows how to use USART as serial console.
//
// You need terminal emulator (eg. screen, minicom, hyperterm) and some USB to
// 3.3V TTL serial adapter (eg. FT232RL, CP2102 based). Warninig! If you use
// USB to RS232 it can destroy your Discovery board.
//
// Connct adapter's GND, Rx and Tx pins respectively to Discovery's GND, PA2,
// PA3 pins.
package main

import (
	"delay"
	"fmt"
	"rtos"

	"stm32/serial"

	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/osclk/systick"
	"stm32/hal/system"
	"stm32/hal/usart"
)

const (
	Green  = gpio.Pin12
	Orange = gpio.Pin13
	Red    = gpio.Pin14
	Blue   = gpio.Pin15
)

var (
	leds *gpio.Port
	con  *serial.Dev
)

func init() {
	system.Setup168(8)
	systick.Setup()

	// GPIO

	gpio.A.EnableClock(true)
	port, tx, rx := gpio.A, gpio.Pin2, gpio.Pin3
	gpio.D.EnableClock(false)
	leds = gpio.D

	// LEDS

	cfg := &gpio.Config{Mode: gpio.Out, Speed: gpio.Low}
	leds.Setup(Green|Orange|Red|Blue, cfg)

	// USART

	port.Setup(tx, &gpio.Config{Mode: gpio.Alt})
	port.Setup(rx, &gpio.Config{Mode: gpio.AltIn})
	port.SetAltFunc(tx|rx, gpio.USART2)

	s := usart.USART2

	s.EnableClock(true)
	s.SetBaudRate(115200)
	s.SetConf(usart.RxEna | usart.TxEna)
	s.EnableIRQs(usart.RxNotEmptyIRQ)
	s.Enable()

	con = serial.New(s, 80, 8)
	con.SetUnix(true)
	fmt.DefaultWriter = con

	rtos.IRQ(irq.USART2).Enable()
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

func conISR() {
	// blink(Blue, -10) // Uncoment to see "hardware buffer overrun".
	con.IRQ()
}

//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.USART2: conISR,
}

func checkErr(err error) {
	if err != nil {
		blink(Red, 10)
		con.WriteString("\nError: ")
		con.WriteString(err.Error())
		con.WriteByte('\n')
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

	fmt.Println("\nrtos.Uptime() in loop:")
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
		n, err := con.Read(buf[:])
		checkErr(err)
		ns := rtos.Nanosec()
		fmt.Printf(" %d ns '%s'\n", ns, buf[:n])
		blink(Green, 10)
	}
}
