#!/bin/sh
# gdb + Black Magic Probe

set -e

arm-none-eabi-gdb --tui \
	-ex 'target extended-remote /dev/ttyACM0' \
	-ex 'monitor swdp_scan' \
	-ex 'attach 1' \
	-ex 'set mem inaccessible-by-default off' \
	-ex 'set remote hardware-breakpoint-limit 4' \
	-ex 'set remote hardware-watchpoint-limit 2' \
	$EGARCH.elf
