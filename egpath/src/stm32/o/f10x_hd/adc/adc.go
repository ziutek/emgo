// +build f10x_hd

// Peripheral: ADC_Periph  Analog to Digital Converter.
// Instances:
//  ADC1  mmap.ADC1_BASE
//  ADC2  mmap.ADC2_BASE
//  ADC3  mmap.ADC3_BASE
// Registers:
//  0x00 32  SR
//  0x04 32  CR1
//  0x08 32  CR2
//  0x0C 32  SMPR1
//  0x10 32  SMPR2
//  0x14 32  JOFR1
//  0x18 32  JOFR2
//  0x1C 32  JOFR3
//  0x20 32  JOFR4
//  0x24 32  HTR
//  0x28 32  LTR
//  0x2C 32  SQR1
//  0x30 32  SQR2
//  0x34 32  SQR3
//  0x38 32  JSQR
//  0x3C 32  JDR1
//  0x40 32  JDR2
//  0x44 32  JDR3
//  0x48 32  JDR4
//  0x4C 32  DR
// Import:
//  stm32/o/f10x_hd/mmap
package adc

const (
	AWD   SR_Bits = 0x01 << 0 //+ Analog watchdog flag.
	EOC   SR_Bits = 0x01 << 1 //+ End of conversion.
	JEOC  SR_Bits = 0x01 << 2 //+ Injected channel end of conversion.
	JSTRT SR_Bits = 0x01 << 3 //+ Injected channel Start flag.
	STRT  SR_Bits = 0x01 << 4 //+ Regular channel Start flag.
)

const (
	AWDCH     CR1_Bits = 0x1F << 0  //+ AWDCH[4:0] bits (Analog watchdog channel select bits).
	AWDCH_0   CR1_Bits = 0x01 << 0  //  Bit 0.
	AWDCH_1   CR1_Bits = 0x02 << 0  //  Bit 1.
	AWDCH_2   CR1_Bits = 0x04 << 0  //  Bit 2.
	AWDCH_3   CR1_Bits = 0x08 << 0  //  Bit 3.
	AWDCH_4   CR1_Bits = 0x10 << 0  //  Bit 4.
	EOCIE     CR1_Bits = 0x01 << 5  //+ Interrupt enable for EOC.
	AWDIE     CR1_Bits = 0x01 << 6  //+ Analog Watchdog interrupt enable.
	JEOCIE    CR1_Bits = 0x01 << 7  //+ Interrupt enable for injected channels.
	SCAN      CR1_Bits = 0x01 << 8  //+ Scan mode.
	AWDSGL    CR1_Bits = 0x01 << 9  //+ Enable the watchdog on a single channel in scan mode.
	JAUTO     CR1_Bits = 0x01 << 10 //+ Automatic injected group conversion.
	DISCEN    CR1_Bits = 0x01 << 11 //+ Discontinuous mode on regular channels.
	JDISCEN   CR1_Bits = 0x01 << 12 //+ Discontinuous mode on injected channels.
	DISCNUM   CR1_Bits = 0x07 << 13 //+ DISCNUM[2:0] bits (Discontinuous mode channel count).
	DISCNUM_0 CR1_Bits = 0x01 << 13 //  Bit 0.
	DISCNUM_1 CR1_Bits = 0x02 << 13 //  Bit 1.
	DISCNUM_2 CR1_Bits = 0x04 << 13 //  Bit 2.
	DUALMOD   CR1_Bits = 0x0F << 16 //+ DUALMOD[3:0] bits (Dual mode selection).
	DUALMOD_0 CR1_Bits = 0x01 << 16 //  Bit 0.
	DUALMOD_1 CR1_Bits = 0x02 << 16 //  Bit 1.
	DUALMOD_2 CR1_Bits = 0x04 << 16 //  Bit 2.
	DUALMOD_3 CR1_Bits = 0x08 << 16 //  Bit 3.
	JAWDEN    CR1_Bits = 0x01 << 22 //+ Analog watchdog enable on injected channels.
	AWDEN     CR1_Bits = 0x01 << 23 //+ Analog watchdog enable on regular channels.
)

