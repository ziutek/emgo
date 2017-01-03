package main

import (
	"delay"

	"stm32/hal/gpio"
	"stm32/hal/spi"
)

type ILI9341 struct {
	spi           *spi.Driver
	reset, cs, dc gpio.Pin
}

func (ili *ILI9341) Reset() {
	ili.cs.Set()
	ili.reset.Clear()
	delay.Millisec(1)
	ili.reset.Set()
}

func (ili *ILI9341) Select() {
	ili.cs.Clear()
}

func (ili *ILI9341) Deselect() {
	ili.cs.Set()
}

func (ili *ILI9341) Cmd(cmd byte) {
	ili.dc.Clear()
	ili.spi.WriteReadByte(cmd)
	ili.dc.Set()
}

func (ili *ILI9341) WriteByte(b byte) {
	ili.spi.WriteReadByte(b)
}

func (ili *ILI9341) Write(data []byte) {
	ili.spi.WriteRead(data, nil)
}

func (ili *ILI9341) Cmd16(cmd uint16) {
	ili.dc.Clear()
	ili.spi.WriteReadWord16(cmd)
	ili.dc.Set()
}

func (ili *ILI9341) WriteWord16(w uint16) {
	ili.spi.WriteReadWord16(w)
}

func (ili *ILI9341) Write16(data []uint16) {
	ili.spi.WriteRead16(data, nil)
}

func (ili *ILI9341) Fill16(w uint16, n int) {
	ili.spi.RepeatWord16(w, n)
}

const (
	ILI9341_SLPOUT  = 0x11
	ILI9341_DISPOFF = 0x28
	ILI9341_DISPON  = 0x29
	ILI9341_RAMWR   = 0x2C
	ILI9341_MADCTL  = 0x36
	ILI9341_PIXFMT  = 0x3A
	ILI9341_CASET   = 0x2A
	ILI9341_PASET   = 0x2B
)
