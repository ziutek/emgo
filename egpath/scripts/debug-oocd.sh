#!/bin/sh

set -e

arch=`grep 'EGARCH=' ../build.sh |sed 's/.*EGARCH=\([[:alnum:]_]\+\).*/\1/g'`
if [ -z "$arch" ]; then
	arch=$EGARCH
fi

if [ -z "$cfg" ]; then
	cfg="'set __NOP _'"
fi

oocd_cmd="openocd -d0 -f interface/$INTERFACE.cfg -f target/$TARGET.cfg -c "$cfg" -c 'gdb_port pipe' -c 'log_output /dev/null'"

brkpnt=6
wchpnt=4

case "$arch" in
cortexm0)
	brkpnt=4
	wchpnt=2
	;;
cortexm7)
	brkpnt=8
	;;
esac

arm-none-eabi-gdb --tui \
	-ex "target extended-remote | $oocd_cmd" \
	-ex 'set mem inaccessible-by-default off' \
	-ex "set remote hardware-breakpoint-limit $brkpnt" \
	-ex "set remote hardware-watchpoint-limit $wchpnt" \
	-ex 'set history save on' \
	-ex 'set history filename ~/.gdb-history-emgo' \
	-ex 'set history size 1000' \
	$arch.elf
