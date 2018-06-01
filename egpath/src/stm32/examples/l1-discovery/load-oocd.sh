#!/bin/sh

INTERFACE=stlink
TARGET=stm32l1
TRACECLKIN=32000000

. ../../../../../scripts/load-oocd.sh $@
