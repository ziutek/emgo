package main

import "rtos"

func main() {
	for {
		rtos.Debug(0).WriteString("Hello world!\n")
	}
}
