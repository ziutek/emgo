// Peripheral: DBGMCU_Periph  Debug MCU.
// Instances:
//  DBGMCU  mmap.DBGMCU_BASE
// Registers:
//  0x00 32  IDCODE MCU device ID code.
//  0x04 32  CR     Debug MCU configuration register.
//  0x08 32  APB1FZ Debug MCU APB1 freeze register.
//  0x0C 32  APB2FZ Debug MCU APB2 freeze register.
// Import:
//  stm32/o/f40_41xxx/mmap
package dbgmcu

const (
	DEV_ID IDCODE_Bits = 0xFFF << 0   //+
	REV_ID IDCODE_Bits = 0xFFFF << 16 //+
)

const (
	DBG_SLEEP    CR_Bits = 0x01 << 0 //+
	DBG_STOP     CR_Bits = 0x01 << 1 //+
	DBG_STANDBY  CR_Bits = 0x01 << 2 //+
	TRACE_IOEN   CR_Bits = 0x01 << 5 //+
	TRACE_MODE   CR_Bits = 0x03 << 6 //+
	TRACE_MODE_0 CR_Bits = 0x01 << 6 //  Bit 0.
	TRACE_MODE_1 CR_Bits = 0x02 << 6 //  Bit 1.
)

const (
	DBG_TIM2_STOP          APB1FZ_Bits = 0x01 << 0  //+
	DBG_TIM3_STOP          APB1FZ_Bits = 0x01 << 1  //+
	DBG_TIM4_STOP          APB1FZ_Bits = 0x01 << 2  //+
	DBG_TIM5_STOP          APB1FZ_Bits = 0x01 << 3  //+
	DBG_TIM6_STOP          APB1FZ_Bits = 0x01 << 4  //+
	DBG_TIM7_STOP          APB1FZ_Bits = 0x01 << 5  //+
	DBG_TIM12_STOP         APB1FZ_Bits = 0x01 << 6  //+
	DBG_TIM13_STOP         APB1FZ_Bits = 0x01 << 7  //+
	DBG_TIM14_STOP         APB1FZ_Bits = 0x01 << 8  //+
	DBG_RTC_STOP           APB1FZ_Bits = 0x01 << 10 //+
	DBG_WWDG_STOP          APB1FZ_Bits = 0x01 << 11 //+
	DBG_IWDG_STOP          APB1FZ_Bits = 0x01 << 12 //+
	DBG_I2C1_SMBUS_TIMEOUT APB1FZ_Bits = 0x01 << 21 //+
	DBG_I2C2_SMBUS_TIMEOUT APB1FZ_Bits = 0x01 << 22 //+
	DBG_I2C3_SMBUS_TIMEOUT APB1FZ_Bits = 0x01 << 23 //+
	DBG_CAN1_STOP          APB1FZ_Bits = 0x01 << 25 //+
	DBG_CAN2_STOP          APB1FZ_Bits = 0x01 << 26 //+
	DBG_TIM1_STOP          APB1FZ_Bits = 0x01 << 0
	DBG_TIM8_STOP          APB1FZ_Bits = 0x01 << 1
	DBG_TIM9_STOP          APB1FZ_Bits = 0x01 << 16 //+
	DBG_TIM10_STOP         APB1FZ_Bits = 0x01 << 17 //+
	DBG_TIM11_STOP         APB1FZ_Bits = 0x01 << 18 //+
)
