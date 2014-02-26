#!/bin/bash

set -e

egc runtime
egc sync

egc math/matrix32

egc delay
egc cortexm/startup
egc cortexm/nvic
egc cortexm/systick

egc stm32/f4/clock
egc stm32/f4/flash
egc stm32/f4/gpio
egc stm32/f4/periph
egc stm32/f4/setup

egc stm32/l1/clock
egc stm32/l1/flash
egc stm32/l1/gpio
egc stm32/l1/periph
egc stm32/l1/setup
