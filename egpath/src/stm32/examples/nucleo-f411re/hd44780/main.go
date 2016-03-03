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

	"stm32/hal/dma"
	"stm32/hal/gpio"
	"stm32/hal/i2c"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
)

var twi *i2c.DriverDMA

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
	//twi = i2c.NewDriver(i2c.I2C1)
	d := dma.DMA1
	d.EnableClock(true)
	twi = i2c.NewDriverDMA(i2c.I2C1, d.Channel(5, 1), d.Channel(6, 1))
	twi.EnableClock(true)
	twi.Reset() // Mandatory!
	twi.Setup(&i2c.Config{Speed: 240e3, Duty: i2c.Duty16_9})
	twi.SetIntMode(true, true)
	twi.Enable()
	rtos.IRQ(irq.I2C1_EV).Enable()
	rtos.IRQ(irq.I2C1_ER).Enable()
	rtos.IRQ(irq.DMA1_Stream5).Enable()
	rtos.IRQ(irq.DMA1_Stream6).Enable()
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
	twi.I2CISR()
}

func twiRxDMAISR() {
	twi.DMAISR(twi.RxDMA)
}

func twiTxDMAISR() {
	twi.DMAISR(twi.TxDMA)
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.I2C1_EV:      twiISR,
	irq.I2C1_ER:      twiISR,
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
