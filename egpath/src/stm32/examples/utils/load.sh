#!/bin/sh

#arm-none-eabi-size $EGARCH.elf >$EGARCH.sizes
#arm-none-eabi-objcopy -O binary -R .noload $EGARCH.elf $EGARCH.bin

arm-none-eabi-objcopy -O binary $EGARCH.elf $EGARCH.bin
addr=0x20000000
if [ $# -eq 1 -a "$1" = 'flash' ]; then
	addr=0x8000000
fi       
echo "Loading at $addr..."
st-flash --reset write $EGARCH.bin $addr
