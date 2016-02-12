StackISR = 2000; /* 2 KB */

StackTaskLog2 = 10; /* 1 KiB */
StackTaskFrac = 8;  /* 8/8 */
MaxTasks = 1;

/* bootSRAM code should by placed at 0x200001E0. */
/*bootOffset = 0x1E0;*/

INCLUDE stm32/f103x8
INCLUDE stm32/loadflash
INCLUDE noos-cortexm

/* Include bootRAM and set it as entry point. */
/*ENTRY(bootRAM)*/
