INCLUDE stm32/f407_g
INCLUDE stm32/loadram
INCLUDE cortexm

HeapSize = ALIGN(FreeSize / 3 , 4);
MaxTasks = 7;