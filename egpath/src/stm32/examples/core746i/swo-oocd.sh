#!/bin/sh

INTERFACE=stlink
TARGET=stm32f7x
TRACECLKIN=192000000

. ../../../../../scripts/swo-oocd.sh $@
