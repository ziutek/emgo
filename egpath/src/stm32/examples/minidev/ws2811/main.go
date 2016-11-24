package main

import (
	"bytes"
	"delay"
	"fmt"
	"math/rand"
	"rtos"

	"stm32/hal/dma"
	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
	"stm32/hal/usart"
)

var tts *usart.Driver

func init() {
	system.Setup(8, 72/8, false)
	systick.Setup()

	// GPIO

	gpio.A.EnableClock(true)
	port, tx := gpio.A, gpio.Pin2

	// USART

	port.Setup(tx, &gpio.Config{Mode: gpio.Alt})
	d := dma.DMA1
	d.EnableClock(true) // DMA clock must remain enabled in s
	tts = usart.NewDriver(usart.USART2, nil, d.Channel(7, 0), nil)
	tts.P.EnableClock(true)
	tts.P.SetBaudRate(2250e3)    // 36MHz/16: 444 ns/UARTbit, 1333 ns/WS2811bit.
	tts.P.SetConf(usart.Stop0b5) // F103 has no 7-bit mode: save 0.5 bit only.
	tts.P.Enable()
	tts.EnableTx()
	rtos.IRQ(irq.USART1).Enable()
	rtos.IRQ(irq.DMA1_Channel7).Enable()
}

func wsb(bit int) byte {
	return byte(6 &^ (bit << 1))
}

const pixelLen = 8 // One byte in RAM represents three WS2811 bits.

type RGB24 struct {
	R, G, B byte
}

const Zero = (6>>1 + 6<<2 + 6<<5)

func encodeRGB24(pixel []byte, c RGB24) {
	r := gamma[c.R]
	g := gamma[c.G]
	b := gamma[c.B]
	pixel[0] = Zero &^ (r>>7&1 + r>>3&8 + r<<1&0x40)
	pixel[1] = Zero &^ (r>>4&1 + r>>0&8 + r<<4&0x40)
	pixel[2] = Zero &^ (r>>1&1 + r<<3&8 + g>>1&0x40)
	pixel[3] = Zero &^ (g>>6&1 + g>>2&8 + g<<2&0x40)
	pixel[4] = Zero &^ (g>>3&1 + g<<1&8 + g<<5&0x40)
	pixel[5] = Zero &^ (g>>0&1 + b>>4&8 + b>>0&0x40)
	pixel[6] = Zero &^ (b>>5&1 + b>>1&8 + b<<3&0x40)
	pixel[7] = Zero &^ (b>>2&1 + b<<2&8 + b<<6&0x40)
}

func main() {
	delay.Millisec(250)

	ledram := make([]byte, 50*pixelLen)
	pixel := make([]byte, pixelLen)
	var rnd rand.XorShift64
	rnd.Seed(rtos.Nanosec())
	for {
		r := rnd.Uint32()
		c := RGB24{byte(r & 0xff), byte(r >> 8 & 0xff), byte(r >> 16 & 0xff)}
		fmt.Printf("%3d %3d %3d\n", c.R, c.G, c.B)
		encodeRGB24(pixel, c)
		for i := 0; i < len(ledram); i += pixelLen {
			bytes.Fill(ledram, Zero)
			copy(ledram[i:], pixel)
			tts.Write(ledram)
			delay.Millisec(10)
		}
		for i := 0; i < len(ledram); i += pixelLen {
			copy(ledram[i:], pixel)
		}
		tts.Write(ledram)
		delay.Millisec(450)
	}
}

func ttsISR() {
	tts.ISR()
}

func ttsTxDMAISR() {
	tts.TxDMAISR()
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.USART2:        ttsISR,
	irq.DMA1_Channel7: ttsTxDMAISR,
}

// Gamma table according to CIE 1976.
//
//	const max = 255
//
//	for i := 0; i <= max; i++ {
//		var x float64
//		y := 100 * float64(i) / max
//		if y > 8 {
//			x = math.Pow((y+16)/116, 3)
//		} else {
//			x = y / 903.3
//		}
//		fmt.Printf("%d,", int(max*x+0.5))
//	}
//
//emgo.const
var gamma = [256]byte{
	0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 2, 3, 3,
	3, 3, 3, 3, 3, 3, 4, 4, 4, 4, 4, 4, 5, 5, 5, 5, 5, 6, 6, 6, 6, 6, 7, 7, 7,
	7, 8, 8, 8, 8, 9, 9, 9, 10, 10, 10, 10, 11, 11, 11, 12, 12, 12, 13, 13, 13,
	14, 14, 15, 15, 15, 16, 16, 17, 17, 17, 18, 18, 19, 19, 20, 20, 21, 21, 22,
	22, 23, 23, 24, 24, 25, 25, 26, 26, 27, 28, 28, 29, 29, 30, 31, 31, 32, 32,
	33, 34, 34, 35, 36, 37, 37, 38, 39, 39, 40, 41, 42, 43, 43, 44, 45, 46, 47,
	47, 48, 49, 50, 51, 52, 53, 54, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64,
	65, 66, 67, 68, 70, 71, 72, 73, 74, 75, 76, 77, 79, 80, 81, 82, 83, 85, 86,
	87, 88, 90, 91, 92, 94, 95, 96, 98, 99, 100, 102, 103, 105, 106, 108, 109,
	110, 112, 113, 115, 116, 118, 120, 121, 123, 124, 126, 128, 129, 131, 132,
	134, 136, 138, 139, 141, 143, 145, 146, 148, 150, 152, 154, 155, 157, 159,
	161, 163, 165, 167, 169, 171, 173, 175, 177, 179, 181, 183, 185, 187, 189,
	191, 193, 196, 198, 200, 202, 204, 207, 209, 211, 214, 216, 218, 220, 223,
	225, 228, 230, 232, 235, 237, 240, 242, 245, 247, 250, 252, 255,
}
