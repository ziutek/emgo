#!/bin/sh

arm-none-eabi-objcopy -O binary -R .noload $EGARCH.elf $EGARCH.bin
addr=0x20000000
if [ $# -eq 1 -a "$1" = 'flash' ]; then
	addr=0x8000000
fi

if [ -n "$TRACECLKIN" ]; then
	tpiu="tpiu config internal /dev/stdout uart off $TRACECLKIN"
	itm='itm ports on'
fi

echo "Loading at $addr..." >/dev/stderr
openocd -f interface/$INTERFACE.cfg -f target/$TARGET.cfg \
	-c 'init' \
	-c 'reset init' \
	-c "load_image $EGARCH.bin $addr" \
	-c "$tpiu" \
	-c "$itm" \
	-c 'reset run' \
	|itmsplit p0:/dev/stdout /dev/stderr
