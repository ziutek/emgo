package main

import (
	"delay"
	"io"
	"math/rand"
	"rtos"

	"led"
	"led/ws281x/wsuart"

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
	system.SetupPLL(8, 1, 72/8)
	systick.Setup(2e6)

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

	// Set UART baudrate to 2250 kb/s (36 MHz / 16 = 2.25 MHz). This gives
	// 444 ns/UARTbit and 1333 ns/WS2811bit. It would be best to use the 7-bit
	// mode but F103 does not support it. Use 0.5 stop bit to slightly speed-up
	// transmission.

	leds = usart.NewDriver(usart.USART2, d.Channel(7, 0), nil, nil)
	leds.Periph().EnableClock(true)
	leds.Periph().SetBaudRate(2250e3)
	leds.Periph().SetConf2(usart.Stop0b5)
	leds.Periph().Enable()
	leds.EnableTx()
	rtos.IRQ(irq.USART2).Enable()
	rtos.IRQ(irq.DMA1_Channel7).Enable()

	// US-100 USART

	usport.Setup(ustx, &gpio.Config{Mode: gpio.Alt})
	usport.Setup(usrx, &gpio.Config{Mode: gpio.AltIn, Pull: gpio.PullUp})
	d = dma.DMA1
	d.EnableClock(true)
	us100 = usart.NewDriver(
		usart.USART3, d.Channel(2, 0), d.Channel(3, 0), make([]byte, 8),
	)
	us100.Periph().EnableClock(true)
	us100.Periph().SetBaudRate(9600)
	us100.Periph().Enable()
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

func rndColor() led.Color {
	x := rnd.Uint32()
	r := byte(x * 3 / 4)
	g := byte(x >> 8 * 3 / 4)
	b := byte(x >> 16 * 3 / 4)
	return led.RGB(r, g, b)
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
	strip := make(wsuart.Strip, 50)
	strip.Clear()
	rgb := wsuart.RGB
	black := rgb.Pixel(led.RGB(0, 0, 0))
	white := rgb.Pixel(led.RGB(255, 255, 255))
	red := rgb.Pixel(led.RGB(255, 0, 0))

	delay.Millisec(200) // Wait for US-100 startup.

	dist := 14000 // 14 m
	colorMax := 64
	colorIter := colorMax
	color1 := led.Color(0)
	color2 := rndColor()
	for {
		// Slowly change the current color.
		if colorIter++; colorIter > colorMax {
			colorIter = 0
			color1 = color2
			color2 = rndColor()
		}
		color := color1.Blend(color2, byte(255*colorIter/colorMax))
		pixel := rgb.Pixel(color)
		strip.Fill(pixel)
		leds.Write(strip.Bytes())

		// Sparkle effect.
		for i := 0; i < 10; i++ {
			if r := int(rnd.Uint32() & 0x1ff); r < len(strip) {
				strip[r] = white
				leds.Write(strip.Bytes())
				delay.Millisec(20)
				strip[r] = pixel
				leds.Write(strip.Bytes())
			} else {
				delay.Millisec(60)
			}
		}

		// Read distance from ultrasonic sensor.
		checkErr(us100.WriteByte(0x55))
		_, err := io.ReadFull(us100, buf)
		checkErr(err)
		dist = (dist + int(buf[0])<<8 + int(buf[1])) / 2

		if dist < 1500 {
			dist = 14000

			// Slowly dim the current color.
			for i := 255; i >= 0; i-- {
				strip.Fill(rgb.Pixel(color.Scale(byte(i))))
				leds.Write(strip.Bytes())
				delay.Millisec(6)
			}
			delay.Millisec(500)

			// Light the red color and play music.
			strip.Fill(red)
			leds.Write(strip.Bytes())
			switch colorIter % 3 {
			case 0:
				play(melody0, 2)
			case 1:
				play(melody1, 2)
			case 2:
				play(melody2, 3)
			}

			// Turn off LEDs in sequence (starting from both ends).
			for i := 0; i < len(strip)/2; i++ {
				strip[i] = black
				strip[len(strip)-1-i] = black
				leds.Write(strip.Bytes())
				delay.Millisec(10)
			}
			delay.Millisec(400)
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
