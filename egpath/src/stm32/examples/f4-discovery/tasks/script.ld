INCLUDE stm32/f407_g
INCLUDE stm32/loadram
INCLUDE arch/cortexm
INCLUDE os/noos

HeapSize = ALIGN(HeapStackSize / 3 , 4);
MaxTasks = 7;