// This example shows how to use HD44780 + PCF8574T combo.
//
// For PCF8574T and HD44780 LCD. Connections:
// P0 --> RS
// P1 --> R/W
// P2 --> E
// P3 --> Backlight
// P4 <-> DB4
// P5 <-> DB5
// P6 <-> DB6
// P7 <-> DB7
//
// It seems that PCF8574T works well up to 200 kHz (VCC = 5V). Tested up to
// 400 kHz but there are no any speed improvements above 200 kHz.
package main

import (
	"bytes"
	"fmt"
	"rtos"

	"hdc"

	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"

	"stm32/hal/i2c"
)

var (
	twiDrv = &i2c.Driver{Periph: i2c.I2C1}
	twiCon = twiDrv.MasterConn(0x27)
)

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
	twiDrv.EnableClock(true)
	twiDrv.Reset() // Mandatory!
	twiDrv.Setup(&i2c.Config{Speed: 200e3})
	twiDrv.SetIntMode(irq.I2C1_EV, irq.I2C1_ER)
	twiDrv.Enable()
	twiCon.SetAutoStop(true, false) // This speedups writing.
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
		Drv:  &twiCon,
		Cols: 20, Rows: 4,
		DLs: 4,
		RS:  1 << 0, RW: 1 << 1, E: 1 << 2, AUX: 1 << 3,
	}
	buf [4 * 20]byte
)

func main() {
	checkErr(lcd.Init())
	checkErr(lcd.SetDisplayMode(hdc.DisplayOn))
	checkErr(lcd.SetAUX())

	for i := range buf[:] {
		buf[i] = ' '
	}
	bb := bytes.MakeBuffer(buf[20:], true)
	var t1 int64
	for i := 0; ; i++ {
		t := rtos.Nanosec()
		fps := 1e9 / float32(t-t1)
		t1 = t
		bb.Reset()
		fmt.Fprintf(&bb, "%7d %6.3g FPS", i, fps)
		_, err := lcd.Write(buf[:])
		checkErr(err)
	}
}

func twiISR() {
	twiDrv.ISR()
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.I2C1_EV: twiISR,
	irq.I2C1_ER: twiISR,
}
