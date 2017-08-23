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

	Write         Method = 0x12
	unknownMethod Method = 0x14
	PrepareWrite  Method = 0x16
	ExecuteWrite  Method = 0x18
)

//emgo:const
var methodStr = [...]string{
	ExchangeMTU>>1 - 1:     "ExchangeMTU",
	FindInformation>>1 - 1: "FindInformation",
	FindByTypeValue>>1 - 1: "FindByTypeValue",
	ReadByType>>1 - 1:      "ReadByType",
	Read>>1 - 1:            "Read",
	ReadBlob>>1 - 1:        "ReadBlob",
	ReadMultiple>>1 - 1:    "ReadMultiple",
	ReadByGroupType>>1 - 1: "ReadByGroupType",
	Write>>1 - 1:           "Write",
	unknownMethod>>1 - 1:   "unknownMethod",
	PrepareWrite>>1 - 1:    "PrepareWrite",
	ExecuteWrite>>1 - 1:    "ExecuteWrite",
}

func (m Method) String() string {
	n := int(m)>>1 - 1
	if n >= len(methodStr) {
		return methodStr[unknownMethod>>1-1]
	}
	return methodStr[n]
}
