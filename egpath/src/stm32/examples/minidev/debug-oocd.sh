#!/bin/sh

INTERFACE=stlink-v2
TARGET=stm32f1x

cfg='-c reset_config none separate'

. ../../utils/debug-oocd.sh
