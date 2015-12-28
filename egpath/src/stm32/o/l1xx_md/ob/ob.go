// +build l1xx_md

// Peripheral: OB_Periph  Option Bytes Registers.
// Instances:
//  OB  mmap.OB_BASE
// Registers:
//  0x00 32  RDP     Read protection register.
//  0x04 32  USER    user register.
//  0x08 32  WRP01   write protection register 0 1.
//  0x0C 32  WRP23   write protection register 2 3.
//  0x10 32  WRP45   write protection register 4 5.
//  0x14 32  WRP67   write protection register 6 7.
//  0x18 32  WRP89   write protection register 8 9.
//  0x1C 32  WRP1011 write protection register 10 11.
//  0x80 32  WRP1213 write protection register 12 13.
//  0x84 32  WRP1415 write protection register 14 15.
// Import:
//  stm32/o/l1xx_md/mmap
package ob
