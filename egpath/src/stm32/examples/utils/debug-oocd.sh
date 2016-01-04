#!/bin/sh

set -e

oocd_cmd="openocd -f interface/$INTERFACE.cfg -f target/$TARGET.cfg -c 'gdb_port pipe' -c 'log_output /dev/null'"

arm-none-eabi-gdb --tui \
	-ex "target extended-remote | $oocd_cmd" \
	-ex 'set remote hardware-breakpoint-limit 6' \
	-ex 'set remote hardware-watchpoint-limit 4' \
	-ex 'set mem inaccessible-by-default off' \
	-ex 'monitor reset init' \
	$EGARCH.elf
