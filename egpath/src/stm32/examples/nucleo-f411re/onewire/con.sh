#!/bin/sh
set -e
tty=/dev/ttyACM0
speed=115200
stty -F $tty ispeed $speed ospeed $speed igncr
exec cat $tty
