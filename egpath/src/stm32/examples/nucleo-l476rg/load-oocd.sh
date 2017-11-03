#!/bin/sh

INTERFACE=stlink-v2-1
TARGET=stm32l4x
TRACECLKIN=80000000
#TRACECLKIN=48000000

. ../../utils/load-oocd.sh $@
