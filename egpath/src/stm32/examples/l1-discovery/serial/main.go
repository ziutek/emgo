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
	"delay"
	"fmt"
	"rtos"

	"stm32/serial"

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
	leds *gpio.Port
	con  *serial.Dev
)

func init() {
	system.Setup32(0)
	systick.Setup()

	gpio.B.EnableClock(true)
	leds = gpio.B
	port, tx, rx := gpio.B, gpio.Pin10, gpio.Pin11

	// LEDS

	cfg := gpio.Config{Mode: gpio.Out, Speed: gpio.Low}
	leds.Setup(Green|Blue, &cfg)

	// USART

	port.Setup(tx, &gpio.Config{Mode: gpio.Alt})
	port.Setup(rx, &gpio.Config{Mode: gpio.AltIn, Pull: gpio.PullUp})
	port.SetAltFunc(tx|rx, gpio.USART3)

	s := usart.USART3

	s.EnableClock(true)
	s.SetBaudRate(115200)
	s.SetConf(usart.RxEna | usart.TxEna)
	s.EnableIRQs(usart.RxNotEmptyIRQ)
	s.Enable()

	con = serial.New(s, 80, 8)
	con.SetUnix(true)
	fmt.DefaultWriter = con

	rtos.IRQ(irq.USART3).Enable()
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
	//blink(Blue, -10) // Uncoment to see "hardware buffer overrun".
	con.IRQ()
}

func checkErr(err error) {
	if err != nil {
		blink(Blue, 10)
		con.WriteString("\nError: ")
		con.WriteString(err.Error())
		con.WriteByte('\n')
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
		n, err := con.Read(buf[:])
		checkErr(err)
		ns := rtos.Nanosec()
		fmt.Printf(" %d ns '%s'\n", ns, buf[:n])
		blink(Green, 10)
	}
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.USART3: conISR,
}
