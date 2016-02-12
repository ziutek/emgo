StackISR = 2000; /* 2 kB */

StackTaskLog2 = 12; /* 4 KiB */
StackTaskFrac = 6;  /* 6/8 -> 3 KiB */
MaxTasks = 3;

INCLUDE stm32/f407xg
INCLUDE stm32/loadram
INCLUDE noos-cortexm