const (
	ADON      CR2_Bits = 0x01 << 0  //+ A/D Converter ON / OFF.
	CONT      CR2_Bits = 0x01 << 1  //+ Continuous Conversion.
	CAL       CR2_Bits = 0x01 << 2  //+ A/D Calibration.
	RSTCAL    CR2_Bits = 0x01 << 3  //+ Reset Calibration.
	DMA       CR2_Bits = 0x01 << 8  //+ Direct Memory access mode.
	ALIGN     CR2_Bits = 0x01 << 11 //+ Data Alignment.
	JEXTSEL   CR2_Bits = 0x07 << 12 //+ JEXTSEL[2:0] bits (External event select for injected group).
	JEXTSEL_0 CR2_Bits = 0x01 << 12 //  Bit 0.
	JEXTSEL_1 CR2_Bits = 0x02 << 12 //  Bit 1.
	JEXTSEL_2 CR2_Bits = 0x04 << 12 //  Bit 2.
	JEXTTRIG  CR2_Bits = 0x01 << 15 //+ External Trigger Conversion mode for injected channels.
	EXTSEL    CR2_Bits = 0x07 << 17 //+ EXTSEL[2:0] bits (External Event Select for regular group).
	EXTSEL_0  CR2_Bits = 0x01 << 17 //  Bit 0.
	EXTSEL_1  CR2_Bits = 0x02 << 17 //  Bit 1.
	EXTSEL_2  CR2_Bits = 0x04 << 17 //  Bit 2.
	EXTTRIG   CR2_Bits = 0x01 << 20 //+ External Trigger Conversion mode for regular channels.
	JSWSTART  CR2_Bits = 0x01 << 21 //+ Start Conversion of injected channels.
	SWSTART   CR2_Bits = 0x01 << 22 //+ Start Conversion of regular channels.
	TSVREFE   CR2_Bits = 0x01 << 23 //+ Temperature Sensor and VREFINT Enable.
)

const (
	SMP10   SMPR1_Bits = 0x07 << 0  //+ SMP10[2:0] bits (Channel 10 Sample time selection).
	SMP10_0 SMPR1_Bits = 0x01 << 0  //  Bit 0.
	SMP10_1 SMPR1_Bits = 0x02 << 0  //  Bit 1.
	SMP10_2 SMPR1_Bits = 0x04 << 0  //  Bit 2.
	SMP11   SMPR1_Bits = 0x07 << 3  //+ SMP11[2:0] bits (Channel 11 Sample time selection).
	SMP11_0 SMPR1_Bits = 0x01 << 3  //  Bit 0.
	SMP11_1 SMPR1_Bits = 0x02 << 3  //  Bit 1.
	SMP11_2 SMPR1_Bits = 0x04 << 3  //  Bit 2.
	SMP12   SMPR1_Bits = 0x07 << 6  //+ SMP12[2:0] bits (Channel 12 Sample time selection).
	SMP12_0 SMPR1_Bits = 0x01 << 6  //  Bit 0.
	SMP12_1 SMPR1_Bits = 0x02 << 6  //  Bit 1.
	SMP12_2 SMPR1_Bits = 0x04 << 6  //  Bit 2.
	SMP13   SMPR1_Bits = 0x07 << 9  //+ SMP13[2:0] bits (Channel 13 Sample time selection).
	SMP13_0 SMPR1_Bits = 0x01 << 9  //  Bit 0.
	SMP13_1 SMPR1_Bits = 0x02 << 9  //  Bit 1.
	SMP13_2 SMPR1_Bits = 0x04 << 9  //  Bit 2.
	SMP14   SMPR1_Bits = 0x07 << 12 //+ SMP14[2:0] bits (Channel 14 Sample time selection).
	SMP14_0 SMPR1_Bits = 0x01 << 12 //  Bit 0.
	SMP14_1 SMPR1_Bits = 0x02 << 12 //  Bit 1.
	SMP14_2 SMPR1_Bits = 0x04 << 12 //  Bit 2.
	SMP15   SMPR1_Bits = 0x07 << 15 //+ SMP15[2:0] bits (Channel 15 Sample time selection).
	SMP15_0 SMPR1_Bits = 0x01 << 15 //  Bit 0.
	SMP15_1 SMPR1_Bits = 0x02 << 15 //  Bit 1.
	SMP15_2 SMPR1_Bits = 0x04 << 15 //  Bit 2.
	SMP16   SMPR1_Bits = 0x07 << 18 //+ SMP16[2:0] bits (Channel 16 Sample time selection).
	SMP16_0 SMPR1_Bits = 0x01 << 18 //  Bit 0.
	SMP16_1 SMPR1_Bits = 0x02 << 18 //  Bit 1.
	SMP16_2 SMPR1_Bits = 0x04 << 18 //  Bit 2.
	SMP17   SMPR1_Bits = 0x07 << 21 //+ SMP17[2:0] bits (Channel 17 Sample time selection).
	SMP17_0 SMPR1_Bits = 0x01 << 21 //  Bit 0.
	SMP17_1 SMPR1_Bits = 0x02 << 21 //  Bit 1.
	SMP17_2 SMPR1_Bits = 0x04 << 21 //  Bit 2.
)

