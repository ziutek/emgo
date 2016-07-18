// This example shows how to use USART as 1-wire master.
package main

import (
	"delay"
	"fmt"
	"rtos"

	"onewire"

	"stm32/onedrv"

	"stm32/hal/dma"
	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
	"stm32/hal/usart"
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
	uprt, tx, rx := gpio.A, gpio.Pin2, gpio.Pin3
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

	// USART

	uprt.Setup(tx, &gpio.Config{Mode: gpio.Alt})
	uprt.Setup(rx, &gpio.Config{Mode: gpio.AltIn, Pull: gpio.PullUp})
	uprt.SetAltFunc(tx|rx, gpio.USART2)
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
	fmt.DefaultWriter = con

	// 1-wire

	oprt.Setup(opin, &gpio.Config{Mode: gpio.Alt, Driver: gpio.OpenDrain})
	oprt.SetAltFunc(opin, gpio.USART6)
	one = usart.NewDriver(
		usart.USART6, usart6rxdma, usart6txdma, make([]byte, 16),
	)
	one.EnableClock(true)
	one.SetBaudRate(9600)
	one.SetMode(usart.HalfDuplex)
	one.Enable()
	one.EnableRx()
	one.EnableTx()
	rtos.IRQ(irq.USART6).Enable()
	rtos.IRQ(irq.DMA2_Stream1).Enable()
	rtos.IRQ(irq.DMA2_Stream7).Enable()

	/*
		s = usart.USART6

		s.EnableClock(true)
		s.SetConf(usart.TxEna | usart.RxEna)
		s.SetMode(usart.HalfDuplex)
		s.EnableIRQs(usart.RxNotEmptyIRQ)
		s.Enable()

		one = serial.New(s, 8, 8)
		rtos.IRQ(irq.USART6).Enable()
	*/
}

func conUSARTISR() {
	con.USARTISR()
}
func conRxDMAISR() {
	con.RxDMAISR()
}
func conTxDMAISR() {
	con.TxDMAISR()
}

func oneUSARTISR() {
	one.USARTISR()
}
func oneRxDMAISR() {
	one.RxDMAISR()
}
func oneTxDMAISR() {
	one.TxDMAISR()
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
	fmt.Printf("Error: %v\r\n", err)
	for {
		blink(Green, 100)
		delay.Millisec(100)
	}
}

func main() {
	m := onewire.Master{Driver: onedrv.USARTDriver{one}}
	dtypes := []onewire.Type{onewire.DS18S20, onewire.DS18B20, onewire.DS1822}

	fmt.Printf("\r\nConfigure all DS18B20, DS1822 to 10bit resolution.\r\n")

	// This algorithm configures and starts conversion simultaneously on all
	// temperature sensors on the bus. It is fast, but doesn't work in case of
	// parasite power mode.

	checkErr(m.SkipROM())
	checkErr(m.WriteScratchpad(127, -128, onewire.T10bit))
	for {
		fmt.Printf("\r\nSending ConvertT command on the bus (SkipROM addressing).\r\n")
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
		fmt.Printf("\r\nSearching for temperature sensors:\r\n")
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
				fmt.Printf("%6.2f C\r\n", t)
			}
			checkErr(s.Err())
		}
		fmt.Printf("Done.\r\n")
		delay.Millisec(4e3)
	}
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.USART2:       conUSARTISR,
	irq.DMA1_Stream5: conRxDMAISR,
	irq.DMA1_Stream6: conTxDMAISR,

	irq.USART6:       oneUSARTISR,
	irq.DMA2_Stream1: oneRxDMAISR,
	irq.DMA2_Stream7: oneTxDMAISR,
}
