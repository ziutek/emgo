#!/bin/sh

arm-none-eabi-objcopy -Obinary main.elf main.bin
addr=0x20000000
if [ $# -eq 1 -a "$1" = 'flash' ]; then
	addr=0x8000000
fi       
echo "Loading at $addr..."
st-flash --reset write main.bin $addr
