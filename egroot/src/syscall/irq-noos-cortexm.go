// +build noos
// +build cortexm0 cortexm3 cortexm4 cortexm4f

package syscall

import (
	"arch/cortexm/exce"
)

const (
	IRQPrioLowest = exce.Lowest
	IRQPrioHighest = exce.Highest
)