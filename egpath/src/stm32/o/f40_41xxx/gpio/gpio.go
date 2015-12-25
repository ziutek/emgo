// Peripheral: GPIO_Periph  General Purpose I/O.
// Instances:
//  GPIOA  mmap.GPIOA_BASE
//  GPIOB  mmap.GPIOB_BASE
//  GPIOC  mmap.GPIOC_BASE
//  GPIOD  mmap.GPIOD_BASE
//  GPIOE  mmap.GPIOE_BASE
//  GPIOF  mmap.GPIOF_BASE
//  GPIOG  mmap.GPIOG_BASE
//  GPIOH  mmap.GPIOH_BASE
//  GPIOI  mmap.GPIOI_BASE
//  GPIOJ  mmap.GPIOJ_BASE
//  GPIOK  mmap.GPIOK_BASE
// Registers:
//  0x00 32  MODER   Port mode register.
//  0x04 32  OTYPER  Port output type register.
//  0x08 32  OSPEEDR Port output speed register.
//  0x0C 32  PUPDR   Port pull-up/pull-down register.
//  0x10 32  IDR     Port input data register.
//  0x14 32  ODR     Port output data register.
//  0x18 16  BSRRL   Port bit set/reset low register.
//  0x1A 16  BSRRH   Port bit set/reset high register.
//  0x1C 32  LCKR    Port configuration lock register.
//  0x20 32  AFR[2]  Alternate function registers.
// Import:
//  stm32/o/f40_41xxx/mmap
package gpio

const (
	MODER0    MODER_Bits = 0x03 << 0 //+
	MODER0_0  MODER_Bits = 0x01 << 0
	MODER0_1  MODER_Bits = 0x02 << 0
	MODER1    MODER_Bits = 0x03 << 2 //+
	MODER1_0  MODER_Bits = 0x01 << 2
	MODER1_1  MODER_Bits = 0x02 << 2
	MODER2    MODER_Bits = 0x03 << 4 //+
	MODER2_0  MODER_Bits = 0x01 << 4
	MODER2_1  MODER_Bits = 0x02 << 4
	MODER3    MODER_Bits = 0x03 << 6 //+
	MODER3_0  MODER_Bits = 0x01 << 6
	MODER3_1  MODER_Bits = 0x02 << 6
	MODER4    MODER_Bits = 0x03 << 8 //+
	MODER4_0  MODER_Bits = 0x01 << 8
	MODER4_1  MODER_Bits = 0x02 << 8
	MODER5    MODER_Bits = 0x03 << 10 //+
	MODER5_0  MODER_Bits = 0x01 << 10
	MODER5_1  MODER_Bits = 0x02 << 10
	MODER6    MODER_Bits = 0x03 << 12 //+
	MODER6_0  MODER_Bits = 0x01 << 12
	MODER6_1  MODER_Bits = 0x02 << 12
	MODER7    MODER_Bits = 0x03 << 14 //+
	MODER7_0  MODER_Bits = 0x01 << 14
	MODER7_1  MODER_Bits = 0x02 << 14
	MODER8    MODER_Bits = 0x03 << 16 //+
	MODER8_0  MODER_Bits = 0x01 << 16
	MODER8_1  MODER_Bits = 0x02 << 16
	MODER9    MODER_Bits = 0x03 << 18 //+
	MODER9_0  MODER_Bits = 0x01 << 18
	MODER9_1  MODER_Bits = 0x02 << 18
	MODER10   MODER_Bits = 0x03 << 20 //+
	MODER10_0 MODER_Bits = 0x01 << 20
	MODER10_1 MODER_Bits = 0x02 << 20
	MODER11   MODER_Bits = 0x03 << 22 //+
	MODER11_0 MODER_Bits = 0x01 << 22
	MODER11_1 MODER_Bits = 0x02 << 22
	MODER12   MODER_Bits = 0x03 << 24 //+
	MODER12_0 MODER_Bits = 0x01 << 24
	MODER12_1 MODER_Bits = 0x02 << 24
	MODER13   MODER_Bits = 0x03 << 26 //+
	MODER13_0 MODER_Bits = 0x01 << 26
	MODER13_1 MODER_Bits = 0x02 << 26
	MODER14   MODER_Bits = 0x03 << 28 //+
	MODER14_0 MODER_Bits = 0x01 << 28
	MODER14_1 MODER_Bits = 0x02 << 28
	MODER15   MODER_Bits = 0x03 << 30 //+
	MODER15_0 MODER_Bits = 0x01 << 30
	MODER15_1 MODER_Bits = 0x02 << 30
)

