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

	"display/hdc"

	"stm32/hal/dma"
	"stm32/hal/gpio"
	"stm32/hal/i2c"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
)

const (
	driver = 2 // Select different drivers (1,2,3,4) to see their performance.
	addr   = 0x27
)

var (
	drv     *i2c.Driver
	adrv    *i2c.AltDriver
	adrvdma *i2c.AltDriverDMA

	lcd = &hdc.Display{
		Cols: 20, Rows: 4,
		DS: 4,
		RS: 1 << 0, RW: 1 << 1, E: 1 << 2, AUX: 1 << 3,
	}
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
	case 1: // Driver, no DMA
		drv = i2c.NewDriver(twi, nil, nil)
		lcd.ReadWriter = drv.NewMasterConn(addr, i2c.ASRD)
	case 2: // Driver, DMA for Rx and Tx.
		d := dma.DMA1
		d.EnableClock(true) // DMA clock must remain enabled in sleep mode.
		rtos.IRQ(irq.DMA1_Stream5).Enable()
		rtos.IRQ(irq.DMA1_Stream6).Enable()
		drv = i2c.NewDriver(twi, d.Channel(5, 1), d.Channel(6, 1))
		lcd.ReadWriter = drv.NewMasterConn(addr, i2c.ASRD)
	case 3: // AltDriver, interrupt mode.
		adrv = i2c.NewAltDriver(twi)
		adrv.SetIntMode(true)
		lcd.ReadWriter = adrv.NewMasterConn(addr, i2c.ASRD)
	case 4: // AltDriverDMA, interrupt mode.
		d := dma.DMA1
		d.EnableClock(true) // DMA clock must remain enabled in sleep mode.
		rtos.IRQ(irq.DMA1_Stream5).Enable()
		rtos.IRQ(irq.DMA1_Stream6).Enable()
		adrvdma = i2c.NewAltDriverDMA(twi, d.Channel(5, 1), d.Channel(6, 1))
		adrvdma.SetIntMode(true, true)
		lcd.ReadWriter = adrvdma.NewMasterConn(addr, i2c.ASRD)
	}
	twi.Enable()
	rtos.IRQ(irq.I2C1_EV).Enable()
	rtos.IRQ(irq.I2C1_ER).Enable()
}

var screen [20 * 4]byte

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
	case 1, 2:
		drv.I2CISR()
	case 3:
		adrv.ISR()
	default:
		adrvdma.I2CISR()
	}
}

func twiRxDMAISR() {
	if driver == 2 {
		drv.DMAISR(drv.RxDMA)
	} else {
		adrvdma.DMAISR(adrvdma.RxDMA)
	}
}

func twiTxDMAISR() {
	if driver == 2 {
		drv.DMAISR(drv.TxDMA)
	} else {
		adrvdma.DMAISR(adrvdma.TxDMA)
	}
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
	fmt.Printf("Error %v\n", err)
	for {
	}
}