const (
	SMP0   SMPR2_Bits = 0x07 << 0  //+ SMP0[2:0] bits (Channel 0 Sample time selection).
	SMP0_0 SMPR2_Bits = 0x01 << 0  //  Bit 0.
	SMP0_1 SMPR2_Bits = 0x02 << 0  //  Bit 1.
	SMP0_2 SMPR2_Bits = 0x04 << 0  //  Bit 2.
	SMP1   SMPR2_Bits = 0x07 << 3  //+ SMP1[2:0] bits (Channel 1 Sample time selection).
	SMP1_0 SMPR2_Bits = 0x01 << 3  //  Bit 0.
	SMP1_1 SMPR2_Bits = 0x02 << 3  //  Bit 1.
	SMP1_2 SMPR2_Bits = 0x04 << 3  //  Bit 2.
	SMP2   SMPR2_Bits = 0x07 << 6  //+ SMP2[2:0] bits (Channel 2 Sample time selection).
	SMP2_0 SMPR2_Bits = 0x01 << 6  //  Bit 0.
	SMP2_1 SMPR2_Bits = 0x02 << 6  //  Bit 1.
	SMP2_2 SMPR2_Bits = 0x04 << 6  //  Bit 2.
	SMP3   SMPR2_Bits = 0x07 << 9  //+ SMP3[2:0] bits (Channel 3 Sample time selection).
	SMP3_0 SMPR2_Bits = 0x01 << 9  //  Bit 0.
	SMP3_1 SMPR2_Bits = 0x02 << 9  //  Bit 1.
	SMP3_2 SMPR2_Bits = 0x04 << 9  //  Bit 2.
	SMP4   SMPR2_Bits = 0x07 << 12 //+ SMP4[2:0] bits (Channel 4 Sample time selection).
	SMP4_0 SMPR2_Bits = 0x01 << 12 //  Bit 0.
	SMP4_1 SMPR2_Bits = 0x02 << 12 //  Bit 1.
	SMP4_2 SMPR2_Bits = 0x04 << 12 //  Bit 2.
	SMP5   SMPR2_Bits = 0x07 << 15 //+ SMP5[2:0] bits (Channel 5 Sample time selection).
	SMP5_0 SMPR2_Bits = 0x01 << 15 //  Bit 0.
	SMP5_1 SMPR2_Bits = 0x02 << 15 //  Bit 1.
	SMP5_2 SMPR2_Bits = 0x04 << 15 //  Bit 2.
	SMP6   SMPR2_Bits = 0x07 << 18 //+ SMP6[2:0] bits (Channel 6 Sample time selection).
	SMP6_0 SMPR2_Bits = 0x01 << 18 //  Bit 0.
	SMP6_1 SMPR2_Bits = 0x02 << 18 //  Bit 1.
	SMP6_2 SMPR2_Bits = 0x04 << 18 //  Bit 2.
	SMP7   SMPR2_Bits = 0x07 << 21 //+ SMP7[2:0] bits (Channel 7 Sample time selection).
	SMP7_0 SMPR2_Bits = 0x01 << 21 //  Bit 0.
	SMP7_1 SMPR2_Bits = 0x02 << 21 //  Bit 1.
	SMP7_2 SMPR2_Bits = 0x04 << 21 //  Bit 2.
	SMP8   SMPR2_Bits = 0x07 << 24 //+ SMP8[2:0] bits (Channel 8 Sample time selection).
	SMP8_0 SMPR2_Bits = 0x01 << 24 //  Bit 0.
	SMP8_1 SMPR2_Bits = 0x02 << 24 //  Bit 1.
	SMP8_2 SMPR2_Bits = 0x04 << 24 //  Bit 2.
	SMP9   SMPR2_Bits = 0x07 << 27 //+ SMP9[2:0] bits (Channel 9 Sample time selection).
	SMP9_0 SMPR2_Bits = 0x01 << 27 //  Bit 0.
	SMP9_1 SMPR2_Bits = 0x02 << 27 //  Bit 1.
	SMP9_2 SMPR2_Bits = 0x04 << 27 //  Bit 2.
)

const (
	JOFFSET1 JOFR1_Bits = 0xFFF << 0 //+ Data offset for injected channel 1.
)

