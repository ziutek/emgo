// +build l476xx

package dma

const (
	DMA1_ADC1   Request = 0
	DMA1_ADC2   Request = 0
	DMA1_ADC3   Request = 0
	DMA1_DFSDM1 Request = 0

	DMA1_SAI2 Request = 1
	DMA1_SPI1 Request = 1
	DMA1_SPI2 Request = 1

	DMA1_USART1 Request = 2
	DMA1_USART2 Request = 2
	DMA1_USART3 Request = 2

	DMA1_I2C1 Request = 3
	DMA1_I2C2 Request = 3
	DMA1_I2C3 Request = 3

	DMA1_TIM2  Request = 4
	DMA1_TIM16 Request = 4

	DMA1_QUADSPI Request = 5
	DMA1_TIM3    Request = 5
	DMA1_TIM7    Request = 5
	DMA1_TIM17   Request = 5

	DMA1_DAC1 Request = 6
	DMA1_TIM4 Request = 6
	DMA1_TIM6 Request = 6

	DMA1_TIM1  Request = 7
	DMA1_TIM15 Request = 7
)

const (
	DMA2_ADC1     Request = 0
	DMA2_ADC2     Request = 0
	DMA2_ADC3     Request = 0
	DMA2_DCMI_CH6 Request = 0
	DMA2_I2C4     Request = 0

	DMA2_SAI1 Request = 1
	DMA2_SAI2 Request = 1

	DMA2_UART4  Request = 2
	DMA2_UART5  Request = 2
	DMA2_USART1 Request = 2

	DMA2_DAC1    Request = 3
	DMA2_DAC2    Request = 3
	DMA2_QUADSPI Request = 3
	DMA2_SPI3    Request = 3
	DMA2_TIM6    Request = 3
	DMA2_TIM7    Request = 3

	DMA2_DCMI_CH5 Request = 4
	DMA2_LPUART1  Request = 4
	DMA2_SPI1     Request = 4
	DMA2_SWPMI1   Request = 4

	DMA2_I2C1 Request = 5
	DMA2_TIM5 Request = 5

	DMA2_AES  Request = 6
	DMA2_HASH Request = 6

	DMA2_SDMMC1 Request = 7
	DMA2_TIM8   Request = 7
)
