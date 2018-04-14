package main

import (
	"delay"
	"io"

	"stm32/hal/gpio"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
)

var led gpio.Pin

func init() {
	system.SetupPLL(8, 1, 48/8)
	systick.Setup(2e6)

	gpio.A.EnableClock(false)
	led = gpio.A.Pin(4)

	cfg := gpio.Config{Mode: gpio.Out, Driver: gpio.OpenDrain, Speed: gpio.Low}
	led.Setup(&cfg)
}

type Telegraph struct {
	Pin   gpio.Pin
	Dotms int // Dot length [ms]
}

func (t Telegraph) Write(s []byte) (int, error) {
	for _, c := range s {
		switch c {
		case '.':
			t.Pin.Clear()
			delay.Millisec(t.Dotms)
			t.Pin.Set()
			delay.Millisec(t.Dotms)
		case '-':
			t.Pin.Clear()
			delay.Millisec(3 * t.Dotms)
			t.Pin.Set()
			delay.Millisec(t.Dotms)
		case ' ':
			delay.Millisec(3 * t.Dotms)
		}
	}
	return len(s), nil
}

func main() {
	mt := &MorseWriter{Telegraph{led, 100}}
	for {
		io.WriteString(mt, "Hello, World! ")
	}
}

type MorseWriter struct {
	W io.Writer
}

func (w *MorseWriter) Write(s []byte) (int, error) {
	var buf [8]byte
	for n, c := range s {
		switch {
		case c == '\n':
			c = ' ' // Replace new lines with spaces.
		case 'a' <= c && c <= 'z':
			c -= 'a' - 'A' // Convert to upper case.
		}
		if c < ' ' || 'Z' < c {
			continue // c is outside ASCII [' ', 'Z']
		}
		var symbol morseSymbol
		if c == ' ' {
			symbol.length = 1
			buf[0] = ' '
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
		if _, err := w.W.Write(buf[:symbol.length+1]); err != nil {
			return n, err
		}
	}
	return len(s), nil
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
