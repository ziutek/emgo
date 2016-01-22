// +build noos
// +build cortexm0 cortexm3 cortexm4 cortexm4f

package syscall

import (
	"arch/cortexm"
)

const (
	IRQPrioLowest  = cortexm.PrioLowest
	IRQPrioHighest = cortexm.PrioHighest
	IRQPrioStep    = cortexm.PrioStep
	IRQPrioNum     = cortexm.PrioNum
)
