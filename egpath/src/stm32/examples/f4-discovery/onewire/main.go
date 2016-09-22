// This example shows how to use USART as 1-wire master.
package main

import (
	"delay"
	"fmt"
	"rtos"

	"onewire"

	"stm32/hal/dma"
	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
	"stm32/hal/usart"
	"stm32/onedrv"
)

const Red = gpio.Pin14

var (
	leds *gpio.Port
	one  *usart.Driver
)

func init() {
	system.Setup168(8)
	systick.Setup()

	// GPIO

	gpio.C.EnableClock(true)
	oprt, opin := gpio.C, gpio.Pin6
	gpio.D.EnableClock(false)
	leds = gpio.D

	// DMA
	dma2 := dma.DMA2
	dma2.EnableClock(true) // DMA clock must remain enabled in sleep mode.
	usart6rxdma := dma2.Channel(1, 5)
	usart6txdma := dma2.Channel(7, 5)

	// LEDS

	leds.Setup(Red, gpio.Config{Mode: gpio.Out, Speed: gpio.Low})

	// 1-wire

	oprt.Setup(opin, gpio.Config{Mode: gpio.Alt, Driver: gpio.OpenDrain})
	oprt.SetAltFunc(opin, gpio.USART6)
	one = usart.NewDriver(
		usart.USART6, usart6rxdma, usart6txdma, make([]byte, 16),
	)
	one.EnableClock(true)
	one.SetBaudRate(115200)
	one.SetMode(usart.HalfDuplex)
	one.Enable()
	one.EnableRx()
	one.EnableTx()
	rtos.IRQ(irq.USART6).Enable()
	rtos.IRQ(irq.DMA2_Stream1).Enable()
	rtos.IRQ(irq.DMA2_Stream7).Enable()

}

func oneISR() {
	one.ISR()
}
func oneRxDMAISR() {
	one.RxDMAISR()
}
func oneTxDMAISR() {
	one.TxDMAISR()
}

func blink(pins gpio.Pins, d int) {
	leds.SetPins(pins)
	delay.Millisec(d)
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
	m := onewire.Master{Driver: onedrv.USARTDriver{one}}
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
			s := onewire.MakeSearch(typ, false)
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
		delay.Millisec(2e3)
	}
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.USART6:       oneISR,
	irq.DMA2_Stream1: oneRxDMAISR,
	irq.DMA2_Stream7: oneTxDMAISR,
}
