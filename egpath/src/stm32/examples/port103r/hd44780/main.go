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
	"fmt"
	"rtos"

	"hdc"

	"stm32/hal/dma"
	"stm32/hal/gpio"
	"stm32/hal/i2c"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/rtc"
)

var (
	twi i2c.DriverDMA
	lcd = &hdc.Display{
		ReadWriter: twi.NewMasterConn(0x27, i2c.ASRD),
		Rows:       20, Cols: 4,
		DS: 4,
		RS: 1 << 0, RW: 1 << 1, E: 1 << 2, AUX: 1 << 3,
	}
)

func init() {
	system.Setup(8, 72/8, false)
	rtc.Setup(32768)

	gpio.B.EnableClock(true)
	port, pins := gpio.B, gpio.Pin10|gpio.Pin11

	cfg := gpio.Config{
		Mode:   gpio.Alt,
		Driver: gpio.OpenDrain,
	}
	port.Setup(pins, &cfg)
	dma := dma.DMA1
	dma.EnableClock(false)
	twi.Periph = i2c.I2C2
	twi.RxDMA = dma.Channel(5, 0)
	twi.TxDMA = dma.Channel(4, 0)
	twi.EnableClock(true)
	twi.Reset() // Mandatory!
	twi.Setup(&i2c.Config{Speed: 100e3})
	twi.SetIntMode(true, true)
	twi.Enable()
	rtos.IRQ(irq.I2C2_EV).Enable()
	rtos.IRQ(irq.I2C2_ER).Enable()
}

func checkErr(err error) {
	if err == nil {
		return
	}
	fmt.Printf("Error: %s\n", err)
	for {
	}
}

func main() {
	checkErr(lcd.Init())
	checkErr(lcd.SetDisplayMode(hdc.DisplayOn))
	checkErr(lcd.SetAUX())
	var t1 int64
	for i := 0; ; i++ {
		t2 := rtos.Nanosec()
		fps := 1e9 / float32(t2-t1)
		t1 = t2
		c := ' ' + i&15
		fmt.Fprintf(lcd, "  Hitachi  Display  ")
		fmt.Fprintf(lcd, "%-10d  %4.1f FPS", i, fps)
		fmt.Fprintf(lcd, "     Controller     ")
		fmt.Fprintf(lcd, "    %c HD44780 %c     ", c, c)
	}
}

func twiISR() {
	twi.ISR()
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTCAlarm: rtc.ISR,
	irq.I2C2_EV:  twiISR,
	irq.I2C2_ER:  twiISR,
}
