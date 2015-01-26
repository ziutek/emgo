#!/bin/sh

set -e

setsid openocd -f oocd-stlink.cfg >/dev/null 2>&1 </dev/null &

trap /bin/true INT

arm-none-eabi-gdb --tui \
	-ex "target extended-remote localhost:3333" \
	-ex "set remote hardware-breakpoint-limit 4" \
	-ex "set remote hardware-watchpoint-limit 2" \
	main.elf

killall openocd
