#!/bin/sh

INTERFACE=stlink
TARGET=stm32f2x
TRACECLKIN=168000000

. ../../../../../scripts/load-oocd.sh $@
