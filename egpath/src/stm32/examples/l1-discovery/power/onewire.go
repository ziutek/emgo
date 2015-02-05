package main

import (
	"stm32/onedrv"
	"stm32/serial"

	"onewire"
)

var (
	onewSerial = serial.New(onewUART, 8, 8)
	onewDriver = onedrv.UARTDriver{Serial: onewSerial, Clock: onewClk}
	onewMaster = onewire.Master{Driver: &onewDriver}
)
