package main

import (
	"fmt"
	"io"

	"hdc"
	"hdc/hdcfb"
)

var lcd *hdcfb.FB

func initLCD(rw io.ReadWriter) {
	d := new(hdc.Display)
	*d = hdc.Display{
		ReadWriter: rw,
		Cols:       20,
		Rows:       4,
		DS:         4,
		RS:         1 << 0, RW: 1 << 1, E: 1 << 2, AUX: 1 << 3,
	}
	checkErr(d.Init())
	checkErr(d.SetDisplayMode(hdc.DisplayOn))
	checkErr(d.SetAUX()) // Backlight on.

	lcd = hdcfb.NewFB(d)
}

func checkErr(err error) {
	if err == nil {
		return
	}
	fmt.Printf("Error %v\n", err)
	for {
	}
}
