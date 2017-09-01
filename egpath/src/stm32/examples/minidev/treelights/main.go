package main

import (
	"delay"
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
	system.Setup(8, 1, 72/8)
	systick.Setup()

	// GPIO

	gpio.A.EnableClock(true)
	auxpin := gpio.A.Pin(0)
	ledspin := gpio.A.Pin(2)

	gpio.B.EnableClock(true)
	audioport, audiopins := gpio.B, gpio.Pin8|gpio.Pin9
	usport, ustx, usrx := gpio.B, gpio.Pin10, gpio.Pin11

	// LED USART

	ledspin.Setup(&gpio.Config{Mode: gpio.Alt})
	d := dma.DMA1
	d.EnableClock(true)
	leds = usart.NewDriver(usart.USART2, nil, d.Channel(7, 0), nil)
	leds.P.EnableClock(true)
	leds.P.SetBaudRate(2250e3)    // 36MHz/16: 444 ns/UARTbit, 1333 ns/WS2811bit.
	leds.P.SetConf(usart.Stop0b5) // F103 has no 7-bit mode: save 0.5 bit.
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
	audio.Setup(tim.TIM4, system.APB1.Clock(), 14700, 136*3) // max*3 => gain/3.
	rtos.IRQ(irq.TIM4).Enable()

	// AUX output

	auxpin.Setup(&gpio.Config{Mode: gpio.Out, Speed: gpio.Low})

	rnd.Seed(rtos.Nanosec())
}

func checkErr(err error) {
	if err != nil {
		dbg := rtos.Debug(0)
		dbg.WriteString(err.Error())
		dbg.WriteString("\n")
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
			// Full snd and duration 4 correspond to quarter note (650 ms).
			const Q = 650
			if d < 4 {
				// Shorten snd.
				snd = snd[:(len(snd)*d+2)/4]
			} else if d > 4 {
				// Apply some delay after snd.
				dly = (Q*d+2)/4 - Q
			}
			audio.Play(snd)
			delay.Millisec(dly)
		}
	}
}

func main() {
	buf := make([]byte, 2)
	ledram := ws281x.MakeFBU(50)
	pixel := ws281x.MakeFBU(1)
	white := ws281x.MakeFBU(1)
	red := ws281x.MakeFBU(1)
	white.EncodeRGB(ws281x.Color(0xffffff))
	red.EncodeRGB(ws281x.Color(0xff0000))

	delay.Millisec(200) // Wait for US-100 startup.

	var color ws281x.Color
	color1 := rndColor()
	color2 := rndColor()
	dist := 14000 // 14 m
	iter := 0
	for {
		// Read distance from ultrasonic sensor.
		checkErr(us100.WriteByte(0x55))
		_, err := io.ReadFull(us100, buf)
		checkErr(err)
		dist = (dist + int(buf[0])<<8 + int(buf[1])) / 2

		if dist < 1500 {
			dist = 14000

			// Slowly dim the current color.
			r, g, b := color.RGB()
			const N = 256
			for i := N; i >= 0; i-- {
				pixel.EncodeRGB(ws281x.RGB(r*i/N, g*i/N, b*i/N).Gamma())
				ledram.Fill(pixel)
				leds.Write(ledram.Bytes())
				delay.Millisec(6)
			}
			delay.Millisec(500)

			// Light the red color and play music.
			ledram.Fill(red)
			leds.Write(ledram.Bytes())
			switch iter % 3 {
			case 0:
				play(melody0, 2)
			case 1:
				play(melody1, 2)
			case 2:
				play(melody2, 3)
			}

			// Turn off LEDs in sequence (starting from both ends).
			for i := 0; i < ledram.Len()/2; i++ {
				ledram.At(i).Head(1).Clear()
				ledram.At(49 - i).Head(1).Clear()
				leds.Write(ledram.Bytes())
				delay.Millisec(10)
			}
			delay.Millisec(400)
		}

		// Slowly change the current color.
		const N = 64
		r := (color1.Red()*(N-iter) + color2.Red()*iter) / N
		g := (color1.Green()*(N-iter) + color2.Green()*iter) / N
		b := (color1.Blue()*(N-iter) + color2.Blue()*iter) / N
		color = ws281x.RGB(r, g, b)
		if iter++; iter > N {
			iter = 0
			color1 = color2
			color2 = rndColor()
		}
		pixel.EncodeRGB(color.Gamma())
		ledram.Fill(pixel)
		leds.Write(ledram.Bytes())

		// Sparkle effect.
		for i := 0; i < 10; i++ {
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
