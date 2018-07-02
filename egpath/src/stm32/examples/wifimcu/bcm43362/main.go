package main

import (
	"bufio"
	"delay"
	"fmt"
	"rtos"
	"text/linewriter"

	"sdcard"

	"stm32/hal/dma"
	"stm32/hal/exti"
	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/sdmmc"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
	"stm32/hal/usart"
)

var (
	led gpio.Pin
	sd  *sdmmc.DriverDMA
	tts *usart.Driver
)

func init() {
	system.Setup96(26)
	systick.Setup(2e6)

	// GPIO

	gpio.A.EnableClock(false)
	//irq := gpio.A.Pin(0)
	port, tx, rx := gpio.A, gpio.Pin2, gpio.Pin3
	led = gpio.A.Pin(4)
	cmd := gpio.A.Pin(6)
	d1 := gpio.A.Pin(8)
	d2 := gpio.A.Pin(9)

	gpio.B.EnableClock(true)
	d3 := gpio.B.Pin(5)
	d0 := gpio.B.Pin(7)
	clk := gpio.B.Pin(15)

	// LED

	cfg := &gpio.Config{
		Mode:   gpio.Out,
		Driver: gpio.OpenDrain,
		Speed:  gpio.Low,
	}
	led.Setup(cfg)

	// UART

	port.Setup(tx, &gpio.Config{Mode: gpio.Alt})
	port.Setup(rx, &gpio.Config{Mode: gpio.AltIn, Pull: gpio.PullUp})
	port.SetAltFunc(tx|rx, gpio.USART2)
	d := dma.DMA1
	d.EnableClock(true) // DMA clock must remain enabled in sleep mode.
	tts = usart.NewDriver(
		usart.USART2, d.Channel(6, 4), d.Channel(5, 4), make([]byte, 88),
	)
	tts.Periph().EnableClock(true)
	tts.Periph().SetBaudRate(115200)
	tts.Periph().Enable()
	tts.EnableRx()
	tts.EnableTx()
	rtos.IRQ(irq.USART2).Enable()
	rtos.IRQ(irq.DMA1_Stream5).Enable()
	rtos.IRQ(irq.DMA1_Stream6).Enable()
	fmt.DefaultWriter = linewriter.New(
		bufio.NewWriterSize(tts, 88),
		linewriter.CRLF,
	)

	// SDIO (BCM43362)

	cfg = &gpio.Config{Mode: gpio.Alt, Speed: gpio.VeryHigh, Pull: gpio.PullUp}
	for _, pin := range []gpio.Pin{clk, cmd, d0, d1, d2, d3} {
		pin.Setup(cfg)
		pin.SetAltFunc(gpio.SDIO)
	}
	d = dma.DMA2
	d.EnableClock(true)
	sd = sdmmc.NewDriverDMA(sdmmc.SDIO, d.Channel(6, 4), d0)
	sd.Periph().EnableClock(true)
	sd.Periph().Enable()

	rtos.IRQ(irq.SDIO).Enable()
	rtos.IRQ(irq.EXTI9_5).Enable()
}

func checkErr(what string, err error) {
	if err == nil {
		return
	}
	fmt.Printf("%s: %v\n", what, err)
	for {
	}
}

func main() {
	fmt.Printf("Try to communicate with BCM43362...\n")

	// Set SDIO_CK to no more than 400 kHz (max. open-drain freq)..
	sd.SetClock(400e3)

	// SD card power-up takes maximum of 1 ms or 74 SDIO_CK cycles.
	delay.Millisec(1)

	ocr, mem, numIO := sd.SendCmd(sdcard.CMD5(0)).R4()
	checkErr("CMD5", sd.Err(true))

	fmt.Printf("ocr=%08x mem=%t numIO=%d\n", ocr, mem, numIO)

	for {
		led.Clear()
		delay.Millisec(900)
		led.Set()
		delay.Millisec(100)
	}
}

func ttsISR() {
	tts.ISR()
}

func ttsRxDMAISR() {
	tts.RxDMAISR()
}

func ttsTxDMAISR() {
	tts.TxDMAISR()
}

func sdioISR() {
	sd.ISR()
}

func exti9_5ISR() {
	pending := exti.Pending() & 0x3E0
	pending.ClearPending()
	if pending&sd.BusyLine() != 0 {
		sd.BusyISR()
	}
}

//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.USART2:       ttsISR,
	irq.DMA1_Stream5: ttsRxDMAISR,
	irq.DMA1_Stream6: ttsTxDMAISR,

	irq.SDIO:    sdioISR,
	irq.EXTI9_5: exti9_5ISR,
}
