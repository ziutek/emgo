#!/bin/sh

INTERFACE=stlink-v2
TARGET=stm32f1x
TRACECLKIN=72000000

cfg='reset_config none separate' # You need to press reset before connect.
#cfg='reset_config srst_only srst_nogate connect_assert_srst'

#. ../../utils/load-oocd.sh
. ../load-oocd-oversized.sh