const (
	JOFFSET2 JOFR2_Bits = 0xFFF << 0 //+ Data offset for injected channel 2.
)

const (
	JOFFSET3 JOFR3_Bits = 0xFFF << 0 //+ Data offset for injected channel 3.
)

const (
	JOFFSET4 JOFR4_Bits = 0xFFF << 0 //+ Data offset for injected channel 4.
)

const (
	HT HTR_Bits = 0xFFF << 0 //+ Analog watchdog high threshold.
)

const (
	LT LTR_Bits = 0xFFF << 0 //+ Analog watchdog low threshold.
)

const (
	SQ13   SQR1_Bits = 0x1F << 0  //+ SQ13[4:0] bits (13th conversion in regular sequence).
	SQ13_0 SQR1_Bits = 0x01 << 0  //  Bit 0.
	SQ13_1 SQR1_Bits = 0x02 << 0  //  Bit 1.
	SQ13_2 SQR1_Bits = 0x04 << 0  //  Bit 2.
	SQ13_3 SQR1_Bits = 0x08 << 0  //  Bit 3.
	SQ13_4 SQR1_Bits = 0x10 << 0  //  Bit 4.
	SQ14   SQR1_Bits = 0x1F << 5  //+ SQ14[4:0] bits (14th conversion in regular sequence).
	SQ14_0 SQR1_Bits = 0x01 << 5  //  Bit 0.
	SQ14_1 SQR1_Bits = 0x02 << 5  //  Bit 1.
	SQ14_2 SQR1_Bits = 0x04 << 5  //  Bit 2.
	SQ14_3 SQR1_Bits = 0x08 << 5  //  Bit 3.
	SQ14_4 SQR1_Bits = 0x10 << 5  //  Bit 4.
	SQ15   SQR1_Bits = 0x1F << 10 //+ SQ15[4:0] bits (15th conversion in regular sequence).
	SQ15_0 SQR1_Bits = 0x01 << 10 //  Bit 0.
	SQ15_1 SQR1_Bits = 0x02 << 10 //  Bit 1.
	SQ15_2 SQR1_Bits = 0x04 << 10 //  Bit 2.
	SQ15_3 SQR1_Bits = 0x08 << 10 //  Bit 3.
	SQ15_4 SQR1_Bits = 0x10 << 10 //  Bit 4.
	SQ16   SQR1_Bits = 0x1F << 15 //+ SQ16[4:0] bits (16th conversion in regular sequence).
	SQ16_0 SQR1_Bits = 0x01 << 15 //  Bit 0.
	SQ16_1 SQR1_Bits = 0x02 << 15 //  Bit 1.
	SQ16_2 SQR1_Bits = 0x04 << 15 //  Bit 2.
	SQ16_3 SQR1_Bits = 0x08 << 15 //  Bit 3.
	SQ16_4 SQR1_Bits = 0x10 << 15 //  Bit 4.
	L      SQR1_Bits = 0x0F << 20 //+ L[3:0] bits (Regular channel sequence length).
	L_0    SQR1_Bits = 0x01 << 20 //  Bit 0.
	L_1    SQR1_Bits = 0x02 << 20 //  Bit 1.
	L_2    SQR1_Bits = 0x04 << 20 //  Bit 2.
	L_3    SQR1_Bits = 0x08 << 20 //  Bit 3.
)

