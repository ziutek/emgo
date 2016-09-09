package main

import (
	"fmt"
	"io"
	"rtos"

	"hdc"
	"hdc/hdcfb"
)

func checkErr(err error) {
	if err == nil {
		return
	}
	fmt.Printf("Error %v\n", err)
	for {
	}
}

var lcd *hdcfb.FB

func startLCD(rw io.ReadWriter) {
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
	go lcdLoop()
}

func lcdLoop() {
	line := lcd.NewSyncSlice(15, 20)
	for i := 0; ; i++ {
		lcd.WaitAndSwap(rtos.Nanosec() + 1e9)
		fmt.Fprintf(line, "%5d", i)
		line.SetPos(0)
		checkErr(lcd.Draw())
	}
}
