#!/bin/sh

set -e

ud='unifdef -k -f undef.h -D'

$ud STM32F40_41xxx stm32f4xx.h |stm32xgen stm32/o/f40_41xxx
$ud STM32F411xE    stm32f4xx.h |stm32xgen stm32/o/f411xe
$ud STM32L1XX_MD   stm32l1xx.h |stm32xgen stm32/o/l1xx_md
$ud STM32F10X_HD   stm32f10x.h |stm32xgen stm32/o/f10x_hd