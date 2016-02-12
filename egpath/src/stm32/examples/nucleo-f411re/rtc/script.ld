StackISR = 2000; /* 2 KB */

/* Stack size: 6 KiB */
StackTaskLog2 = 13; /* 8 KiB */
StackTaskFrac = 6;  /* 6/8 */

/* Two gorutines, 3 stacks. */
MaxTasks = 2;

INCLUDE stm32/f411xe
INCLUDE stm32/loadram
INCLUDE noos-cortexm

