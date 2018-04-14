package main

import (
	"io"
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
	system.SetupPLL(8, 1, 48/8)
	systick.Setup(2e6)

	gpio.A.EnableClock(true)
	tx := gpio.A.Pin(9)

	tx.Setup(&gpio.Config{Mode: gpio.Alt})
	tx.SetAltFunc(gpio.USART1_AF1)
	d := dma.DMA1
	d.EnableClock(true)
	tts = usart.NewDriver(usart.USART1, d.Channel(2, 0), nil, nil)
	tts.Periph().EnableClock(true)
	tts.Periph().SetBaudRate(115200)
	tts.Periph().Enable()
	tts.EnableTx()

	rtos.IRQ(irq.USART1).Enable()
	rtos.IRQ(irq.DMA1_Channel2_3).Enable()
}

type MorseWriter struct {
	W io.Writer
}

func (w *MorseWriter) Write(s []byte) (n int, err error) {
	var buf [8]byte
	for _, c := range s {
		if c < ' ' {
			continue
		}
		if 'a' <= c && c <= 'z' {
			c -= 'a' - 'A' // Convert to upper case.
		}
		if c > 'Z' {
			continue
		}
		var symbol morseSymbol
		if c == ' ' || c == '\n' {
			symbol.length = 2
			buf[0] = ' '
			buf[1] = ' '
		} else {
			symbol = morseSymbols[c-'!']
			for i := uint(0); i < uint(symbol.length); i++ {
				if (symbol.code>>i)&1 != 0 {
					buf[i] = '-'
				} else {
					buf[i] = '.'
				}
			}
		}
		buf[symbol.length] = ' '
		m, err := w.W.Write(buf[:symbol.length+1])
		n += m
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

func main() {
	s := "Hello, World!\r\n"
	w := &MorseWriter{tts}

	io.WriteString(tts, s)
	io.WriteString(w, s)
}

func ttsISR() {
	tts.ISR()
}

func ttsDMAISR() {
	tts.TxDMAISR()
}

//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.USART1:          ttsISR,
	irq.DMA1_Channel2_3: ttsDMAISR,
}

type morseSymbol struct {
	code, length byte
}

//emgo:const
var morseSymbols = [...]morseSymbol{
	{1<<0 | 1<<1 | 1<<2, 4}, // ! ---.
	{1<<1 | 1<<4, 6},        // " .-..-.
	{},                      // #
	{1<<3 | 1<<6, 7},        // $ ...-..-
	{},                      // %
	{},                      // &
	{1<<1 | 1<<2 | 1<<3 | 1<<4, 6}, // ' .----.
	{1<<0 | 1<<2 | 1<<3, 5},        // ( -.--.
	{1<<0 | 1<<2 | 1<<3 | 1<<5, 6}, // ) -.--.-
	{},                                    // *
	{1<<1 | 1<<3, 5},                      // + .-.-.
	{1<<0 | 1<<1 | 1<<4 | 1<<5, 6},        // , --..--
	{1<<0 | 1<<5, 6},                      // - -....-
	{1<<1 | 1<<3 | 1<<5, 6},               // . .-.-.-
	{1<<0 | 1<<3, 5},                      // / -..-.
	{1<<0 | 1<<1 | 1<<2 | 1<<3 | 1<<4, 5}, // 0 -----
	{1<<1 | 1<<2 | 1<<3 | 1<<4, 5},        // 1 .----
	{1<<2 | 1<<3 | 1<<4, 5},               // 2 ..---
	{1<<3 | 1<<4, 5},                      // 3 ...--
	{1 << 4, 5},                           // 4 ....-
	{0, 5},                                // 5 .....
	{1 << 0, 5},                           // 6 -....
	{1<<0 | 1<<1, 5},                      // 7 --...
	{1<<0 | 1<<1 | 1<<2, 5},               // 8 ---..
	{1<<0 | 1<<1 | 1<<2 | 1<<3, 5},        // 9 ----.
	{1<<0 | 1<<1 | 1<<2, 6},               // : ---...
	{1<<0 | 1<<2 | 1<<4, 6},               // ; -.-.-.
	{},                      // <
	{1<<0 | 1<<4, 5},        // = -...-
	{},                      // >
	{1<<2 | 1<<3, 6},        // ? ..--..
	{1<<1 | 1<<2 | 1<<4, 6}, // @ .--.-.
	{1 << 1, 2},             // A .-
	{1 << 0, 4},             // B -...
	{1<<0 | 1<<2, 4},        // C -.-.
	{1 << 0, 3},             // D -..
	{0, 1},                  // E .
	{1 << 2, 4},             // F ..-.
	{1<<0 | 1<<1, 3},        // G --.
	{0, 4},                  // H ....
	{0, 2},                  // I ..
	{1<<1 | 1<<2 | 1<<3, 4}, // J .---
	{1<<0 | 1<<2, 3},        // K -.-
	{1 << 1, 4},             // L .-..
	{1<<0 | 1<<1, 2},        // M --
	{1 << 0, 2},             // N -.
	{1<<0 | 1<<1 | 1<<2, 3}, // O ---
	{1<<1 | 1<<2, 4},        // P .--.
	{1<<0 | 1<<1 | 1<<3, 4}, // Q --.-
	{1 << 1, 3},             // R .-.
	{0, 3},                  // S ...
	{1 << 0, 1},             // T -
	{1 << 2, 3},             // U ..-
	{1 << 3, 4},             // V ...-
	{1<<1 | 1<<2, 3},        // W .--
	{1<<0 | 1<<3, 4},        // X -..-
	{1<<0 | 1<<2 | 1<<3, 4}, // Y -.--
	{1<<0 | 1<<1, 4},        // Z --..
}
