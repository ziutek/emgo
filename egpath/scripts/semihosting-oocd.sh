#!/bin/sh

openocd -f interface/$INTERFACE.cfg -f target/$TARGET.cfg \
	-c 'init' \
	-c 'arm semihosting enable' \
	-c 'reset run'
