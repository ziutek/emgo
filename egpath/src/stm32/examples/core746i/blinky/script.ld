ISRStack = 64 * 32;
MainStack = 64 * 32;
TaskStack = 0;
MaxTasks = 1;

INCLUDE stm32/f7x6xg
INCLUDE stm32/loaditcram
INCLUDE noos-cortexm

/*
Core746I BOOT CONFIG switch determines which BOOT_ADD register is used to obtain
the boot addres:

- FLASH:  BOOT_ADD0 (0x1FFF 0010),
- SYSTEM: BOOT_ADD1 (0x1FFF 0018).

Boot addres is equal to BOOT_ADDx&0xFFFF << 14. Example values:

- 0x0000: ITCM RAM      (0x0000 0000), required by stm32/loaditcram,
- 0x0040: Bootloader    (0x0010 0000)
- 0x0080: Flash on ITCM (0x0020 0000), required by stm32/loadflash,
- 0x2000: Flash on AXIM (0x0800 0000)
- 0x8000: DTCM RAM      (0x2000 0000), required by stm32/loadram,
- 0x8004: SRAM          (0x2001 0000).

Boot addresses can be changed by modify FLASH_OPTCR1 control register. Use oppenocd command stm32f2x options_write 0 ... or gdb (see ../change_optcr_gdb.txt).

BE CAREFUL! Manipulating FLASH_OPTCR in a wrong way can PERMANENTLY BRICK your device.
*/
