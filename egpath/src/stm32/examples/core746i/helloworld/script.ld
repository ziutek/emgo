ISRStack = 64 * 32;
MainStack = 128 * 32;
TaskStack = 0;
MaxTasks = 1;

/*
To boot from RAM BOOT_ADDR1 default value 0x0040 must be changed to to 0x8000.
Use oppenocd commands: stm32f2x options_read 0, stm32f2x options_write 0.
*/
INCLUDE stm32/f407xg
INCLUDE stm32/loadram
INCLUDE noos-cortexm

