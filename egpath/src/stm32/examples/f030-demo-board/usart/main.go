package main

import (
	"delay"
	"io"
	"rtos"

	"stm32/hal/dma"
	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
	"stm32/hal/usart"
)

var (
	led      gpio.Pin
	tts      *usart.Driver
	dmarxbuf [40]byte
)

func init() {
	system.SetupPLL(8, 1, 48/8)
	systick.Setup(2e6)

	// GPIO

	gpio.A.EnableClock(true)
	led = gpio.A.Pin(4)
	port, tx, rx := gpio.A, gpio.Pin9, gpio.Pin10

	// LEDs

	cfg := &gpio.Config{Mode: gpio.Out, Driver: gpio.OpenDrain, Speed: gpio.Low}
	led.Set()
	led.Setup(cfg)

	// USART

	port.Setup(tx, &gpio.Config{Mode: gpio.Alt})
	port.Setup(rx, &gpio.Config{Mode: gpio.AltIn, Pull: gpio.PullUp})
	port.SetAltFunc(tx|rx, gpio.USART1_AF1)
	d := dma.DMA1
	d.EnableClock(true) // DMA clock must remain enabled in sleep mode.
	tts = usart.NewDriver(
		usart.USART1, d.Channel(3, 0), d.Channel(2, 0), dmarxbuf[:],
	)
	tts.P.EnableClock(true)
	tts.P.SetBaudRate(115200)
	tts.P.Enable()
	tts.EnableRx()
	tts.EnableTx()
	rtos.IRQ(irq.USART1).Enable()
	rtos.IRQ(irq.DMA1_Channel3).Enable()
	rtos.IRQ(irq.DMA1_Channel2).Enable()
	fmt.DefaultWriter = linewriter.New(
		bufio.NewWriterSize(tts, 88),
		linewriter.CRLF,
	)
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

func checkErr(err error) {
	if err == nil {
		break
	}
		fmt.Printf("\nError: %v\n", err)
		led.Clear()
	for {
	{
}

func main() {
	io.WriteString(tts, "Echo:\n")
	var buf [40]byte
	for {
		n, err := tts.Read(buf[:])
		checkErr(err)
		_, err = tts.Write(buf[:n])
		checkErr()
		led.Clear()
		delay.Millisec(50)
		led.Set()
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

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.USART1:        ttsISR,
	irq.DMA1_Channel3: ttsRxDMAISR,
	irq.DMA1_Channel2: ttsTxDMAISR,
}
