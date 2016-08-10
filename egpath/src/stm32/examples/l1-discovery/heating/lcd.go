package main

import (
	"fmt"
	"io"

	"hdc"
)

var lcd = &hdc.Display{
	Cols: 20, Rows: 4,
	DS: 4,
	RS: 1 << 0, RW: 1 << 1, E: 1 << 2, AUX: 1 << 3,
}

func checkErr(err error) {
	if err != nil {
		fmt.Printf("Error %v\n", err)
	}
}

func initLCD(rw io.ReadWriter) {
	lcd.ReadWriter = rw
	checkErr(lcd.Init())
	checkErr(lcd.SetDisplayMode(hdc.DisplayOn))
	checkErr(lcd.SetAUX()) // Backlight on.

}
