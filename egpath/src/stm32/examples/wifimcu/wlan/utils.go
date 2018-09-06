package main

import (
	"delay"
	"fmt"
	"io"
	"text/linewriter"
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

func checkErr(what string, err error) {
	if err == nil {
		return
	}
	fmt.Printf(" %s: %v\n", what, err)
	for {
		led.Clear()
		delay.Millisec(200)
		led.Set()
		delay.Millisec(200)
	}
}

/*
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
*/
