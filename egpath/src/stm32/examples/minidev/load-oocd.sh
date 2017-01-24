#!/bin/sh

INTERFACE=stlink-v2
TARGET=stm32f1x
TRACECLKIN=72000000

cfg='reset_config srst_only srst_nogate connect_assert_srst'

. ../load-oocd-oversized.sh
