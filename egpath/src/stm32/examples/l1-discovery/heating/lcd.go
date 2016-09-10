package main

import (
	"delay"
	"fmt"
	"io"
	"rtos"

	"hdc"
	"hdc/hdcfb"
	"onewire"

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
	lcd        *hdcfb.FB
	lcdHead    *hdcfb.SyncSlice
	searchResp = make(chan onewire.Dev, 1)
	tempResp   = make(chan float32, 1)
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
	lcdHead = lcd.NewSyncSlice(0, 20)
	t := rtos.Nanosec()
	var lastPrint int64
	for {
		lcd.WaitAndSwap(t + 2e9)
		t = rtos.Nanosec()
		if t-lastPrint >= 2e9 {
			lastPrint = t
			owd.Cmd <- SearchCmd{Typ: onewire.DS18B20, Resp: searchResp}
			var (
				devs [2]onewire.Dev
				i    int
			)
			for d := range searchResp {
				if d.Type() == 0 {
					break
				}
				devs[i] = d
				i++
			}
			var temps [len(devs)]float32
			for i, d := range devs {
				if d.Type() == 0 {
					temps[i] = -99
					continue
				}
				owd.Cmd <- TempCmd{Dev: d, Resp: tempResp}
				temps[i] = <-tempResp
			}
			fmt.Fprintf(lcdHead, "%8d ", t/1e9)
			for _, temp := range temps {
				fmt.Fprintf(lcdHead, "%2.1f ", temp)
			}
			lcdHead.SetPos(0)
		}
		if logLCDErr(lcd.Draw()) {
			hdcInit(lcd.Display())
		}
	}
}