const (
	OT_0  OTYPER_Bits = 0x01 << 0  //+
	OT_1  OTYPER_Bits = 0x01 << 1  //+
	OT_2  OTYPER_Bits = 0x01 << 2  //+
	OT_3  OTYPER_Bits = 0x01 << 3  //+
	OT_4  OTYPER_Bits = 0x01 << 4  //+
	OT_5  OTYPER_Bits = 0x01 << 5  //+
	OT_6  OTYPER_Bits = 0x01 << 6  //+
	OT_7  OTYPER_Bits = 0x01 << 7  //+
	OT_8  OTYPER_Bits = 0x01 << 8  //+
	OT_9  OTYPER_Bits = 0x01 << 9  //+
	OT_10 OTYPER_Bits = 0x01 << 10 //+
	OT_11 OTYPER_Bits = 0x01 << 11 //+
	OT_12 OTYPER_Bits = 0x01 << 12 //+
	OT_13 OTYPER_Bits = 0x01 << 13 //+
	OT_14 OTYPER_Bits = 0x01 << 14 //+
	OT_15 OTYPER_Bits = 0x01 << 15 //+
)

const (
	OSPEEDR0    OSPEEDR_Bits = 0x03 << 0 //+
	OSPEEDR0_0  OSPEEDR_Bits = 0x01 << 0
	OSPEEDR0_1  OSPEEDR_Bits = 0x02 << 0
	OSPEEDR1    OSPEEDR_Bits = 0x03 << 2 //+
	OSPEEDR1_0  OSPEEDR_Bits = 0x01 << 2
	OSPEEDR1_1  OSPEEDR_Bits = 0x02 << 2
	OSPEEDR2    OSPEEDR_Bits = 0x03 << 4 //+
	OSPEEDR2_0  OSPEEDR_Bits = 0x01 << 4
	OSPEEDR2_1  OSPEEDR_Bits = 0x02 << 4
	OSPEEDR3    OSPEEDR_Bits = 0x03 << 6 //+
	OSPEEDR3_0  OSPEEDR_Bits = 0x01 << 6
	OSPEEDR3_1  OSPEEDR_Bits = 0x02 << 6
	OSPEEDR4    OSPEEDR_Bits = 0x03 << 8 //+
	OSPEEDR4_0  OSPEEDR_Bits = 0x01 << 8
	OSPEEDR4_1  OSPEEDR_Bits = 0x02 << 8
	OSPEEDR5    OSPEEDR_Bits = 0x03 << 10 //+
	OSPEEDR5_0  OSPEEDR_Bits = 0x01 << 10
	OSPEEDR5_1  OSPEEDR_Bits = 0x02 << 10
	OSPEEDR6    OSPEEDR_Bits = 0x03 << 12 //+
	OSPEEDR6_0  OSPEEDR_Bits = 0x01 << 12
	OSPEEDR6_1  OSPEEDR_Bits = 0x02 << 12
	OSPEEDR7    OSPEEDR_Bits = 0x03 << 14 //+
	OSPEEDR7_0  OSPEEDR_Bits = 0x01 << 14
	OSPEEDR7_1  OSPEEDR_Bits = 0x02 << 14
	OSPEEDR8    OSPEEDR_Bits = 0x03 << 16 //+
	OSPEEDR8_0  OSPEEDR_Bits = 0x01 << 16
	OSPEEDR8_1  OSPEEDR_Bits = 0x02 << 16
	OSPEEDR9    OSPEEDR_Bits = 0x03 << 18 //+
	OSPEEDR9_0  OSPEEDR_Bits = 0x01 << 18
	OSPEEDR9_1  OSPEEDR_Bits = 0x02 << 18
	OSPEEDR10   OSPEEDR_Bits = 0x03 << 20 //+
	OSPEEDR10_0 OSPEEDR_Bits = 0x01 << 20
	OSPEEDR10_1 OSPEEDR_Bits = 0x02 << 20
	OSPEEDR11   OSPEEDR_Bits = 0x03 << 22 //+
	OSPEEDR11_0 OSPEEDR_Bits = 0x01 << 22
	OSPEEDR11_1 OSPEEDR_Bits = 0x02 << 22
	OSPEEDR12   OSPEEDR_Bits = 0x03 << 24 //+
	OSPEEDR12_0 OSPEEDR_Bits = 0x01 << 24
	OSPEEDR12_1 OSPEEDR_Bits = 0x02 << 24
	OSPEEDR13   OSPEEDR_Bits = 0x03 << 26 //+
	OSPEEDR13_0 OSPEEDR_Bits = 0x01 << 26
	OSPEEDR13_1 OSPEEDR_Bits = 0x02 << 26
	OSPEEDR14   OSPEEDR_Bits = 0x03 << 28 //+
	OSPEEDR14_0 OSPEEDR_Bits = 0x01 << 28
	OSPEEDR14_1 OSPEEDR_Bits = 0x02 << 28
	OSPEEDR15   OSPEEDR_Bits = 0x03 << 30 //+
	OSPEEDR15_0 OSPEEDR_Bits = 0x01 << 30
	OSPEEDR15_1 OSPEEDR_Bits = 0x02 << 30
)

