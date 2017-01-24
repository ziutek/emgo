#!/bin/sh

set -e

arch=`grep 'EGARCH=' ../build.sh |sed 's/.*EGARCH=\([[:alnum:]_]\+\).*/\1/g'`
if [ -z "$arch" ]; then
	arch=$EGARCH
fi

if [ $# -eq 1 -a "$1" = 'flash' ]; then
	load="program $arch.elf"
else
	arm-none-eabi-objcopy -O binary $arch.elf $arch.bin
	load="load_image $arch.bin 0x20000000"
fi

if [ -n "$TRACECLKIN" ]; then
	tpiu="tpiu config internal /dev/stdout uart off $TRACECLKIN"
	itm='itm ports on'
fi

echo CFG: $cfg

echo "Loading at $addr..." >/dev/stderr
openocd -f interface/$INTERFACE.cfg -f target/$TARGET.cfg -c "$cfg" \
	-c 'init' \
	-c 'reset init' \
	-c "$load" \
	-c "$tpiu" \
	-c "$itm" \
	-c 'reset run' \
	|itmsplit p0:/dev/stdout /dev/stderr
