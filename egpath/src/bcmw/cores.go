package bcmw

const (
	coreCommon   = 0
	coreDot11MAC = 1
	coreSDIO     = 2
	coreARMCM3   = 3
	coreSOCSRAM  = 4
)

//emgo:const
var wrapBase = [5]uint32{
	coreCommon:   0x18000000,
	coreDot11MAC: 0x18001000,
	coreSDIO:     0x18002000,
	coreARMCM3:   0x18003000 + 0x100000,
	coreSOCSRAM:  0x18004000 + 0x100000,
}
