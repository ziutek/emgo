package main

import (
	"delay"
	"fmt"

	"sdcard"
)

func checkErr(what string, err error, status sdcard.IOStatus) {
	switch {
	case err != nil:
		fmt.Printf(" %s: %v\n", what, err)
	case status&^sdcard.IO_CURRENT_STATE != 0:
		fmt.Printf(" %s: 0x%X", what, status)
	default:
		return
	}
	for {
		led.Clear()
		delay.Millisec(200)
		led.Set()
		delay.Millisec(200)
	}
}

func checkRetry(retry int) {
	if retry > 0 {
		return
	}
	fmt.Printf(" retry timeout")
	for {
		led.Clear()
		delay.Millisec(200)
		led.Set()
		delay.Millisec(200)
	}
}

func printOK() {
	fmt.Printf(" OK\n")
}

func sendCMD52(h sdcard.Host, fn, addr int, flags sdcard.IORWFlags, val byte) byte {
	val, st := h.SendCmd(sdcard.CMD52(fn, addr, flags, val)).R5()
	checkErr("CMD52", h.Err(true), st)
	return val
}
