package main

import (
	"stm32/hal/dma"
	"stm32/hal/i2c"
)

var i2cdrv *i2c.Driver

func initI2C(twi *i2c.Periph, rxdma, txdma *dma.Channel) {
	twi.EnableClock(true)
	twi.Reset() // Mandatory!
	twi.Setup(&i2c.Config{Speed: 100e3})
	i2cdrv = i2c.NewDriver(twi, rxdma, txdma)
	twi.Enable()
}

func i2cISR() {
	i2cdrv.I2CISR()
}

func i2cRxDMAISR() {
	i2cdrv.DMAISR(i2cdrv.RxDMA)
}

func i2cTxDMAISR() {
	i2cdrv.DMAISR(i2cdrv.TxDMA)
}