const (
	SQ7    SQR2_Bits = 0x1F << 0  //+ SQ7[4:0] bits (7th conversion in regular sequence).
	SQ7_0  SQR2_Bits = 0x01 << 0  //  Bit 0.
	SQ7_1  SQR2_Bits = 0x02 << 0  //  Bit 1.
	SQ7_2  SQR2_Bits = 0x04 << 0  //  Bit 2.
	SQ7_3  SQR2_Bits = 0x08 << 0  //  Bit 3.
	SQ7_4  SQR2_Bits = 0x10 << 0  //  Bit 4.
	SQ8    SQR2_Bits = 0x1F << 5  //+ SQ8[4:0] bits (8th conversion in regular sequence).
	SQ8_0  SQR2_Bits = 0x01 << 5  //  Bit 0.
	SQ8_1  SQR2_Bits = 0x02 << 5  //  Bit 1.
	SQ8_2  SQR2_Bits = 0x04 << 5  //  Bit 2.
	SQ8_3  SQR2_Bits = 0x08 << 5  //  Bit 3.
	SQ8_4  SQR2_Bits = 0x10 << 5  //  Bit 4.
	SQ9    SQR2_Bits = 0x1F << 10 //+ SQ9[4:0] bits (9th conversion in regular sequence).
	SQ9_0  SQR2_Bits = 0x01 << 10 //  Bit 0.
	SQ9_1  SQR2_Bits = 0x02 << 10 //  Bit 1.
	SQ9_2  SQR2_Bits = 0x04 << 10 //  Bit 2.
	SQ9_3  SQR2_Bits = 0x08 << 10 //  Bit 3.
	SQ9_4  SQR2_Bits = 0x10 << 10 //  Bit 4.
	SQ10   SQR2_Bits = 0x1F << 15 //+ SQ10[4:0] bits (10th conversion in regular sequence).
	SQ10_0 SQR2_Bits = 0x01 << 15 //  Bit 0.
	SQ10_1 SQR2_Bits = 0x02 << 15 //  Bit 1.
	SQ10_2 SQR2_Bits = 0x04 << 15 //  Bit 2.
	SQ10_3 SQR2_Bits = 0x08 << 15 //  Bit 3.
	SQ10_4 SQR2_Bits = 0x10 << 15 //  Bit 4.
	SQ11   SQR2_Bits = 0x1F << 20 //+ SQ11[4:0] bits (11th conversion in regular sequence).
	SQ11_0 SQR2_Bits = 0x01 << 20 //  Bit 0.
	SQ11_1 SQR2_Bits = 0x02 << 20 //  Bit 1.
	SQ11_2 SQR2_Bits = 0x04 << 20 //  Bit 2.
	SQ11_3 SQR2_Bits = 0x08 << 20 //  Bit 3.
	SQ11_4 SQR2_Bits = 0x10 << 20 //  Bit 4.
	SQ12   SQR2_Bits = 0x1F << 25 //+ SQ12[4:0] bits (12th conversion in regular sequence).
	SQ12_0 SQR2_Bits = 0x01 << 25 //  Bit 0.
	SQ12_1 SQR2_Bits = 0x02 << 25 //  Bit 1.
	SQ12_2 SQR2_Bits = 0x04 << 25 //  Bit 2.
	SQ12_3 SQR2_Bits = 0x08 << 25 //  Bit 3.
	SQ12_4 SQR2_Bits = 0x10 << 25 //  Bit 4.
)

const (
	SQ1   SQR3_Bits = 0x1F << 0  //+ SQ1[4:0] bits (1st conversion in regular sequence).
	SQ1_0 SQR3_Bits = 0x01 << 0  //  Bit 0.
	SQ1_1 SQR3_Bits = 0x02 << 0  //  Bit 1.
	SQ1_2 SQR3_Bits = 0x04 << 0  //  Bit 2.
	SQ1_3 SQR3_Bits = 0x08 << 0  //  Bit 3.
	SQ1_4 SQR3_Bits = 0x10 << 0  //  Bit 4.
	SQ2   SQR3_Bits = 0x1F << 5  //+ SQ2[4:0] bits (2nd conversion in regular sequence).
	SQ2_0 SQR3_Bits = 0x01 << 5  //  Bit 0.
	SQ2_1 SQR3_Bits = 0x02 << 5  //  Bit 1.
	SQ2_2 SQR3_Bits = 0x04 << 5  //  Bit 2.
	SQ2_3 SQR3_Bits = 0x08 << 5  //  Bit 3.
	SQ2_4 SQR3_Bits = 0x10 << 5  //  Bit 4.
	SQ3   SQR3_Bits = 0x1F << 10 //+ SQ3[4:0] bits (3rd conversion in regular sequence).
	SQ3_0 SQR3_Bits = 0x01 << 10 //  Bit 0.
	SQ3_1 SQR3_Bits = 0x02 << 10 //  Bit 1.
	SQ3_2 SQR3_Bits = 0x04 << 10 //  Bit 2.
	SQ3_3 SQR3_Bits = 0x08 << 10 //  Bit 3.
	SQ3_4 SQR3_Bits = 0x10 << 10 //  Bit 4.
	SQ4   SQR3_Bits = 0x1F << 15 //+ SQ4[4:0] bits (4th conversion in regular sequence).
	SQ4_0 SQR3_Bits = 0x01 << 15 //  Bit 0.
	SQ4_1 SQR3_Bits = 0x02 << 15 //  Bit 1.
	SQ4_2 SQR3_Bits = 0x04 << 15 //  Bit 2.
	SQ4_3 SQR3_Bits = 0x08 << 15 //  Bit 3.
	SQ4_4 SQR3_Bits = 0x10 << 15 //  Bit 4.
	SQ5   SQR3_Bits = 0x1F << 20 //+ SQ5[4:0] bits (5th conversion in regular sequence).
	SQ5_0 SQR3_Bits = 0x01 << 20 //  Bit 0.
	SQ5_1 SQR3_Bits = 0x02 << 20 //  Bit 1.
	SQ5_2 SQR3_Bits = 0x04 << 20 //  Bit 2.
	SQ5_3 SQR3_Bits = 0x08 << 20 //  Bit 3.
	SQ5_4 SQR3_Bits = 0x10 << 20 //  Bit 4.
	SQ6   SQR3_Bits = 0x1F << 25 //+ SQ6[4:0] bits (6th conversion in regular sequence).
	SQ6_0 SQR3_Bits = 0x01 << 25 //  Bit 0.
	SQ6_1 SQR3_Bits = 0x02 << 25 //  Bit 1.
	SQ6_2 SQR3_Bits = 0x04 << 25 //  Bit 2.
	SQ6_3 SQR3_Bits = 0x08 << 25 //  Bit 3.
	SQ6_4 SQR3_Bits = 0x10 << 25 //  Bit 4.
)

