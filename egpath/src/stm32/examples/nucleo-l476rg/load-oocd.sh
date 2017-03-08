#!/bin/sh

INTERFACE=stlink-v2-1
TARGET=stm32l4x
TRACECLKIN=40000000

. ../../utils/load-oocd.sh $@
