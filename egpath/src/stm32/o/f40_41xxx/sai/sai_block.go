// +build f40_41xxx

// Peripheral: SAI_Block_Periph  Serial Audio Interface.
// Instances:
//  SAI1_Block_A  mmap.SAI1_Block_A_BASE
//  SAI1_Block_B  mmap.SAI1_Block_B_BASE
// Registers:
//  0x00 32  CR1   SAI block x configuration register 1.
//  0x04 32  CR2   SAI block x configuration register 2.
//  0x08 32  FRCR  SAI block x frame configuration register.
//  0x0C 32  SLOTR SAI block x slot register.
//  0x10 32  IMR   SAI block x interrupt mask register.
//  0x14 32  SR    SAI block x status register.
//  0x18 32  CLRFR SAI block x clear flag register.
//  0x1C 32  DR    SAI block x data register.
// Import:
//  stm32/o/f40_41xxx/mmap
package sai
