#!/bin/sh

INTERFACE=stlink-v2
TARGET=stm32f4x
TRACECLKIN=16000000

. ../../../../../scripts/load-oocd.sh $@