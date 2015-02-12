package main

import (
	"delay"
	"rtos"
)

func heatingTask() {
	for {
		rtos.DbgOut.WriteString("Hello debugger!\n")
		rtos.DbgErr.WriteString("Error!\n")
		heatPort.ClearAndSet(1<<(16+heat0) | 1<<heat1)
		delay.Millisec(500)
		heatPort.ClearAndSet(1<<(16+heat1) | 1<<heat0)
		delay.Millisec(500)
	}
}
