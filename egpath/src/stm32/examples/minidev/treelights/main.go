package main

import (
	"delay"
	"fmt"
	"io"
	"math/rand"
	"rtos"

	"ws281x"

	"stm32/hal/dma"
	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
	"stm32/hal/usart"

	"stm32/hal/raw/rcc"
	"stm32/hal/raw/tim"
)

var (
	leds, us100 *usart.Driver
	rnd         rand.XorShift64
)

func init() {
	system.Setup(8, 72/8, false)
	systick.Setup()

	// GPIO

	gpio.A.EnableClock(true)
	auxport, auxpin := gpio.A, 0
	ledport, ledpin := gpio.A, gpio.Pin2

	gpio.B.EnableClock(true)
	audioport, audiopins := gpio.B, gpio.Pin8|gpio.Pin9
	usport, ustx, usrx := gpio.B, gpio.Pin10, gpio.Pin11

	// LED USART

	ledport.Setup(ledpin, &gpio.Config{Mode: gpio.Alt})
	d := dma.DMA1
	d.EnableClock(true)
	leds = usart.NewDriver(usart.USART2, nil, d.Channel(7, 0), nil)
	leds.P.EnableClock(true)
	leds.P.SetBaudRate(2250e3)    // 36MHz/16: 444 ns/UARTbit, 1333 ns/WS2811bit.
	leds.P.SetConf(usart.Stop0b5) // F103 has no 7-bit mode: save 0.5 bit only.
	leds.P.Enable()
	leds.EnableTx()
	rtos.IRQ(irq.USART2).Enable()
	rtos.IRQ(irq.DMA1_Channel7).Enable()

	// US-100 USART

	usport.Setup(ustx, &gpio.Config{Mode: gpio.Alt})
	usport.Setup(usrx, &gpio.Config{Mode: gpio.AltIn, Pull: gpio.PullUp})
	d = dma.DMA1
	d.EnableClock(true)
	us100 = usart.NewDriver(
		usart.USART3, d.Channel(3, 0), d.Channel(2, 0), make([]byte, 8),
	)
	us100.P.EnableClock(true)
	us100.P.SetBaudRate(9600)
	us100.P.Enable()
	us100.EnableRx()
	us100.EnableTx()
	rtos.IRQ(irq.USART3).Enable()
	rtos.IRQ(irq.DMA1_Channel2).Enable()
	rtos.IRQ(irq.DMA1_Channel3).Enable()

	// Audio PWM

	audioport.Setup(audiopins, &gpio.Config{Mode: gpio.Alt, Speed: gpio.Low})
	rcc.RCC.TIM4EN().Set()
	// For TIMclk=72Mhz and sr=14700, max=136 gives PSC=9-1 and actual SR=14706.
	audio.Setup(tim.TIM4, system.APB1.Clock(), 14700, 136*3)
	rtos.IRQ(irq.TIM4).Enable()

	// AUX output

	auxport.SetupPin(auxpin, &gpio.Config{Mode: gpio.Out, Speed: gpio.Low})

	rnd.Seed(rtos.Nanosec())
}

func checkErr(err error) {
	if err != nil {
		fmt.Printf("\nError: %v\n", err)
	}
}

func rndColor() ws281x.Color {
	x := int(rnd.Uint32())
	r := x & 0xff * 3 / 4
	g := (x >> 8) & 0xff * 3 / 4
	b := (x >> 16) & 0xff * 3 / 4
	return ws281x.RGB(r, g, b)
}

func play(melody []Note, n int) {
	for i := 0; i < n; i++ {
		for _, note := range melody {
			d := note.Duration
			snd := note.Sample[:]
			dly := 0
			if d < 4 {
				snd = snd[:(len(snd)*d+2)/4]
			} else {
				dly = (650*d+2)/4 - 650
			}
			audio.Play(snd)
			delay.Millisec(dly)
		}
	}
}

func main() {
	buf := make([]byte, 2)
	ledram := ws281x.MakeS3(50)
	pixel := ws281x.MakeS3(1)
	white := ws281x.MakeS3(1)
	red := ws281x.MakeS3(1)
	white.EncodeRGB(ws281x.Color(0xffffff))
	red.EncodeRGB(ws281x.Color(0xff0000))

	delay.Millisec(200) // Wait for US-100 startup.

	var color ws281x.Color
	color1 := rndColor()
	color2 := rndColor()
	x := 14000
	k := 0
	for {
		checkErr(us100.WriteByte(0x55))
		_, err := io.ReadFull(us100, buf)
		checkErr(err)
		x = (x + int(buf[0])<<8 + int(buf[1])) / 2

		switch {
		case x < 1500:
			x = 14000
			r, g, b := color.RGB()
			for n := 256; n >= 0; n-- {
				pixel.EncodeRGB(ws281x.RGB(r*n/256, g*n/256, b*n/256).Gamma())
				ledram.Fill(pixel)
				leds.Write(ledram.Bytes())
				delay.Millisec(6)
			}
			delay.Millisec(500)

			ledram.Fill(red)
			leds.Write(ledram.Bytes())
			if k&1 == 0 {
				play(melody1, 1)
			} else {
				play(melody2, 3)
			}
			for n := 0; n < 25; n++ {
				ledram.At(n).Head(1).Clear()
				ledram.At(49-n).Head(1).Clear()
				leds.Write(ledram.Bytes())
				delay.Millisec(10)
			}
			delay.Millisec(400)

		default:
			const N = 64
			r := (color1.Red()*(N-k) + color2.Red()*k) / N
			g := (color1.Green()*(N-k) + color2.Green()*k) / N
			b := (color1.Blue()*(N-k) + color2.Blue()*k) / N
			color = ws281x.RGB(r, g, b)
			if k++; k > N {
				k = 0
				color1 = color2
				color2 = rndColor()
			}
			pixel.EncodeRGB(color.Gamma())
			ledram.Fill(pixel)
			leds.Write(ledram.Bytes())
			for n := 0; n < 10; n++ {
				if r := int(rnd.Uint32() & 0x1ff); r < ledram.Len() {
					ledram.At(r).Write(white)
					leds.Write(ledram.Bytes())
					delay.Millisec(20)
					ledram.At(r).Write(pixel)
					leds.Write(ledram.Bytes())
				} else {
					delay.Millisec(60)
				}
			}
		}
	}
}

func ledsISR() {
	leds.ISR()
}

func ledsTxDMAISR() {
	leds.TxDMAISR()
}

func us100ISR() {
	us100.ISR()
}

func us100RxDMAISR() {
	us100.RxDMAISR()
}

func us100TxDMAISR() {
	us100.TxDMAISR()
}

func audioISR() {
	audio.ISR()
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.USART2:        ledsISR,
	irq.DMA1_Channel7: ledsTxDMAISR,

	irq.USART3:        us100ISR,
	irq.DMA1_Channel2: us100TxDMAISR,
	irq.DMA1_Channel3: us100RxDMAISR,

	irq.TIM4: audioISR,
}
