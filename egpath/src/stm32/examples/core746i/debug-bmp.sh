#!/bin/sh
# gdb + Black Magic Probe

set -e


reset='monitor connect_srst enable'
if [ $# -eq 1 -a "$1" = 'noreset' ]; then
	reset=''
fi


arch=`grep 'EGARCH=' ../build.sh |sed 's/.*EGARCH=\([[:alnum:]_]\+\).*/\1/g'`
if [ -z "$arch" ]; then
	arch=$EGARCH
fi

arm-none-eabi-gdb --tui \
	-ex 'target extended-remote /dev/ttyACM0' \
	-ex "$reset" \
	-ex 'monitor swdp_scan' \
	-ex 'attach 1' \
	-ex 'set mem inaccessible-by-default off' \
	-ex 'set remote hardware-breakpoint-limit 8' \
	-ex 'set remote hardware-watchpoint-limit 4' \
	-ex 'set history save on' \
	-ex 'set history filename ~/.gdb-history-emgo' \
	-ex 'set history size 1000' \
	$arch.elf
