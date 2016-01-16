StackLog2 = 10; /* 1 KB */
StackFrac = 8;  /* 8/8 */

/* One task, two stacks. */
MaxTasks = 1;

/* bootSRAM code should by placed at 0x200001E0. */
bootOffset = 0x1E0;

INCLUDE stm32/f103x8
INCLUDE stm32/loadram
INCLUDE noos-cortexm

/* Include bootRAM and set it as entry point. */
ENTRY(bootRAM)
