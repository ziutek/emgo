#!/bin/sh

INTERFACE=stlink-v2
TARGET=stm32f1x
TRACECLKIN=72000000

cfg='set _FLASHNAME oversized.flash;flash bank $_FLASHNAME stm32f1x 0 0x20000 0 0 $_TARGETNAME'
cfg="$cfg;reset_config none separate"
#cfg="$cfg;reset_config srst_only srst_nogate connect_assert_srst"

. ../../../../../scripts/load-oocd.sh $@
