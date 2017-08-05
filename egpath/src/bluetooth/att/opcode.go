package att

type Method byte

const (
	ExchangeMTU     Method = 0x02 >> 1
	FindInformation Method = 0x04 >> 1
	FindByTypeValue Method = 0x06 >> 1
	ReadByType      Method = 0x08 >> 1
	Read            Method = 0x0A >> 1
	ReadBlob        Method = 0x0C >> 1

	ReadByGroup Method = 0x10 >> 1
)
