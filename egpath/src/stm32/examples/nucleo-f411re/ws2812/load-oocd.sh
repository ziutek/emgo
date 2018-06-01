#!/bin/sh

INTERFACE=stlink
TARGET=stm32f4x
TRACECLKIN=102000000

. ../../utils/load-oocd.sh $@
