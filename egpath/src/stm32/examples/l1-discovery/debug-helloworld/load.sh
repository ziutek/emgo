#!/bin/sh

INTERFACE=stlink-v2
TARGET=stm32l1
TRACECLKIN=2097000

. ../../utils/load-oocd.sh $@
