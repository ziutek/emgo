#!/bin/sh

set -e

arch=`grep 'EGARCH=' ../build.sh |sed 's/.*EGARCH=\([[:alnum:]_]\+\).*/\1/g'`
if [ -z "$arch" ]; then
	arch=$EGARCH
fi

openocd -f interface/$INTERFACE.cfg -f target/nrf51.cfg \
	-c 'telnet_port pipe' \
	-c 'init' \
	-c 'reset init' \
	-c "program $arch.elf" \
	-c 'reset run' \
	-c 'exit'
