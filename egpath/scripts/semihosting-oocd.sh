#!/bin/sh

openocd -d0 -f interface/$INTERFACE.cfg -f target/$TARGET.cfg \
	-c 'init' \
	-c 'arm semihosting enable' \
	-c 'reset run'
