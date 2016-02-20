package gpio

import (
	"mmio"

	"arch/cortexm/bitband"
)

func bit(p *Port, reg *mmio.U32) bitband.Bit {
	return bitband.Alias32(reg).Bit(pnum(p))
}
