#!/bin/sh
# gdb + Black Magic Probe

set -e

arch=`grep 'EGARCH=' ../build.sh |sed 's/.*EGARCH=\([[:alnum:]_]\+\).*/\1/g'`
if [ -z "$arch" ]; then
	arch=$EGARCH
fi

arm-none-eabi-gdb --tui \
	-ex 'target extended-remote /dev/ttyACM0' \
	-ex 'monitor swdp_scan' \
	-ex 'attach 1' \
	-ex 'set mem inaccessible-by-default off' \
	-ex 'set remote hardware-breakpoint-limit 4' \
	-ex 'set remote hardware-watchpoint-limit 2' \
	$arch.elf
