// This example shows how to use USART as serial console.
//
// You need terminal emulator (eg. screen, minicom, hyperterm) and some USB to
// 3.3V TTL serial adapter (eg. FT232RL, CP2102 based). Warninig! If you use
// USB to RS232 it can destroy your Discovery board.

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
	"stm32/hal/setup"
	"stm32/hal/usart"
)

const (
	Green = 12 + iota
	Orange
	Red
	Blue
)

var (
	leds *gpio.Port
	con  *serial.Dev
)

func init() {
	setup.Performance168(8)

	// LEDS

	leds = gpio.D

	leds.EnableClock(false)
	leds.SetMode(Green, gpio.Out)
	leds.SetMode(Orange, gpio.Out)
	leds.SetMode(Red, gpio.Out)
	leds.SetMode(Blue, gpio.Out)

	// USART

	port, tx, rx := gpio.A, 2, 3

	port.EnableClock(true)
	port.SetMode(tx, gpio.Alt)
	port.SetAltFunc(tx, gpio.USART2)
	port.SetMode(rx, gpio.Alt)
	port.SetAltFunc(rx, gpio.USART2)

	tts := usart.USART2

	tts.EnableClock(true)
	tts.SetBaudRate(115200)
	tts.SetConf(usart.RxEna | usart.TxEna)
	tts.EnableIRQs(usart.RxNotEmptyIRQ)
	tts.Enable()

	con = serial.New(tts, 80, 8)
	con.SetUnix(true)
	fmt.DefaultWriter = con

	rtos.IRQ(irq.USART2).Enable()
}

func blink(c, dly int) {
	leds.SetPin(c)
	if dly > 0 {
		delay.Millisec(dly)
	} else {
		delay.Loop(-1e4 * dly)
	}
	leds.ClearPin(c)
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
		uts[i] = rtos.Uptime()
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
		ns := rtos.Uptime()
		fmt.Printf(" %d ns '%s'\n", ns, buf[:n])
		blink(Green, 10)
	}
}
