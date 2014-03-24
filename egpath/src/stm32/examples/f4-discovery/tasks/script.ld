INCLUDE stm32/f407_g
INCLUDE stm32/loadram
INCLUDE arch/cortexm
INCLUDE os/noos

/*
MaxTasks defines the maximum number of tasks.

MaxTasks == 0 means:
- no tasks (gorutines),
- one stack (MSP),
- runtime doesn't touch any Cortex-M peripherals,
- program runs in privileged mode.

MaxTasks > 0 mens:
- separate stack for any task,
- MPU (if avilable) is used to protect the stacks and runtime data,
- all tasks runs in user mode.
*/
MaxTasks = 4;

/* StackSize defines the encoded size of stack for one task. */
stackExp = 10; /* 1 KB */
stackNum = 7;  /* 7/8 */
StackSize = (stackExp - 1) << 8 | (stackNum - 1) << 5;

/* 
MaxIRQ defines the maximum number of exceptions as power of 2.

MaxIRQ should be >= 5 (2^5 allows to handle 16 system exceptions and 16
external intnerrupts.)
*/
MaxIRQ = 5;
