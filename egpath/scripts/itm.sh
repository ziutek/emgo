#!/bin/sh

itmsplit=cat

tpiu="tpiu config external uart off $TRACECLKIN 2000000"
if [ "$INTERFACE" = 'stlink-v2' -o "$INTERFACE" = 'stlink-v2-1' ]; then
	tpiu="tpiu config internal /dev/stdout uart off $TRACECLKIN"
	itmsplit='itmsplit p0:/dev/stdout /dev/stderr'
	exit=''
fi

openocd -f interface/$INTERFACE.cfg -f target/$TARGET.cfg $cfg \
	-c 'init' \
	-c "tpiu config internal /dev/stdout uart off $TRACECLKIN" \
	-c 'itm ports on' \
	|$itmsplit