#!/bin/sh

INTERFACE=stlink-v2
TARGET=stm32f1x

#cfg='reset_config none separate'
cfg='reset_config srst_only srst_nogate connect_assert_srst'

. ../../utils/debug-oocd.sh
