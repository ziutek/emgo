#!/bin/sh

set -e

arch=`grep 'EGARCH=' ../build.sh |sed 's/.*EGARCH=\([[:alnum:]_]\+\).*/\1/g'`
if [ -z "$arch" ]; then
	arch=$EGARCH
fi

setsid st-util >/dev/null 2>&1 </dev/null &

trap /bin/true INT

arm-none-eabi-gdb --tui \
	-ex 'target extended-remote localhost:4242' \
	-ex 'set remote hardware-breakpoint-limit 6' \
	-ex 'set remote hardware-watchpoint-limit 4' \
	-ex 'set mem inaccessible-by-default off' \
	$arch.elf

killall st-util