const (
	PUPDR0    PUPDR_Bits = 0x03 << 0 //+
	PUPDR0_0  PUPDR_Bits = 0x01 << 0
	PUPDR0_1  PUPDR_Bits = 0x02 << 0
	PUPDR1    PUPDR_Bits = 0x03 << 2 //+
	PUPDR1_0  PUPDR_Bits = 0x01 << 2
	PUPDR1_1  PUPDR_Bits = 0x02 << 2
	PUPDR2    PUPDR_Bits = 0x03 << 4 //+
	PUPDR2_0  PUPDR_Bits = 0x01 << 4
	PUPDR2_1  PUPDR_Bits = 0x02 << 4
	PUPDR3    PUPDR_Bits = 0x03 << 6 //+
	PUPDR3_0  PUPDR_Bits = 0x01 << 6
	PUPDR3_1  PUPDR_Bits = 0x02 << 6
	PUPDR4    PUPDR_Bits = 0x03 << 8 //+
	PUPDR4_0  PUPDR_Bits = 0x01 << 8
	PUPDR4_1  PUPDR_Bits = 0x02 << 8
	PUPDR5    PUPDR_Bits = 0x03 << 10 //+
	PUPDR5_0  PUPDR_Bits = 0x01 << 10
	PUPDR5_1  PUPDR_Bits = 0x02 << 10
	PUPDR6    PUPDR_Bits = 0x03 << 12 //+
	PUPDR6_0  PUPDR_Bits = 0x01 << 12
	PUPDR6_1  PUPDR_Bits = 0x02 << 12
	PUPDR7    PUPDR_Bits = 0x03 << 14 //+
	PUPDR7_0  PUPDR_Bits = 0x01 << 14
	PUPDR7_1  PUPDR_Bits = 0x02 << 14
	PUPDR8    PUPDR_Bits = 0x03 << 16 //+
	PUPDR8_0  PUPDR_Bits = 0x01 << 16
	PUPDR8_1  PUPDR_Bits = 0x02 << 16
	PUPDR9    PUPDR_Bits = 0x03 << 18 //+
	PUPDR9_0  PUPDR_Bits = 0x01 << 18
	PUPDR9_1  PUPDR_Bits = 0x02 << 18
	PUPDR10   PUPDR_Bits = 0x03 << 20 //+
	PUPDR10_0 PUPDR_Bits = 0x01 << 20
	PUPDR10_1 PUPDR_Bits = 0x02 << 20
	PUPDR11   PUPDR_Bits = 0x03 << 22 //+
	PUPDR11_0 PUPDR_Bits = 0x01 << 22
	PUPDR11_1 PUPDR_Bits = 0x02 << 22
	PUPDR12   PUPDR_Bits = 0x03 << 24 //+
	PUPDR12_0 PUPDR_Bits = 0x01 << 24
	PUPDR12_1 PUPDR_Bits = 0x02 << 24
	PUPDR13   PUPDR_Bits = 0x03 << 26 //+
	PUPDR13_0 PUPDR_Bits = 0x01 << 26
	PUPDR13_1 PUPDR_Bits = 0x02 << 26
	PUPDR14   PUPDR_Bits = 0x03 << 28 //+
	PUPDR14_0 PUPDR_Bits = 0x01 << 28
	PUPDR14_1 PUPDR_Bits = 0x02 << 28
	PUPDR15   PUPDR_Bits = 0x03 << 30 //+
	PUPDR15_0 PUPDR_Bits = 0x01 << 30
	PUPDR15_1 PUPDR_Bits = 0x02 << 30
)

