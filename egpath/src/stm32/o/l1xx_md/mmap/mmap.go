// Package mmap provides base memory adresses for all peripherals.
package mmap

const (
	FLASH_BASE     uintptr = 0x08000000 // FLASH base address in the alias region
	SRAM_BASE      uintptr = 0x20000000 // SRAM base address in the alias region
	PERIPH_BASE    uintptr = 0x40000000 // Peripheral base address in the alias region
	SRAM_BB_BASE   uintptr = 0x22000000 // SRAM base address in the bit-band region
	PERIPH_BB_BASE uintptr = 0x42000000 // Peripheral base address in the bit-band region
	FSMC_R_BASE    uintptr = 0xA0000000 // FSMC registers base address
)

// Peripheral memory map
const (
	APB1PERIPH_BASE    uintptr = PERIPH_BASE
	APB2PERIPH_BASE    uintptr = PERIPH_BASE + 0x10000
	AHBPERIPH_BASE     uintptr = PERIPH_BASE + 0x20000
	TIM2_BASE          uintptr = APB1PERIPH_BASE + 0x0000
	TIM3_BASE          uintptr = APB1PERIPH_BASE + 0x0400
	TIM4_BASE          uintptr = APB1PERIPH_BASE + 0x0800
	TIM5_BASE          uintptr = APB1PERIPH_BASE + 0x0C00
	TIM6_BASE          uintptr = APB1PERIPH_BASE + 0x1000
	TIM7_BASE          uintptr = APB1PERIPH_BASE + 0x1400
	LCD_BASE           uintptr = APB1PERIPH_BASE + 0x2400
	RTC_BASE           uintptr = APB1PERIPH_BASE + 0x2800
	WWDG_BASE          uintptr = APB1PERIPH_BASE + 0x2C00
	IWDG_BASE          uintptr = APB1PERIPH_BASE + 0x3000
	SPI2_BASE          uintptr = APB1PERIPH_BASE + 0x3800
	SPI3_BASE          uintptr = APB1PERIPH_BASE + 0x3C00
	USART2_BASE        uintptr = APB1PERIPH_BASE + 0x4400
	USART3_BASE        uintptr = APB1PERIPH_BASE + 0x4800
	UART4_BASE         uintptr = APB1PERIPH_BASE + 0x4C00
	UART5_BASE         uintptr = APB1PERIPH_BASE + 0x5000
	I2C1_BASE          uintptr = APB1PERIPH_BASE + 0x5400
	I2C2_BASE          uintptr = APB1PERIPH_BASE + 0x5800
	PWR_BASE           uintptr = APB1PERIPH_BASE + 0x7000
	DAC_BASE           uintptr = APB1PERIPH_BASE + 0x7400
	COMP_BASE          uintptr = APB1PERIPH_BASE + 0x7C00
	RI_BASE            uintptr = APB1PERIPH_BASE + 0x7C04
	OPAMP_BASE         uintptr = APB1PERIPH_BASE + 0x7C5C
	SYSCFG_BASE        uintptr = APB2PERIPH_BASE + 0x0000
	EXTI_BASE          uintptr = APB2PERIPH_BASE + 0x0400
	TIM9_BASE          uintptr = APB2PERIPH_BASE + 0x0800
	TIM10_BASE         uintptr = APB2PERIPH_BASE + 0x0C00
	TIM11_BASE         uintptr = APB2PERIPH_BASE + 0x1000
	ADC1_BASE          uintptr = APB2PERIPH_BASE + 0x2400
	ADC_BASE           uintptr = APB2PERIPH_BASE + 0x2700
	SDIO_BASE          uintptr = APB2PERIPH_BASE + 0x2C00
	SPI1_BASE          uintptr = APB2PERIPH_BASE + 0x3000
	USART1_BASE        uintptr = APB2PERIPH_BASE + 0x3800
	GPIOA_BASE         uintptr = AHBPERIPH_BASE + 0x0000
	GPIOB_BASE         uintptr = AHBPERIPH_BASE + 0x0400
	GPIOC_BASE         uintptr = AHBPERIPH_BASE + 0x0800
	GPIOD_BASE         uintptr = AHBPERIPH_BASE + 0x0C00
	GPIOE_BASE         uintptr = AHBPERIPH_BASE + 0x1000
	GPIOH_BASE         uintptr = AHBPERIPH_BASE + 0x1400
	GPIOF_BASE         uintptr = AHBPERIPH_BASE + 0x1800
	GPIOG_BASE         uintptr = AHBPERIPH_BASE + 0x1C00
	CRC_BASE           uintptr = AHBPERIPH_BASE + 0x3000
	RCC_BASE           uintptr = AHBPERIPH_BASE + 0x3800
	FLASH_R_BASE       uintptr = AHBPERIPH_BASE + 0x3C00 // FLASH registers base address
	OB_BASE            uintptr = 0x1FF80000              // FLASH Option Bytes base address
	DMA1_BASE          uintptr = AHBPERIPH_BASE + 0x6000
	DMA1_Channel1_BASE uintptr = DMA1_BASE + 0x0008
	DMA1_Channel2_BASE uintptr = DMA1_BASE + 0x001C
	DMA1_Channel3_BASE uintptr = DMA1_BASE + 0x0030
	DMA1_Channel4_BASE uintptr = DMA1_BASE + 0x0044
	DMA1_Channel5_BASE uintptr = DMA1_BASE + 0x0058
	DMA1_Channel6_BASE uintptr = DMA1_BASE + 0x006C
	DMA1_Channel7_BASE uintptr = DMA1_BASE + 0x0080
	DMA2_BASE          uintptr = AHBPERIPH_BASE + 0x6400
	DMA2_Channel1_BASE uintptr = DMA2_BASE + 0x0008
	DMA2_Channel2_BASE uintptr = DMA2_BASE + 0x001C
	DMA2_Channel3_BASE uintptr = DMA2_BASE + 0x0030
	DMA2_Channel4_BASE uintptr = DMA2_BASE + 0x0044
	DMA2_Channel5_BASE uintptr = DMA2_BASE + 0x0058
	AES_BASE           uintptr = 0x50060000
	FSMC_Bank1_R_BASE  uintptr = FSMC_R_BASE + 0x0000 // FSMC Bank1 registers base address
	FSMC_Bank1E_R_BASE uintptr = FSMC_R_BASE + 0x0104 // FSMC Bank1E registers base address
	DBGMCU_BASE        uintptr = 0xE0042000           // Debug MCU registers base address
)
