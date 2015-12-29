#!/bin/sh


if [ $# -eq 1 -a "$1" = 'flash' ]; then
	load="program $EGARCH.elf"
else
	arm-none-eabi-objcopy -O binary $EGARCH.elf $EGARCH.bin
	load="load_image $EGARCH.bin 0x20000000"
fi

if [ -n "$TRACECLKIN" ]; then
	tpiu="tpiu config internal /dev/stdout uart off $TRACECLKIN"
	itm='itm ports on'
fi

echo "Loading at $addr..." >/dev/stderr
openocd -f interface/$INTERFACE.cfg -f target/$TARGET.cfg \
	-c 'init' \
	-c 'reset init' \
	-c "$load" \
	-c "$tpiu" \
	-c "$itm" \
	-c 'reset run' \
	|itmsplit p0:/dev/stdout /dev/stderr
