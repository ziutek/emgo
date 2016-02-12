StackISR = 2000;    /* 2 kB */

/* Stack size: 7/8 KB */
StackTaskLog2 = 10;
StackTaskFrac = 7;

/* One task, two stacks. */
MaxTasks = 1;

INCLUDE stm32/f407xg
INCLUDE stm32/loadram
INCLUDE noos-cortexm

