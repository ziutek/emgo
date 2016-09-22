// This example blinks leds connected to pins P4, P5, P6, P7 of PCF8574T.
// Button can be connected to other pins and its state observed on SWO.
//
// This application uses simple I2C error recovery. To generate some recoverable
// error disconnect and after that connect PCF8574T.
package main

import (
	"delay"
	"fmt"
	"rtos"

	"stm32/hal/dma"
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
	port.Setup(pins, cfg)
	port.SetAltFunc(pins, gpio.I2C1)
	d := dma.DMA1
	d.EnableClock(true) // DMA clock must remain enabled in sleep mode.
	twi = i2c.NewDriver(i2c.I2C1, d.Channel(5, 1), d.Channel(6, 1))
	twi.EnableClock(true)
	rtos.IRQ(irq.I2C1_EV).Enable()
	rtos.IRQ(irq.I2C1_ER).Enable()
	rtos.IRQ(irq.DMA1_Stream5).Enable()
	rtos.IRQ(irq.DMA1_Stream6).Enable()
}

func twiConfigure() {
	fmt.Printf("Reset\n")
	twi.Reset() // Mandatory!
	twi.Setup(i2c.Config{Speed: 5000})
	twi.Enable()
}

func recover(err error) {
	printError(err)
	delay.Millisec(500) // Reduce CPU load in case of permanent error.
	twiConfigure()
	twi.Unlock()
}

func main() {
	delay.Millisec(200)
	twiConfigure()
	c := twi.MasterConn(0x27, i2c.NOAS)

	out := []byte{
		0xef, 0xef, 0xef, 0xef, 0xef, 0xef, 0xef, 0xef, 0xef, 0xef, 0xef,
		0xef, 0xef, 0xef, 0xef, 0xef, 0xef, 0xef, 0xef, 0xef, 0xef, 0xef,
		0xdf, 0xdf, 0xdf, 0xdf, 0xdf, 0xdf, 0xdf, 0xdf, 0xdf, 0xdf, 0xdf,
		0xdf, 0xdf, 0xdf, 0xdf, 0xdf, 0xdf, 0xdf, 0xdf, 0xdf, 0xdf, 0xdf,
		0xbf, 0xbf, 0xbf, 0xbf, 0xbf, 0xbf, 0xbf, 0xbf, 0xbf, 0xbf, 0xbf,
		0xbf, 0xbf, 0xbf, 0xbf, 0xbf, 0xbf, 0xbf, 0xbf, 0xbf, 0xbf, 0xbf,
		0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f,
		0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f,
		0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f,
		0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f,
		0xbf, 0xbf, 0xbf, 0xbf, 0xbf, 0xbf, 0xbf, 0xbf, 0xbf, 0xbf, 0xbf,
		0xbf, 0xbf, 0xbf, 0xbf, 0xbf, 0xbf, 0xbf, 0xbf, 0xbf, 0xbf, 0xbf,
		0xdf, 0xdf, 0xdf, 0xdf, 0xdf, 0xdf, 0xdf, 0xdf, 0xdf, 0xdf, 0xdf,
		0xdf, 0xdf, 0xdf, 0xdf, 0xdf, 0xdf, 0xdf, 0xdf, 0xdf, 0xdf, 0xdf,
		0xef, 0xef, 0xef, 0xef, 0xef, 0xef, 0xef, 0xef, 0xef, 0xef, 0xef,
		0xef, 0xef, 0xef, 0xef, 0xef, 0xef, 0xef, 0xef, 0xef, 0xef, 0xef,
	}
	for i := 0; ; i++ {
		_, err := c.Write(out)
		if err != nil {
			recover(err)
			continue
		}
		_, err = c.Write(out)
		if err != nil {
			recover(err)
			continue
		}
		var in [8]byte
		n, err := c.Read(in[:4])
		fmt.Printf("%d %2x\n", i, in[:n])
		if err != nil {
			recover(err)
			continue
		}
		c.SetStopRead()
		n, err = c.Read(in[:2])
		fmt.Printf("%d %2x\n", i, in[:n])
		if err != nil {
			recover(err)
			continue
		}
	}
}

func twiISR() {
	twi.I2CISR()
}

func twiTxDMAISR() {
	twi.DMAISR(twi.TxDMA)
}

func twiRxDMAISR() {
	twi.DMAISR(twi.RxDMA)
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.I2C1_EV:      twiISR,
	irq.I2C1_ER:      twiISR,
	irq.DMA1_Stream5: twiRxDMAISR,
	irq.DMA1_Stream6: twiTxDMAISR,
}

func printError(err error) {
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
}
