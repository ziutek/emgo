#!/bin/sh

INTERFACE=stlink-v2-1
TARGET=stm32f1x
TRACECLKIN=72000000

. ../../utils/load-oocd.sh $@
