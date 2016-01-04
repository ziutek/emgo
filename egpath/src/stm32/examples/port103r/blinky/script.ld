StackLog2 = 12; /* 4 KB */
StackFrac = 8;  /* 8/8 */

/* One task, two stacks. */
MaxTasks = 1;

/* bootSRAM code should by placed at 0x200001E0. */
bootOffset = 0x1E0;

INCLUDE stm32/f103rc
INCLUDE stm32/loadram
INCLUDE noos-cortexm

/* Include bootSRAM and set it as entry point. */
ENTRY(bootSRAM)