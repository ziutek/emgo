// Peripheral: FLASH_Periph  FLASH Registers.
// Instances:
//  FLASH  mmap.FLASH_R_BASE
// Registers:
//  0x00 32  ACR     Access control register.
//  0x04 32  PECR    Program/erase control register.
//  0x08 32  PDKEYR  Power down key register.
//  0x0C 32  PEKEYR  Program/erase key register.
//  0x10 32  PRGKEYR Program memory key register.
//  0x14 32  OPTKEYR Option byte key register.
//  0x18 32  SR      Status register.
//  0x1C 32  OBR     Option byte register.
//  0x20 32  WRPR    Write protection register.
//  0x80 32  WRPR1   Write protection register 1.
//  0x84 32  WRPR2   Write protection register 2.
//  0x88 32  WRPR3   Write protection register 3.
// Import:
//  stm32/o/l1xx_md/mmap
package flash

const (
	LATENCY  ACR_Bits = 0x01 << 0 //+ Latency.
	PRFTEN   ACR_Bits = 0x01 << 1 //+ Prefetch Buffer Enable.
	ACC64    ACR_Bits = 0x01 << 2 //+ Access 64 bits.
	SLEEP_PD ACR_Bits = 0x01 << 3 //+ Flash mode during sleep mode.
	RUN_PD   ACR_Bits = 0x01 << 4 //+ Flash mode during RUN mode.
)

const (
	PELOCK     PECR_Bits = 0x01 << 0  //+ FLASH_PECR and Flash data Lock.
	PRGLOCK    PECR_Bits = 0x01 << 1  //+ Program matrix Lock.
	OPTLOCK    PECR_Bits = 0x01 << 2  //+ Option byte matrix Lock.
	PROG       PECR_Bits = 0x01 << 3  //+ Program matrix selection.
	DATA       PECR_Bits = 0x01 << 4  //+ Data matrix selection.
	FTDW       PECR_Bits = 0x01 << 8  //+ Fixed Time Data write for Word/Half Word/Byte programming.
	ERASE      PECR_Bits = 0x01 << 9  //+ Page erasing mode.
	FPRG       PECR_Bits = 0x01 << 10 //+ Fast Page/Half Page programming mode.
	PARALLBANK PECR_Bits = 0x01 << 15 //+ Parallel Bank mode.
	EOPIE      PECR_Bits = 0x01 << 16 //+ End of programming interrupt.
	ERRIE      PECR_Bits = 0x01 << 17 //+ Error interrupt.
	OBL_LAUNCH PECR_Bits = 0x01 << 18 //+ Launch the option byte loading.
)

const ()

const ()

const ()

const ()

const (
	BSY        SR_Bits = 0x01 << 0  //+ Busy.
	EOP        SR_Bits = 0x01 << 1  //+ End Of Programming.
	ENHV       SR_Bits = 0x01 << 2  //+ End of high voltage.
	READY      SR_Bits = 0x01 << 3  //+ Flash ready after low power mode.
	WRPERR     SR_Bits = 0x01 << 8  //+ Write protected error.
	PGAERR     SR_Bits = 0x01 << 9  //+ Programming Alignment Error.
	SIZERR     SR_Bits = 0x01 << 10 //+ Size error.
	OPTVERR    SR_Bits = 0x01 << 11 //+ Option validity error.
	OPTVERRUSR SR_Bits = 0x01 << 12 //+ Option User validity error.
	RDERR      SR_Bits = 0x01 << 13 //+ Read protected error.
)

const (
	RDPRT      OBR_Bits = 0x55 << 1  //+ Read Protection.
	SPRMOD     OBR_Bits = 0x01 << 8  //+ Selection of protection mode of WPRi bits.
	BOR_LEV    OBR_Bits = 0x0F << 16 //+ BOR_LEV[3:0] Brown Out Reset Threshold Level.
	IWDG_SW    OBR_Bits = 0x01 << 20 //+ IWDG_SW.
	nRST_STOP  OBR_Bits = 0x01 << 21 //+ nRST_STOP.
	nRST_STDBY OBR_Bits = 0x01 << 22 //+ nRST_STDBY.
	BFB2       OBR_Bits = 0x01 << 23 //+ BFB2(available only in STM32L1xx High-density devices).
)

const (
	WRP WRPR_Bits = 0xFFFFFFFF << 0 //+ Write Protection bits.
)

const (
	WRP WRPR1_Bits = 0xFFFFFFFF << 0 //+ Write Protection bits (available only in STM32L1xx.
)

const (
	WRP WRPR2_Bits = 0xFFFFFFFF << 0 //+ Write Protection bits (available only in STM32L1xx.
)
