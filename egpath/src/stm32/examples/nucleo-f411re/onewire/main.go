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
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
	"stm32/hal/usart"
)

const Green = gpio.Pin5

var (
	leds gpio.Port
	con  *serial.Dev
	one  *serial.Dev
)

func init() {
	system.Setup96(8)
	systick.Setup()

	// GPIO
	gpio.A.EnableClock(true)
	leds = gpio.A
	uprt, tx, rx := gpio.A, gpio.Pin2, gpio.Pin3
	gpio.C.EnableClock(true)
	oprt, opin := gpio.C, gpio.Pin6

	// LEDS

	cfg := gpio.Config{Mode: gpio.Out, Speed: gpio.Low}
	leds.Setup(Green, &cfg)

	// USART

	uprt.Setup(tx, &gpio.Config{Mode: gpio.Alt})
	uprt.Setup(rx, &gpio.Config{Mode: gpio.AltIn, Pull: gpio.PullUp})
	uprt.SetAltFunc(tx|rx, gpio.USART2)

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

	// 1-wire

	oprt.Setup(opin, &gpio.Config{Mode: gpio.Alt, Driver: gpio.OpenDrain})
	oprt.SetAltFunc(opin, gpio.USART6)

	s = usart.USART6

	s.EnableClock(true)
	s.SetConf(usart.TxEna | usart.RxEna)
	s.SetMode(usart.HalfDuplex)
	s.EnableIRQs(usart.RxNotEmptyIRQ)
	s.Enable()

	one = serial.New(s, 8, 8)
	rtos.IRQ(irq.USART6).Enable()
}

func conISR() {
	con.IRQ()
}

func oneISR() {
	one.IRQ()
}

func blink(c gpio.Pins, d int) {
	leds.SetPins(c)
	if d > 0 {
		delay.Millisec(d)
	} else {
		delay.Loop(-1e4 * d)
	}
	leds.ClearPins(c)
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

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.USART2: conISR,
	irq.USART6: oneISR,
}
