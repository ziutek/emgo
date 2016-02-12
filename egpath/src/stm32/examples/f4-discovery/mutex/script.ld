StackISR = 2000; /* 2 kB */

/* Stack size: 7/8 KiB */
StackTaskLog2 = 10;
StackTaskFrac = 7;
MaxTasks = 5;

INCLUDE stm32/f407xg
INCLUDE stm32/loadram
INCLUDE noos-cortexm

