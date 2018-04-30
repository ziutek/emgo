#!/bin/sh

INTERFACE=stlink-v2-1
TARGET=stm32f4x
TRACECLKIN=102000000

. ../../utils/load-oocd.sh $@
