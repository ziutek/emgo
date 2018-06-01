#!/bin/sh

INTERFACE=stlink
TARGET=stm32l4x
TRACECLKIN=80000000
#TRACECLKIN=48000000

. ../../../../../scripts/load-oocd.sh $@
