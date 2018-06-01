#!/bin/sh

INTERFACE=stlink
TARGET=nrf52
TRACECLKIN=4000000

. ../../../../../scripts/load-oocd.sh $@
