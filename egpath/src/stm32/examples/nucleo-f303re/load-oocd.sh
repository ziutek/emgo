#!/bin/sh

INTERFACE=stlink
TARGET=stm32f3x
TRACECLKIN=72000000

. ../../../../../scripts/load-oocd.sh $@
