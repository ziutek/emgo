#!/bin/sh

INTERFACE=stlink
TARGET=stm32f4x
TRACECLKIN=168000000

. ../../../../../scripts/load-oocd.sh $@
