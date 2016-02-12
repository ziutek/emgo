StackISR = 2000; /* 2 KB */

StackTaskLog2 = 12; /* 4 KiB */
StackTaskFrac = 8;  /* 8/8 */

/* Two task, three stacks. */
MaxTasks = 2;

/* bootSRAM code should by placed at 0x200001E0. */
bootOffset = 0x1E0;

INCLUDE stm32/f103xc
INCLUDE stm32/loadram
INCLUDE noos-cortexm

/* Include bootRAM and set it as entry point. */
ENTRY(bootRAM)
