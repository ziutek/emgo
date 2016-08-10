// This example shows how to use USART as 1-wire master.
package main

import (
	"bufio"
	"delay"
	"fmt"
	"rtos"
	"text/linewriter"

	"onewire"

	"stm32/hal/dma"
	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
	"stm32/hal/usart"
	"stm32/onedrv"
)

const Green = gpio.Pin5

var (
	leds *gpio.Port
	con  *usart.Driver
	one  *usart.Driver
)

func init() {
	system.Setup96(8)
	systick.Setup()

	// GPIO
	gpio.A.EnableClock(true)
	leds = gpio.A
	cprt, tx, rx := gpio.A, gpio.Pin2, gpio.Pin3
	gpio.C.EnableClock(true)
	oprt, opin := gpio.C, gpio.Pin6

	// DMA
	dma1 := dma.DMA1
	dma1.EnableClock(true) // DMA clock must remain enabled in sleep mode.
	usart2rxdma := dma1.Channel(5, 4)
	usart2txdma := dma1.Channel(6, 4)

	dma2 := dma.DMA2
	dma2.EnableClock(true) // DMA clock must remain enabled in sleep mode.
	usart6rxdma := dma2.Channel(1, 5)
	usart6txdma := dma2.Channel(7, 5)

	// LEDS

	cfg := gpio.Config{Mode: gpio.Out, Speed: gpio.Low}
	leds.Setup(Green, &cfg)

	// Console

	cprt.Setup(tx, &gpio.Config{Mode: gpio.Alt})
	cprt.Setup(rx, &gpio.Config{Mode: gpio.AltIn, Pull: gpio.PullUp})
	cprt.SetAltFunc(tx|rx, gpio.USART2)
	con = usart.NewDriver(
		usart.USART2, usart2rxdma, usart2txdma, make([]byte, 80),
	)
	con.EnableClock(true)
	con.SetBaudRate(115200)
	con.Enable()
	con.EnableRx()
	con.EnableTx()
	rtos.IRQ(irq.USART2).Enable()
	rtos.IRQ(irq.DMA1_Stream5).Enable()
	rtos.IRQ(irq.DMA1_Stream6).Enable()
	fmt.DefaultWriter = linewriter.New(
		bufio.NewWriterSize(con, 88),
		linewriter.CRLF,
	)

	// 1-wire

	oprt.Setup(opin, &gpio.Config{Mode: gpio.Alt, Driver: gpio.OpenDrain})
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

func conISR() {
	con.ISR()
}
func conRxDMAISR() {
	con.RxDMAISR()
}
func conTxDMAISR() {
	con.TxDMAISR()
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

func blink(c gpio.Pins, d int) {
	leds.SetPins(c)
	delay.Millisec(d)
	leds.ClearPins(c)
}

func printErr(err error) bool {
	if err == nil {
		return false
	}
	fmt.Printf("Error: %v\n", err)
	for i := 0; i < 5; i++ {
		blink(Green, 100)
		delay.Millisec(100)
	}
	return true
}

func printOK() {
	fmt.Println("OK.")
}

func main() {
	m := onewire.Master{Driver: onedrv.USARTDriver{one}}
	dtypes := []onewire.Type{onewire.DS18S20, onewire.DS18B20, onewire.DS1822}

	// This algorithm configures and starts conversion simultaneously on all
	// temperature sensors on the bus. It is fast, but doesn't work in case of
	// parasite power mode.


start:
	for {
		fmt.Print("\nConfigure all DS18B20, DS1822 to 10bit resolution: ")
		if printErr(m.SkipROM()) {
			continue start
		}
		if printErr(m.WriteScratchpad(127, -128, onewire.T10bit)) {
			continue start
		}
		printOK()

		fmt.Print("Sending ConvertT command (SkipROM addressing): ")
		if printErr(m.SkipROM()) {
			continue start
		}
		if printErr(m.ConvertT()) {
			continue start
		}
		printOK()

		fmt.Print("Waiting until all devices finish the conversion: ")
		for {
			delay.Millisec(50)
			b, err := m.ReadBit()
			if printErr(err) {
				continue start
			}
			fmt.Print(". ")
			if b != 0 {
				printOK()
				break
			}
		}
		fmt.Print("Searching for temperature sensors: ")
		for _, typ := range dtypes {
			s := onewire.MakeSearch(typ, false)
			for m.SearchNext(&s) {
				d := s.Dev()
				fmt.Printf("\n %v : ", d)
				if printErr(m.MatchROM(d)) {
					continue start
				}
				s, err := m.ReadScratchpad()
				if printErr(err) {
					continue start
				}
				t, err := s.Temp(typ)
				if printErr(err) {
					continue start
				}
				fmt.Printf("%6.2f C", t)
			}
			if printErr(s.Err()) {
				continue start
			}
		}
		fmt.Println("\nDone.")
		delay.Millisec(4e3)
	}
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.USART2:       conISR,
	irq.DMA1_Stream5: conRxDMAISR,
	irq.DMA1_Stream6: conTxDMAISR,

	irq.USART6:       oneISR,
	irq.DMA2_Stream1: oneRxDMAISR,
	irq.DMA2_Stream7: oneTxDMAISR,
}
