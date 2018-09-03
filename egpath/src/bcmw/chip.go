package bcmw

// The chip config concept borrowed from NuttX (http://nuttx.org/)

type Chip struct {
	ramSize  int
	baseAddr [5]uint32
}

const (
	coreCommon   = 0
	coreDot11MAC = 1
	coreSDIO     = 2
	coreARMCM3   = 3
	coreSOCSRAM  = 4
)

//emgo:const
var (
	chip43362 = Chip{
		ramSize: 240 * 1024,
		baseAddr: [5]uint32{
			coreCommon:   0x18000000,
			coreDot11MAC: 0x18001000,
			coreSDIO:     0x18002000,
			coreARMCM3:   0x18003000 + 0x100000,
			coreSOCSRAM:  0x18004000 + 0x100000,
		},
	}
	chip43438 = Chip{
		ramSize: 512 * 1024,
		baseAddr: [5]uint32{
			coreCommon:   0x18000000,
			coreDot11MAC: 0x18001000,
			coreSDIO:     0x18002000,
			coreARMCM3:   0x18003000 + 0x100000,
			coreSOCSRAM:  0x18004000 + 0x100000,
		},
	}
)

//emgo:const
var (
	Chip43362 = &chip43362
	Chip43438 = &chip43438
)
