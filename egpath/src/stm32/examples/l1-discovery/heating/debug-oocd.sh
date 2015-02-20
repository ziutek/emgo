#!/bin/sh

arm-none-eabi-gdb --tui \
	-ex "target extended-remote | openocd -f openocd.cfg" \
	-ex "set remote hardware-breakpoint-limit 6" \
	-ex "set remote hardware-watchpoint-limit 4" \
	-ex "monitor halt" \
	main.elf
#	-ex "monitor gdb_sync" \
#	-ex "stepi" \
