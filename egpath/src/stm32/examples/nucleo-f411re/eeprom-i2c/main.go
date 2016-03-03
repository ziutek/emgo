// Example of communication to 24C0x EEPROM.
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
	d.EnableClock(true) // DMA clock must remain enabled in sleep mode.
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

func main() {
	delay.Millisec(100)

	c := twi.NewMasterConn(0x50, i2c.ASRD)
	addr := []byte{0}

	fmt.Printf("Sending data to EEPROM... ")
	_, err := c.Write(addr)
	checkErr(err)
	_, err = c.Write([]byte("**Hello EEPROM**"))
	checkErr(err)
	c.StopWrite()
	fmt.Printf("OK.\n")

	fmt.Printf("Waiting for writing... ")
	for {
		_, err = c.Write(addr)
		if err == nil {
			break
		}
		if e, ok := err.(i2c.Error); !ok || e != i2c.AckFail {
			checkErr(err)
		}
		fmt.Printf(".")
	}
	fmt.Printf(" OK.\n")

	var buf [16]byte
	_, err = c.Read(buf[:])
	checkErr(err)
	fmt.Printf("%s\n", buf[:])
}

func twiI2CISR() {
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
	irq.I2C1_EV:      twiI2CISR,
	irq.I2C1_ER:      twiI2CISR,
	irq.DMA1_Stream5: twiRxDMAISR,
	irq.DMA1_Stream6: twiTxDMAISR,
}

func checkErr(err error) {
	if err == nil {
		return
	}
	fmt.Printf("Error:")
	if e, ok := err.(i2c.Error); ok {
		if e&i2c.BusErr != 0 {
			fmt.Printf(" BusErr")
		} else if e&i2c.ArbLost != 0 {
			fmt.Printf(" ArbLost")
		} else if e&i2c.AckFail != 0 {
			fmt.Printf(" AckFail")
		} else if e&i2c.Overrun != 0 {
			fmt.Printf(" Overrun")
		} else if e&i2c.PECErr != 0 {
			fmt.Printf(" PECErr")
		} else if e&i2c.Timeout != 0 {
			fmt.Printf(" Timeout")
		} else if e&i2c.SMBAlert != 0 {
			fmt.Printf(" SMBAlert")
		} else if e&i2c.SoftTimeout != 0 {
			fmt.Printf(" SoftTimeout")
		} else if e&i2c.BelatedStop != 0 {
			fmt.Printf(" BelatedStop")
		} else if e&i2c.ActiveRead != 0 {
			fmt.Printf(" ActiveRead")
		} else if e&i2c.DMAErr != 0 {
			fmt.Printf(" DMAErr")
		}
	} else {
		fmt.Printf(" %v", err)
	}
	fmt.Println()
	for {
	}
}
