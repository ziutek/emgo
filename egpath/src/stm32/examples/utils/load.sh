#!/bin/sh

set -e
 
arch=`grep 'EGARCH=' ../build.sh |sed 's/.*EGARCH=\([[:alnum:]_]\+\).*/\1/g'`
if [ -z "$arch" ]; then
	arch=$EGARCH
fi

arm-none-eabi-objcopy -O binary $arch.elf $arch.bin
addr=0x20000000
if [ $# -eq 1 -a "$1" = 'flash' ]; then
	addr=0x8000000
fi       
echo "Loading at $addr..."
st-flash --reset write $arch.bin $addr

