package main

import (
	"delay"
	"fmt"
	"io"

	"hdc"
	"hdc/hdcfb"

	"stm32/hal/i2c"
)

func logLCDErr(err error) bool {
	o := err != nil
	if o {
		fmt.Printf("LCD: %v\n", err)
		delay.Millisec(1000)
		if _, ok := err.(i2c.Error); ok {
			resetI2C()
			i2cdrv.Unlock()
		}
	}
	return o
}

func hdcInit(d *hdc.Display) {
start:
	if logLCDErr(d.Init()) {
		goto start
	}
	if logLCDErr(d.SetDisplayMode(hdc.DisplayOn)) {
		goto start
	}
	// Backlight on.
	if logLCDErr(d.SetAUX()) {
		goto start
	}

}

var (
	lcd      *hdcfb.FB
	tempResp = make(chan float32, 0)
)

func startLCD(rw io.ReadWriter) {
	d := new(hdc.Display)
	*d = hdc.Display{
		ReadWriter: rw,
		Cols:       20,
		Rows:       4,
		DS:         4,
		RS:         1 << 0, RW: 1 << 1, E: 1 << 2, AUX: 1 << 3,
	}
	hdcInit(d)
	lcd = hdcfb.NewFB(d)
	go lcdLoop()
}

func lcdLoop() {
	for {
		lcd.WaitAndSwap(0)
		if logLCDErr(lcd.Draw()) {
			hdcInit(lcd.Display())
		}
	}
}
