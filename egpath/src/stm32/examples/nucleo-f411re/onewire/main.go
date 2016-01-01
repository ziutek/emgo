// This example shows how to use USART as 1-wire master.
package main

import (
	"delay"
	"fmt"
	"rtos"

	"onewire"

	"stm32/onedrv"
	"stm32/serial"

	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/setup"
	"stm32/hal/usart"
)

const Green = 5

var (
	LEDport  = gpio.A
	con, one *serial.Dev
)

func init() {
	setup.Performance96(8)

	// LEDS

	LEDport.EnableClock(false)
	LEDport.SetMode(Green, gpio.Out)

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

	// 1-wire

	port, tx = gpio.C, 6

	port.EnableClock(true)
	port.SetMode(tx, gpio.Alt)
	port.SetOutType(tx, gpio.OpenDrain)
	port.SetAltFunc(tx, gpio.USART6)

	ow := usart.USART6

	ow.EnableClock(true)
	ow.SetConf(usart.TxEna | usart.RxEna)
	ow.SetMode(usart.HalfDuplex)
	ow.EnableIRQs(usart.RxNotEmptyIRQ)
	ow.Enable()

	one = serial.New(ow, 8, 8)
	rtos.IRQ(irq.USART6).Enable()
}

func conISR() {
	con.IRQ()
}

func oneISR() {
	one.IRQ()
}

//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.USART2: conISR,
	irq.USART6: oneISR,
}

func blink(c, d int) {
	LEDport.SetPin(c)
	if d > 0 {
		delay.Millisec(d)
	} else {
		delay.Loop(-1e4 * d)
	}
	LEDport.ClearPin(c)
}

func checkErr(err error) {
	if err == nil {
		return
	}
	fmt.Println("Error: ", err)
	for {
		blink(Green, 100)
		delay.Millisec(100)
	}
}

func main() {
	m := onewire.Master{Driver: onedrv.SerialDriver{one}}
	dtypes := []onewire.Type{onewire.DS18S20, onewire.DS18B20, onewire.DS1822}

	fmt.Println("\nConfigure all DS18B20, DS1822 to 10bit resolution.")

	checkErr(m.SkipROM())
	checkErr(m.WriteScratchpad(127, -128, onewire.T10bit))

	// This algorithm doesn't work in case of parasite power mode.
	for {
		fmt.Println("\nSending ConvertT command on the bus (SkipROM addressing).")
		checkErr(m.SkipROM())
		checkErr(m.ConvertT())
		fmt.Print("Waiting until all devices finish the conversion")
		for {
			delay.Millisec(50)
			fmt.Print(".")
			b, err := m.ReadBit()
			checkErr(err)
			if b != 0 {
				break
			}
		}
		fmt.Println("\nSearching for temperature sensors:")
		for _, typ := range dtypes {
			s := onewire.NewSearch(typ, false)
			for m.SearchNext(&s) {
				d := s.Dev()
				fmt.Print(d, ": ")
				checkErr(m.MatchROM(d))
				s, err := m.ReadScratchpad()
				checkErr(err)
				t, err := s.Temp(typ)
				checkErr(err)
				fmt.Printf("%6.2f C\n", t)
			}
			checkErr(s.Err())
		}
		fmt.Println("Done.")
		delay.Millisec(4e3)
	}
}
