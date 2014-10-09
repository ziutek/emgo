package gpio

type AltFunc byte

const (
	AF0 AltFunc = iota
	AF1
	AF2
	AF3
	AF4
	AF5
	AF6
	AF7
	AF8
	AF9
	AF10
	AF11
	AF12
	AF13
	AF14
	AF15

	Sys = AF0

	Tim2 = AF1

	Tim3 = AF2
	Tim4 = AF2
	Tim5 = AF2

	Tim9  = AF3
	Tim10 = AF3
	Tim11 = AF3

	I2C1 = AF4
	I2C2 = AF4

	SPI1 = AF5
	SPI2 = AF5

	SPI3 = AF6

	USART1 = AF7
	USART2 = AF7
	USART3 = AF7

	UART4 = AF8
	UART5 = AF8

	USB = AF10

	LCD = AF11

	FSMC = AF12

	RI = AF14

	EventOut = AF15
)
