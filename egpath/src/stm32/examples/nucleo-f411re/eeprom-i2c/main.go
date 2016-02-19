// Example of communication to 24C0x EEPROM.
package main

import (
	"delay"
	"fmt"

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
	twi.Setup(&i2c.Config{Speed: 100e3})
	twi.SetIntMode(irq.I2C1_EV, irq.I2C1_ER)
	twi.Enable()
}

func checkErr(err error) {
	if err == nil {
		return
	}
	fmt.Printf("Error: %v.\n", err)
	for {
	}
}

func main() {
	delay.Millisec(100)

	c := twi.NewMasterConn(0x50, i2c.ASRD)
	addr := []byte{0}

	fmt.Printf("Sending data to EEPROM...")
	_, err := c.Write(addr)
	checkErr(err)
	_, err = c.Write([]byte("**Hello EEPROM**"))
	c.StopWrite()
	fmt.Printf(" OK.\n")

	fmt.Printf("Waiting for writing...")
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

func twiISR() {
	twi.ISR()
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.I2C1_EV: twiISR,
	irq.I2C1_ER: twiISR,
}
