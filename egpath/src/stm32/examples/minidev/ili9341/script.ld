ISRStack = 64*32;
MainStack = 192*32;
TaskStack = 0;
MaxTasks = 1;

/* bootSRAM code should by placed at 0x200001E0. */
/*bootOffset = 0x1E0;*/

/* F103x8 usually contains 128 KB of Flash like F103xB. Use this "property". */
INCLUDE stm32/f103xb
INCLUDE stm32/loadflash
INCLUDE noos-cortexm

/* Include bootRAM and set it as entry point. */
/*ENTRY(bootRAM)*/
