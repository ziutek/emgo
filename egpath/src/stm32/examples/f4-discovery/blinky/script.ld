INCLUDE stm32/f407_g
INCLUDE stm32/loadram
INCLUDE arch/cortexm
INCLUDE os/noos

stackExp = 10; /* 1 KB */
stackNum = 7;  /* 7/8 */
StackSize = (stackExp - 1) << 8 | (stackNum - 1) << 5;
