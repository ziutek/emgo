#!/bin/sh

set -e

arch=`grep 'EGARCH=' ../build.sh |sed 's/.*EGARCH=\([[:alnum:]_]\+\).*/\1/g'`
if [ -z "$arch" ]; then
	arch=$EGARCH
fi

oocd_cmd="openocd -f interface/$INTERFACE.cfg -f target/$TARGET.cfg $cfg -c 'gdb_port pipe' -c 'log_output /dev/null'"

arm-none-eabi-gdb --tui \
	-ex "target extended-remote | $oocd_cmd" \
	-ex 'set remote hardware-breakpoint-limit 6' \
	-ex 'set remote hardware-watchpoint-limit 4' \
	-ex 'set mem inaccessible-by-default off' \
	$arch.elf
