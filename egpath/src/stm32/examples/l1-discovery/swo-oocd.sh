#!/bin/sh

INTERFACE=stlink-v2
TARGET=stm32l1
TRACECLKIN=32000000

. ../../../../../scripts/swo-oocd.sh $@
