#!/bin/sh
set -e
tty=/dev/ttyACM0
#stty -F $tty 115200
exec cat $tty
