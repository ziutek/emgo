#!/bin/sh

set -e

oocd_cmd="openocd -f interface/$INTERFACE.cfg -f target/nrf51.cfg -c 'gdb_port pipe' -c 'log_output /dev/null'"

arm-none-eabi-gdb --tui \
	-ex "target extended-remote | $oocd_cmd" \
	-ex "set remote hardware-breakpoint-limit 4" \
	-ex "set remote hardware-watchpoint-limit 2" \
	main.elf
