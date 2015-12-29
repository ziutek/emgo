#!/bin/sh

set -e

arm-none-eabi-objcopy -O binary $EGARCH.elf $EGARCH.bin

openocd -f interface/$INTERFACE.cfg -f target/nrf51.cfg \
	-c 'telnet_port pipe' \
	-c 'init' \
	-c 'reset init' \
	-c "program $EGARCH.bin" \
	-c 'reset run' \
	-c 'exit'
