package main

import (
	"delay"
	"fmt"
	"io"
	"text/linewriter"

	"sdcard"
)

func print(s string) {
	io.WriteString(fmt.DefaultWriter, s)
	if s[len(s)-1] == '\n' {
		return
	}
	fmt.DefaultWriter.(*linewriter.Writer).Flush()
}

func printOK() {
	print(" OK\n")
}

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
	print(" retry timeout\n")
	for {
		led.Clear()
		delay.Millisec(200)
		led.Set()
		delay.Millisec(200)
	}
}

func sendCMD52(h sdcard.Host, fn, addr int, flags sdcard.IORWFlags, val byte) byte {
	val, st := h.SendCmd(sdcard.CMD52(fn, addr, flags, val)).R5()
	checkErr("CMD52", h.Err(true), st)
	return val
}