const (
	JSQ1   JSQR_Bits = 0x1F << 0  //+ JSQ1[4:0] bits (1st conversion in injected sequence).
	JSQ1_0 JSQR_Bits = 0x01 << 0  //  Bit 0.
	JSQ1_1 JSQR_Bits = 0x02 << 0  //  Bit 1.
	JSQ1_2 JSQR_Bits = 0x04 << 0  //  Bit 2.
	JSQ1_3 JSQR_Bits = 0x08 << 0  //  Bit 3.
	JSQ1_4 JSQR_Bits = 0x10 << 0  //  Bit 4.
	JSQ2   JSQR_Bits = 0x1F << 5  //+ JSQ2[4:0] bits (2nd conversion in injected sequence).
	JSQ2_0 JSQR_Bits = 0x01 << 5  //  Bit 0.
	JSQ2_1 JSQR_Bits = 0x02 << 5  //  Bit 1.
	JSQ2_2 JSQR_Bits = 0x04 << 5  //  Bit 2.
	JSQ2_3 JSQR_Bits = 0x08 << 5  //  Bit 3.
	JSQ2_4 JSQR_Bits = 0x10 << 5  //  Bit 4.
	JSQ3   JSQR_Bits = 0x1F << 10 //+ JSQ3[4:0] bits (3rd conversion in injected sequence).
	JSQ3_0 JSQR_Bits = 0x01 << 10 //  Bit 0.
	JSQ3_1 JSQR_Bits = 0x02 << 10 //  Bit 1.
	JSQ3_2 JSQR_Bits = 0x04 << 10 //  Bit 2.
	JSQ3_3 JSQR_Bits = 0x08 << 10 //  Bit 3.
	JSQ3_4 JSQR_Bits = 0x10 << 10 //  Bit 4.
	JSQ4   JSQR_Bits = 0x1F << 15 //+ JSQ4[4:0] bits (4th conversion in injected sequence).
	JSQ4_0 JSQR_Bits = 0x01 << 15 //  Bit 0.
	JSQ4_1 JSQR_Bits = 0x02 << 15 //  Bit 1.
	JSQ4_2 JSQR_Bits = 0x04 << 15 //  Bit 2.
	JSQ4_3 JSQR_Bits = 0x08 << 15 //  Bit 3.
	JSQ4_4 JSQR_Bits = 0x10 << 15 //  Bit 4.
	JL     JSQR_Bits = 0x03 << 20 //+ JL[1:0] bits (Injected Sequence length).
	JL_0   JSQR_Bits = 0x01 << 20 //  Bit 0.
	JL_1   JSQR_Bits = 0x02 << 20 //  Bit 1.
)

const (
	JDATA JDR1_Bits = 0xFFFF << 0 //+ Injected data.
)

const (
	JDATA JDR2_Bits = 0xFFFF << 0 //+ Injected data.
)

const (
	JDATA JDR3_Bits = 0xFFFF << 0 //+ Injected data.
)

const (
	JDATA JDR4_Bits = 0xFFFF << 0 //+ Injected data.
)

const (
	DATA     DR_Bits = 0xFFFF << 0  //+ Regular data.
	ADC2DATA DR_Bits = 0xFFFF << 16 //+ ADC2 data.
)
