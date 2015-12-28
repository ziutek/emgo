// +build f40_41xxx

// Peripheral: FSMC_Bank2_Periph  Flexible Static Memory Controller Bank2.
// Instances:
//  FSMC_Bank2  mmap.FSMC_Bank2_R_BASE
// Registers:
//  0x00 32  PCR2  NAND Flash control register 2.
//  0x04 32  SR2   NAND Flash FIFO status and interrupt register 2.
//  0x08 32  PMEM2 NAND Flash Common memory space timing register 2.
//  0x0C 32  PATT2 NAND Flash Attribute memory space timing register 2.
//  0x14 32  ECCR2 NAND Flash ECC result registers 2.
// Import:
//  stm32/o/f40_41xxx/mmap
package fsmc
