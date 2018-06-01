#!/bin/sh

INTERFACE=stlink
TARGET=stm32f4x
TRACECLKIN=96000000

. ../../../../../scripts/swo-oocd.sh $@
