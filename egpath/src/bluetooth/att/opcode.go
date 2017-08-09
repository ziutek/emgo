package att

type Method byte

const (
	ExchangeMTU     Method = 0x02
	FindInformation Method = 0x04
	FindByTypeValue Method = 0x06

	ReadByType      Method = 0x08
	Read            Method = 0x0A
	ReadBlob        Method = 0x0C
	ReadMultiple    Method = 0x0E
	ReadByGroupType Method = 0x10

	Write        Method = 0x12
	unusedMethod Method = 0x14
	PrepareWrite Method = 0x16
	ExecuteWrite Method = 0x18
)

