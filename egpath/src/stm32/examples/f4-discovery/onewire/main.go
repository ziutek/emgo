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
	"stm32/hal/osclk/systick"
	"stm32/hal/system"
	"stm32/hal/usart"
)

const Red = gpio.Pin14

var (
	leds *gpio.Port
	one  *serial.Dev
)

func init() {
	system.Setup168(8)
	systick.Setup()

	// GPIO

	gpio.C.EnableClock(true)
	port, tx := gpio.C, gpio.Pin6
	gpio.D.EnableClock(false)
	leds = gpio.D

	// LEDS

	cfg := &gpio.Config{Mode: gpio.Out, Speed: gpio.Low}
	leds.Setup(Red, cfg)

	// 1-wire

	port.Setup(tx, &gpio.Config{Mode: gpio.Alt, Driver: gpio.OpenDrain})
	port.SetAltFunc(tx, gpio.USART6)

	s := usart.USART6

	s.EnableClock(true)
	s.SetConf(usart.TxEna | usart.RxEna)
	s.SetMode(usart.HalfDuplex)
	s.EnableIRQs(usart.RxNotEmptyIRQ)
	s.Enable()

	one = serial.New(s, 8, 8)
	rtos.IRQ(irq.USART6).Enable()
}

func oneISR() {
	one.IRQ()
}

//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.USART6: oneISR,
}

func blink(pins gpio.Pins, d int) {
	leds.SetPins(pins)
	if d > 0 {
		delay.Millisec(d)
	} else {
		delay.Loop(-1e4 * d)
	}
	leds.ClearPins(pins)
}

func checkErr(err error) {
	if err == nil {
		return
	}
	fmt.Println("Error: ", err)
	for {
		blink(Red, 100)
		delay.Millisec(100)
	}
}

func main() {
	m := onewire.Master{Driver: onedrv.SerialDriver{one}}
	dtypes := []onewire.Type{onewire.DS18S20, onewire.DS18B20, onewire.DS1822}

	delay.Millisec(100)

	fmt.Println("\nConfigure all DS18B20, DS1822 to 10bit resolution.")

	// This algorithm configures and starts conversion simultaneously on all
	// temperature sensors on the bus. It is fast, but doesn't work in case of
	// parasite power mode.

	checkErr(m.SkipROM())
	checkErr(m.WriteScratchpad(127, -128, onewire.T10bit))
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
