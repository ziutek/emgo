#!/bin/sh

set -e

cd ../../stm32/xgen
rm -rf stm32/o

ud='unifdef -k -f undef.h -D'

$ud STM32F40_41xxx stm32f4xx.h |stm32xgen stm32/o/f40_41xxx
$ud STM32F411xE    stm32f4xx.h |stm32xgen stm32/o/f411xe
#$ud STM32F303xE    stm32f30x.h |stm32xgen stm32/o/f303xe
$ud STM32L1XX_MD   stm32l1xx.h |stm32xgen stm32/o/l1xx_md
$ud STM32F10X_MD   stm32f10x.h |stm32xgen stm32/o/f10x_md
$ud STM32F10X_HD   stm32f10x.h |stm32xgen stm32/o/f10x_hd

stm32xgen stm32/o/f303xe <stm32f3xx/stm32f303xe.h
stm32xgen stm32/o/f746xx <stm32f7xx/stm32f746xx.h
stm32xgen stm32/o/l476xx <stm32l4xx/stm32l476xx.h

cd stm32/o
for target in *; do
	cd $target
	for pkg in *; do
		cd $pkg
		pwd
		xgen *.go
		cd ..
	done
	cd ..
done

cd ../../..
rm -rf o
mv xgen/stm32/o .
