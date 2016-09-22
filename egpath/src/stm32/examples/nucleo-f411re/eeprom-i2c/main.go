// Example of communication to 24C0x EEPROM.
package main

import (
	"delay"
	"fmt"
	"io"
	"rtos"

	"stm32/hal/dma"
	"stm32/hal/gpio"
	"stm32/hal/i2c"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
)

const (
	driver = 2
	addr   = 0x50
)

type conn interface {
	io.ReadWriter
	StopWrite()
	UnlockDriver()
}

var (
	drv     *i2c.Driver
	adrv    *i2c.AltDriver
	adrvdma *i2c.AltDriverDMA

	eeprom conn
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
	port.Setup(pins, cfg)
	port.SetAltFunc(pins, gpio.I2C1)

	twi := i2c.I2C1
	twi.EnableClock(true)
	twi.Reset() // Mandatory!
	twi.Setup(i2c.Config{Speed: 480e3, Duty: i2c.Duty16_9})
	switch driver {
	case 1: // Driver, no DMA
		drv = i2c.NewDriver(twi, nil, nil)
		eeprom = drv.NewMasterConn(addr, i2c.ASRD)
	case 2: // Driver, DMA for Rx and Tx.
		d := dma.DMA1
		d.EnableClock(true) // DMA clock must remain enabled in sleep mode.
		rtos.IRQ(irq.DMA1_Stream5).Enable()
		rtos.IRQ(irq.DMA1_Stream6).Enable()
		drv = i2c.NewDriver(twi, d.Channel(5, 1), d.Channel(6, 1))
		eeprom = drv.NewMasterConn(addr, i2c.ASRD)
	case 3: // AltDriver, interrupt mode.
		adrv = i2c.NewAltDriver(twi)
		adrv.SetIntMode(true)
		eeprom = adrv.NewMasterConn(addr, i2c.ASRD)
	case 4: // AltDriverDMA, interrupt mode.
		d := dma.DMA1
		d.EnableClock(true) // DMA clock must remain enabled in sleep mode.
		rtos.IRQ(irq.DMA1_Stream5).Enable()
		rtos.IRQ(irq.DMA1_Stream6).Enable()
		adrvdma = i2c.NewAltDriverDMA(twi, d.Channel(5, 1), d.Channel(6, 1))
		adrvdma.SetIntMode(true, true)
		eeprom = adrvdma.NewMasterConn(addr, i2c.ASRD)
	}
	twi.Enable()
	rtos.IRQ(irq.I2C1_EV).Enable()
	rtos.IRQ(irq.I2C1_ER).Enable()
}

func main() {
	delay.Millisec(100)

	addr := []byte{0}

	fmt.Printf("Sending data to EEPROM... ")
	_, err := eeprom.Write(addr)
	checkErr(err)
	_, err = eeprom.Write([]byte("-+Hello EEPROM+-"))
	checkErr(err)
	eeprom.StopWrite()
	fmt.Printf("OK.\n")

	fmt.Printf("Waiting for writing... ")
	for {
		_, err = eeprom.Write(addr)
		if err == nil {
			break
		}
		if e, ok := err.(i2c.Error); !ok || e != i2c.AckFail {
			checkErr(err)
		}
		eeprom.UnlockDriver()
		fmt.Printf(".")
	}
	fmt.Printf(" OK.\n")

	var buf [16]byte
	_, err = eeprom.Read(buf[:])
	checkErr(err)
	fmt.Printf("Read string: \"%s\"\n", buf[:])
	t1 := rtos.Nanosec()
	n := 0
	for {
		_, err = eeprom.Read(buf[:])
		checkErr(err)
		if n++; n == 2000 {
			t2 := rtos.Nanosec()
			fmt.Printf(
				"Read speed: %d bit/s\n",
				int64(len(buf)*n)*8e9/(t2-t1),
			)
			n = 0
			t1 = t2
		}
	}
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
