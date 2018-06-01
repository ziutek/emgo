#!/bin/sh

INTERFACE=stlink
TARGET=stm32f1x
TRACECLKIN=72000000

. ../../../../../scripts/load-oocd.sh $@
