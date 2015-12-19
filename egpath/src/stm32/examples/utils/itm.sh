#!/bin/sh

openocd -f interface/$INTERFACE.cfg -f target/$TARGET.cfg \
	-c 'init' \
	-c "tpiu config internal /dev/stdout uart off $TRACECLKIN" \
	-c 'itm ports on' \
	|itmsplit p0:/dev/stdout /dev/stderr
