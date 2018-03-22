#!/bin/sh

INTERFACE=stlink-v2-1
TARGET=stm32f4x
TRACECLKIN=96000000

. ../../../../../scripts/load-oocd.sh $@
