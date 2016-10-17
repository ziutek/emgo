#!/bin/sh

INTERFACE=stlink-v2
TARGET=stm32f1x
TRACECLKIN=72000000

cfg='-c reset_config none separate'

. ../../utils/itm.sh
