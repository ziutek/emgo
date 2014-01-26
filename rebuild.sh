#!/bin/bash

set -e

rebuild() {
	cd $1
	make clean >/dev/null 2>&1
	make
	cd - >/dev/null
}


cd egroot/src/pkg

rebuild delay
rebuild stm32/clock
rebuild stm32/flash
rebuild stm32/gpio
rebuild stm32/periph


cd ../../../egpath/src

rebuild blinky

