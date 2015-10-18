#!/bin/sh

set -e

setsid st-util >/dev/null 2>&1 </dev/null &

trap /bin/true INT

arm-none-eabi-gdb --tui \
	-ex "target extended-remote localhost:4242" \
	-ex "set remote hardware-breakpoint-limit 6" \
	-ex "set remote hardware-watchpoint-limit 4" \
	main.elf

killall st-util
