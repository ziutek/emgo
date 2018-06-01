#!/bin/sh

INTERFACE=stlink
TARGET=stm32f0x

#cfg='reset_config none separate' # Press reset before connect.
#cfg='reset_config srst_only srst_nogate connect_assert_srst'

. ../../../../../scripts/debug-oocd.sh $@
