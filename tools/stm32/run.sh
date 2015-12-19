#!/bin/sh

unifdef -k -f undef.h -D STM32F40_41xxx stm32f4xx.h |stm32xgen f40_41xxx
unifdef -k -f undef.h -D STM32F411xE stm32f4xx.h |stm32xgen f411xe