// This example shows how to use PCF8574T + HD44780 combo.
//
// Connections:
// P0 --> RS
// P1 --> R/W
// P2 --> E
// P3 --> Backlight
// P4 <-> DB4
// P5 <-> DB5
// P6 <-> DB6
// P7 <-> DB7
//
// It seems that PCF8574T works up to 480 kHz (VCC = 5V, short cables, 16:9).
package main

import (
	"bytes"
	"fmt"
	"rtos"

	"hdc"

	"stm32/hal/gpio"
	"stm32/hal/i2c"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
)

var twi *i2c.Driver

func init() {
	system.Setup96(8)
	systick.Setup()

	gpio.B.EnableClock(true)
	port, pins := gpio.B, gpio.Pin8|gpio.Pin9

	cfg := gpio.Config{
		Mode:   gpio.Alt,
		Driver: gpio.OpenDrain,
	}
	port.Setup(pins, &cfg)
	port.SetAltFunc(pins, gpio.I2C1)
	twi = i2c.NewDriver(i2c.I2C1)
	twi.EnableClock(true)
	twi.Reset() // Mandatory!
	twi.Setup(&i2c.Config{Speed: 240e3, Duty: i2c.Duty16_9})
	twi.SetIntMode(true)
	twi.Enable()
	rtos.IRQ(irq.I2C1_EV).Enable()
	rtos.IRQ(irq.I2C1_ER).Enable()
}

func checkErr(err error) {
	if err == nil {
		return
	}
	fmt.Printf("Error: %v.\n", err)
	for {
	}
}

var (
	lcd = &hdc.Display{
		Cols: 20, Rows: 4,
		DS: 4,
		RS: 1 << 0, RW: 1 << 1, E: 1 << 2, AUX: 1 << 3,
	}
	screen [20 * 4]byte
)

func main() {
	lcd.ReadWriter = twi.NewMasterConn(0x27, i2c.ASRD)

	checkErr(lcd.Init())
	checkErr(lcd.SetDisplayMode(hdc.DisplayOn))
	checkErr(lcd.SetAUX()) // Backlight on.

	for i := range screen[:] {
		screen[i] = ' '
	}
	line := bytes.MakeBuffer(screen[2*20:], true)
	var t1 int64
	for i := 0; ; i++ {
		t := rtos.Nanosec()
		fps := 1e9 / float32(t-t1)
		t1 = t
		line.Reset()
		fmt.Fprintf(&line, "%7d %6.1f FPS", i, fps)
		writeScreen(lcd, &screen)
	}
}

func writeScreen(lcd *hdc.Display, screen *[4 * 20]byte) {
	_, err := lcd.Write(screen[0:20])
	checkErr(err)
	_, err = lcd.Write(screen[40:60])
	checkErr(err)
	_, err = lcd.Write(screen[20:40])
	checkErr(err)
	_, err = lcd.Write(screen[60:80])
	checkErr(err)
}

func twiISR() {
	twi.ISR()
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.I2C1_EV: twiISR,
	irq.I2C1_ER: twiISR,
}
