#!/bin/sh

itmsplit=cat

tpiu="tpiu config external uart off $TRACECLKIN 2000000"
if [ "$INTERFACE" = 'stlink-v2' -o "$INTERFACE" = 'stlink-v2-1' ]; then
	# Reduce speed to 200 kb/s because of problems with some ST-LINK
	tpiu="tpiu config internal /dev/stdout uart off $TRACECLKIN 200000"
	itmsplit='itmsplit p0:/dev/stdout /dev/stderr'
	exit=''
fi

openocd -f interface/$INTERFACE.cfg -f target/$TARGET.cfg $cfg \
	-c 'init' \
	-c "$tpiu" \
	-c 'itm ports on' \
	|$itmsplit