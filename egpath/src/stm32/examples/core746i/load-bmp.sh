#!/bin/sh
# gdb + Black Magic Probe

set -e

arch=`grep 'EGARCH=' ../build.sh |sed 's/.*EGARCH=\([[:alnum:]_]\+\).*/\1/g'`
if [ -z "$arch" ]; then
	arch=$EGARCH
fi

arm-none-eabi-gdb \
	-ex 'set pagination off' \
	-ex 'set confirm off' \
	-ex 'target extended-remote /dev/ttyACM0' \
	-ex 'monitor connect_srst enable' \
	-ex 'monitor swdp_scan' \
	-ex 'attach 1' \
	-ex 'load' \
	-ex 'run' \
	-ex 'quit' \
	$arch.elf
