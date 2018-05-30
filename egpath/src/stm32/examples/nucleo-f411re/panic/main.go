package main

import "rtos"

func main() {
	if rtos.Nanosec() == 32768 {
		panic("bla")
	}
	/*if rtos.Nanosec() == 32768*2 {
		panic("eee")
	}*/
}
