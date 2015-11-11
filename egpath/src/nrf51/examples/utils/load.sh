#!/bin/sh

set -e

arm-none-eabi-objcopy -O binary -R .noload main.elf main.bin

openocd -f interface/$INTERFACE.cfg -f target/nrf51.cfg \
	-c 'telnet_port pipe' \
	-c 'init' \
	-c 'reset init' \
	-c 'program main.bin' \
	-c 'reset run' \
	-c 'exit'
