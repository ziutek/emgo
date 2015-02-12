#!/bin/sh

set -e

setsid openocd -f openocd.cfg >/dev/null 2>&1 </dev/null &

trap /bin/true INT

arm-none-eabi-gdb --tui \
	-ex "target extended-remote localhost:3333" \
	-ex "set remote hardware-breakpoint-limit 6" \
	-ex "set remote hardware-watchpoint-limit 4" \
	-ex "monitor reset halt" \
	main.elf

killall openocd
