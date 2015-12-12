// This example shows how to use USART as 1-wire master.
package main

import (
	"delay"
	"fmt"
	"rtos"

	"stm32/f4/gpio"
	"stm32/f4/irq"
	"stm32/f4/periph"
	"stm32/f4/setup"
	"stm32/f4/usarts"
	"stm32/onedrv"
	"stm32/serial"
	"stm32/usart"

	"onewire"
)

const (
	Green = 12 + iota
	Orange
	Red
	Blue
)

var (
	leds = gpio.D

	ow  = usarts.USART6
	one = serial.New(ow, 8, 8)
)

func init() {
	setup.Performance168(8)

	// LEDS

	periph.AHB1ClockEnable(periph.GPIOD)
	periph.AHB1Reset(periph.GPIOD)

	leds.SetMode(Green, gpio.Out)
	leds.SetMode(Orange, gpio.Out)
	leds.SetMode(Red, gpio.Out)
	leds.SetMode(Blue, gpio.Out)

	// 1-wire

	periph.AHB1ClockEnable(periph.GPIOC)
	periph.AHB1Reset(periph.GPIOC)
	periph.APB2ClockEnable(periph.USART6)
	periph.APB2Reset(periph.USART6)

	port, tx := gpio.C, 6

	port.SetMode(tx, gpio.Alt)
	port.SetOutType(tx, gpio.OpenDrain)
	port.SetAltFunc(tx, gpio.USART6)

	ow.SetWordLen(usart.Bits8)
	ow.SetParity(usart.None)
	ow.SetStopBits(usart.Stop1b)
	ow.SetMode(usart.Tx | usart.Rx)
	ow.SetHalfDuplex(true)
	ow.EnableIRQs(usart.RxNotEmptyIRQ)
	ow.Enable()

	rtos.IRQ(irq.USART6).Enable()
}

func oneISR() {
	one.IRQ()
}

var ISRs = [...]func(){
	irq.USART6: oneISR,
} //c:__attribute__((section(".ISRs")))

func blink(c, d int) {
	leds.SetPin(c)
	if d > 0 {
		delay.Millisec(d)
	} else {
		delay.Loop(-1e4 * d)
	}
	leds.ClearPin(c)
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
	drv := onedrv.UARTDriver{Serial: one, Clock: setup.APB2Clk}
	m := onewire.Master{Driver: &drv}
	dtypes := []onewire.Type{onewire.DS18S20, onewire.DS18B20, onewire.DS1822}

	fmt.Println("\nConfigure all DS18B20, DS1822 to 10bit resolution.")
	checkErr(m.SkipROM())
	checkErr(m.WriteScratchpad(127, -128, onewire.T10bit))

	// This algorithm doesn't work in case of parasite power mode.
	for {
		fmt.Println(
			"\nSending ConvertT command on the bus (SkipROM addressing).",
		)
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
