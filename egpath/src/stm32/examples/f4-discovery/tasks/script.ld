/*
MaxTasks defines the maximum number of tasks.

MaxTasks == 0 means:
- no tasks (gorutines),
- one stack (uses MSP),
- runtime doesn't touch any Cortex-M peripherals,
- program runs in privileged mode.

MaxTasks > 0 mens:
- separate stack for any task + separate stack for ISRs,
- MPU (if avilable) is used to protect the stacks and other areas,
- all tasks runs in user mode.
*/
MaxTasks = 4;

/* StackExp and StackFrac defines the size of the stack for one task. */
StackExp = 11; /* 2 KB */
StackFrac = 7;  /* 7/8 */

/* 
VTExp defines the length of the table for exception vectors. VTExp should be
>= 5 (2^5 allows to handle 16 system exceptions and 16 external intnerrupts.)
*/
VTExp = 5;

INCLUDE stm32/f407_g
INCLUDE stm32/loadram
INCLUDE noos-cortexm

