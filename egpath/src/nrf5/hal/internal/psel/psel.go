package psel

import (
	"nrf5/hal/gpio"
)

func Pin(psel uint32) gpio.Pin {
	return gpio.PortN(int(psel >> 5)).Pin(int(psel & 31))
}

func Sel(pin gpio.Pin) uint32 {
	return uint32(pin.Port().Index()*32 + pin.Index())
}
