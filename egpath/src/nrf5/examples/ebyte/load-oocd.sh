#!/bin/sh

INTERFACE=stlink-v2
TARGET=nrf52
TRACECLKIN=64000000
#. ../../../../stm32/examples/utils/load-oocd.sh $@
. ../../utils/load-oocd.sh $@
