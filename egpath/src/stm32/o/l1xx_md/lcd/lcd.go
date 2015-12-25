// Peripheral: LCD_Periph  LCD.
// Instances:
//  LCD  mmap.LCD_BASE
// Registers:
//  0x00 32  CR      Control register.
//  0x04 32  FCR     Frame control register.
//  0x08 32  SR      Status register.
//  0x0C 32  CLR     Clear register.
//  0x14 32  RAM[16] Display memory.
// Import:
//  stm32/o/l1xx_md/mmap
package lcd

const (
	LCDEN   CR_Bits = 0x01 << 0 //+ LCD Enable Bit.
	VSEL    CR_Bits = 0x01 << 1 //+ Voltage source selector Bit.
	DUTY    CR_Bits = 0x07 << 2 //+ DUTY[2:0] bits (Duty selector).
	DUTY_0  CR_Bits = 0x01 << 2 //  Duty selector Bit 0.
	DUTY_1  CR_Bits = 0x02 << 2 //  Duty selector Bit 1.
	DUTY_2  CR_Bits = 0x04 << 2 //  Duty selector Bit 2.
	BIAS    CR_Bits = 0x03 << 5 //+ BIAS[1:0] bits (Bias selector).
	BIAS_0  CR_Bits = 0x01 << 5 //  Bias selector Bit 0.
	BIAS_1  CR_Bits = 0x02 << 5 //  Bias selector Bit 1.
	MUX_SEG CR_Bits = 0x01 << 7 //+ Mux Segment Enable Bit.
)

const (
	HD       FCR_Bits = 0x01 << 0  //+ High Drive Enable Bit.
	SOFIE    FCR_Bits = 0x01 << 1  //+ Start of Frame Interrupt Enable Bit.
	UDDIE    FCR_Bits = 0x01 << 3  //+ Update Display Done Interrupt Enable Bit.
	PON      FCR_Bits = 0x07 << 4  //+ PON[2:0] bits (Puls ON Duration).
	PON_0    FCR_Bits = 0x01 << 4  //  Bit 0.
	PON_1    FCR_Bits = 0x02 << 4  //  Bit 1.
	PON_2    FCR_Bits = 0x04 << 4  //  Bit 2.
	DEAD     FCR_Bits = 0x07 << 7  //+ DEAD[2:0] bits (DEAD Time).
	DEAD_0   FCR_Bits = 0x01 << 7  //  Bit 0.
	DEAD_1   FCR_Bits = 0x02 << 7  //  Bit 1.
	DEAD_2   FCR_Bits = 0x04 << 7  //  Bit 2.
	CC       FCR_Bits = 0x07 << 10 //+ CC[2:0] bits (Contrast Control).
	CC_0     FCR_Bits = 0x01 << 10 //  Bit 0.
	CC_1     FCR_Bits = 0x02 << 10 //  Bit 1.
	CC_2     FCR_Bits = 0x04 << 10 //  Bit 2.
	BLINKF   FCR_Bits = 0x07 << 13 //+ BLINKF[2:0] bits (Blink Frequency).
	BLINKF_0 FCR_Bits = 0x01 << 13 //  Bit 0.
	BLINKF_1 FCR_Bits = 0x02 << 13 //  Bit 1.
	BLINKF_2 FCR_Bits = 0x04 << 13 //  Bit 2.
	BLINK    FCR_Bits = 0x03 << 16 //+ BLINK[1:0] bits (Blink Enable).
	BLINK_0  FCR_Bits = 0x01 << 16 //  Bit 0.
	BLINK_1  FCR_Bits = 0x02 << 16 //  Bit 1.
	DIV      FCR_Bits = 0x0F << 18 //+ DIV[3:0] bits (Divider).
	PS       FCR_Bits = 0x0F << 22 //+ PS[3:0] bits (Prescaler).
)

const (
	ENS   SR_Bits = 0x01 << 0 //+ LCD Enabled Bit.
	SOF   SR_Bits = 0x01 << 1 //+ Start Of Frame Flag Bit.
	UDR   SR_Bits = 0x01 << 2 //+ Update Display Request Bit.
	UDD   SR_Bits = 0x01 << 3 //+ Update Display Done Flag Bit.
	RDY   SR_Bits = 0x01 << 4 //+ Ready Flag Bit.
	FCRSR SR_Bits = 0x01 << 5 //+ LCD FCR Register Synchronization Flag Bit.
)

const (
	SOFC CLR_Bits = 0x01 << 1 //+ Start Of Frame Flag Clear Bit.
	UDDC CLR_Bits = 0x01 << 3 //+ Update Display Done Flag Clear Bit.
)
