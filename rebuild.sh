#!/bin/bash

set -e

rm -rf egroot/pkg/* egpath/pkg/*

list="
builtin
errors

stack
sync/atomic
sync/barrier
mmio
bits

arch/cortexm
arch/cortexm/fpu
arch/cortexm/exce
arch/cortexm/systick
arch/cortexm/sleep
arch/cortexm/debug
arch/cortexm/debug/itm

syscall
runtime/noos
sync
runtime
delay
time
rtos
io
strconv
fmt
math/matrix32
math/rand

stm32/serial

stm32/stlink
stm32/f4/clock
stm32/f4/flash
stm32/f4/gpio
stm32/f4/periph
stm32/f4/setup
stm32/f4/exti
stm32/f4/irqs
stm32/f4/usarts

stm32/l1/clock
stm32/l1/flash
stm32/l1/gpio
stm32/l1/periph
stm32/l1/setup
stm32/l1/exti
stm32/l1/irqs
stm32/l1/usarts

dcf77
onewire
"

for p in $list; do 
	echo $p
	egc $p
done
