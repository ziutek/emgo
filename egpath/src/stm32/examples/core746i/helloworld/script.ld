ISRStack = 64 * 32;
MainStack = 128 * 32;
TaskStack = 0;
MaxTasks = 1;

INCLUDE stm32/f7x6xg
INCLUDE stm32/loadram
INCLUDE noos-cortexm

/*
Use Core746I BOOT CONFIG switch and STM32F746I BOOT_ADDx registers to select
boot addres required by stm32/load* script. Read ../boot_config.txt for more
information.
*/
