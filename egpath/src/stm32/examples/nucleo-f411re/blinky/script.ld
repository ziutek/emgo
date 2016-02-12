StackISR = 2000; /* 2 KB */

/* Stack size: 2*6/8 KiB */
StackTaskLog2 = 11;
StackTaskFrac = 6;
MaxTasks = 1;

INCLUDE stm32/f411xe
INCLUDE stm32/loadram
INCLUDE noos-cortexm

