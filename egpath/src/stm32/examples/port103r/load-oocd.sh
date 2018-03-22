#!/bin/sh

INTERFACE=stlink-v2
TARGET=stm32f1x
TRACECLKIN=72000000

. ../../../../../scripts/load-oocd.sh $@
