ISRStack = 2048;
MainStack = 4096;
TaskStack = 0;
MaxTasks = 1;

INCLUDE stm32/f303xe
INCLUDE stm32/loadram
INCLUDE noos-cortexm

/* 
To run from RAM the BOOT0 pin must be high and nBOOT1 option bit must be
cleared. To run from Flash change stm32/loadram to stm32/loadflash.
*/