const (
	IDR_0  IDR_Bits = 0x01 << 0  //+
	IDR_1  IDR_Bits = 0x01 << 1  //+
	IDR_2  IDR_Bits = 0x01 << 2  //+
	IDR_3  IDR_Bits = 0x01 << 3  //+
	IDR_4  IDR_Bits = 0x01 << 4  //+
	IDR_5  IDR_Bits = 0x01 << 5  //+
	IDR_6  IDR_Bits = 0x01 << 6  //+
	IDR_7  IDR_Bits = 0x01 << 7  //+
	IDR_8  IDR_Bits = 0x01 << 8  //+
	IDR_9  IDR_Bits = 0x01 << 9  //+
	IDR_10 IDR_Bits = 0x01 << 10 //+
	IDR_11 IDR_Bits = 0x01 << 11 //+
	IDR_12 IDR_Bits = 0x01 << 12 //+
	IDR_13 IDR_Bits = 0x01 << 13 //+
	IDR_14 IDR_Bits = 0x01 << 14 //+
	IDR_15 IDR_Bits = 0x01 << 15 //+
)

const (
	ODR_0  ODR_Bits = 0x01 << 0  //+
	ODR_1  ODR_Bits = 0x01 << 1  //+
	ODR_2  ODR_Bits = 0x01 << 2  //+
	ODR_3  ODR_Bits = 0x01 << 3  //+
	ODR_4  ODR_Bits = 0x01 << 4  //+
	ODR_5  ODR_Bits = 0x01 << 5  //+
	ODR_6  ODR_Bits = 0x01 << 6  //+
	ODR_7  ODR_Bits = 0x01 << 7  //+
	ODR_8  ODR_Bits = 0x01 << 8  //+
	ODR_9  ODR_Bits = 0x01 << 9  //+
	ODR_10 ODR_Bits = 0x01 << 10 //+
	ODR_11 ODR_Bits = 0x01 << 11 //+
	ODR_12 ODR_Bits = 0x01 << 12 //+
	ODR_13 ODR_Bits = 0x01 << 13 //+
	ODR_14 ODR_Bits = 0x01 << 14 //+
	ODR_15 ODR_Bits = 0x01 << 15 //+
)
