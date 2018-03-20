#!/bin/sh

set -e

arch=`grep 'EGARCH=' ../build.sh |sed 's/.*EGARCH=\([[:alnum:]_]\+\).*/\1/g'`
if [ -z "$arch" ]; then
	arch=$EGARCH
fi

if grep -q '^INCLUDE.*/loadflash' script.ld; then
	load="program $arch.elf"
else
	arm-none-eabi-objcopy -O binary $arch.elf $arch.bin
	load="load_image $arch.bin 0x20000000"
fi

itmsplit=cat
exit=exit

if [ -n "$TRACECLKIN" ]; then
	tpiu="tpiu config external uart off $TRACECLKIN 2000000"
	if [ "$INTERFACE" = 'stlink-v2' -o "$INTERFACE" = 'stlink-v2-1' ]; then
		tpiu="tpiu config internal /dev/stdout uart off $TRACECLKIN"
		itmsplit='itmsplit p0:/dev/stdout /dev/stderr'
		exit=''
	fi
	itm='itm ports on'
fi

echo CFG: $cfg

openocd -f interface/$INTERFACE.cfg -f target/$TARGET.cfg -c "$cfg" \
	-c 'init' \
	-c 'reset init' \
	-c "$load" \
	-c "$tpiu" \
	-c "$itm" \
	-c 'reset run' \
	-c "$exit" \
	|$itmsplit
