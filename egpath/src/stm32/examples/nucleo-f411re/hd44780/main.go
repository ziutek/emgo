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
	"delay"
	"fmt"
	"rtos"

	"hdc"

	"stm32/hal/dma"
	"stm32/hal/gpio"
	"stm32/hal/i2c"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
)

const (
	driver  = 2 // Select different drivers (1,2) to see their performance.
	hdcaddr = 0x27
)

var (
	adrv *i2c.AltDriver
	drv  *i2c.Driver

	lcd = &hdc.Display{
		Cols: 20, Rows: 4,
		DS: 4,
		RS: 1 << 0, RW: 1 << 1, E: 1 << 2, AUX: 1 << 3,
	}
	screen [20 * 4]byte
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

	twi := i2c.I2C1
	twi.EnableClock(true)
	twi.Reset() // Mandatory!
	twi.Setup(&i2c.Config{Speed: 480e3, Duty: i2c.Duty16_9})
	switch driver {
	case 1:
		adrv = i2c.NewAltDriver(twi)
		adrv.SetIntMode(true)
		lcd.ReadWriter = adrv.NewMasterConn(hdcaddr, i2c.ASRD)
	case 2:
		d := dma.DMA1
		d.EnableClock(true) // DMA clock must remain enabled in sleep mode.
		drv = i2c.NewDriver(twi, d.Channel(5, 1), d.Channel(6, 1))
		//drv = i2c.NewDriver(twi, nil, nil)
		rtos.IRQ(irq.DMA1_Stream5).Enable()
		rtos.IRQ(irq.DMA1_Stream6).Enable()
		lcd.ReadWriter = drv.NewMasterConn(hdcaddr, i2c.ASRD)
	}
	twi.Enable()
	rtos.IRQ(irq.I2C1_EV).Enable()
	rtos.IRQ(irq.I2C1_ER).Enable()
}

func main() {
	delay.Millisec(250)

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

func twiI2CISR() {
	switch driver {
	case 1:
		adrv.ISR()
	case 2:
		drv.I2CISR()
	}
}

func twiRxDMAISR() {
	drv.DMAISR(drv.RxDMA)
}

func twiTxDMAISR() {
	drv.DMAISR(drv.TxDMA)
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.I2C1_EV:      twiI2CISR,
	irq.I2C1_ER:      twiI2CISR,
	irq.DMA1_Stream5: twiRxDMAISR,
	irq.DMA1_Stream6: twiTxDMAISR,
}

func checkErr(err error) {
	if err == nil {
		return
	}
	if e, ok := err.(i2c.Error); ok {
		fmt.Printf("I2C error:")
		if e&i2c.BusErr != 0 {
			fmt.Printf(" BusErr")
		}
		if e&i2c.ArbLost != 0 {
			fmt.Printf(" ArbLost")
		}
		if e&i2c.AckFail != 0 {
			fmt.Printf(" AckFail")
		}
		if e&i2c.Overrun != 0 {
			fmt.Printf(" Overrun")
		}
		if e&i2c.PECErr != 0 {
			fmt.Printf(" PECErr")
		}
		if e&i2c.Timeout != 0 {
			fmt.Printf(" Timeout")
		}
		if e&i2c.SMBAlert != 0 {
			fmt.Printf(" SMBAlert")
		}
		if e&i2c.SoftTimeout != 0 {
			fmt.Printf(" SoftTimeout")
		}
		if e&i2c.BelatedStop != 0 {
			fmt.Printf(" BelatedStop")
		}
		if e&i2c.BadEvent != 0 {
			fmt.Printf(" BadEvent")
		}
		if e&i2c.ActiveRead != 0 {
			fmt.Printf(" ActiveRead")
		}
		if e&i2c.DMAErr != 0 {
			fmt.Printf(" DMAErr")
		}
		fmt.Println()
	} else {
		fmt.Printf("Error %v\n", err)
	}
	for {
	}
}
