#!/bin/sh

INTERFACE=stlink-v2-1
TARGET=stm32l4x
#TRACECLKIN=800000000
TRACECLKIN=480000000

. ../../utils/load-oocd.sh $@
