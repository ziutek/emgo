#!/bin/sh

unifdef -k -f undef.h -D STM32F40_41xxx stm32f4xx.h |stm32xgen