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

	Tim1 = AF1
	Tim2 = AF1

	Tim3 = AF2
	Tim4 = AF2
	Tim5 = AF2

	Tim8  = AF3
	Tim9  = AF3
	Tim10 = AF3
	Tim11 = AF3

	I2C1 = AF4
	I2C2 = AF4
	I2C3 = AF4

	SPI1 = AF5
	SPI2 = AF5
	SPI4 = AF5
	SPI5 = AF5
	SPI6 = AF5

	SPI3 = AF6
	SAI1 = AF6

	USART1 = AF7
	USART2 = AF7
	USART3 = AF7

	USART4 = AF8
	USART5 = AF8
	USART6 = AF8
	USART7 = AF8
	USART8 = AF8

	CAN1  = AF9
	CAN2  = AF9
	Tim12 = AF9
	Tim13 = AF9
	Tim14 = AF9

	OTGFS = AF10
	OTGHS = AF10

	Eth = AF11

	FSMC    = AF12
	FMC     = AF12
	SDIO    = AF12
	OTGHSFS = AF12

	DCMI = AF13

	LTDC = AF14

	EventOut = AF15
)

// AltFunc
func (g *Port) AltFunc(n int) AltFunc {
	var af uint32
	if n < 8 {
		af = g.afl
	} else {
		af = g.afh
		n -= 8
	}
	n *= 4
	return AltFunc(af>>uint(n)) & 0xf
}

// SetAltFunc
func (g *Port) SetAltFunc(n int, af AltFunc) {
	n *= 4
	if n < 32 {
		g.afl = g.afl&^(0xf<<uint(n)) | uint32(af)<<uint(n)
	} else {
		n -= 32
		g.afh = g.afh&^(0xf<<uint(n)) | uint32(af)<<uint(n)
	}
}
