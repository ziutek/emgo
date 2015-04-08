// +build cortexm0 cortexm3 cortexm4 cortexm4f
package noos

import "arch/cortexm"

func panic_(i interface{}) {
	for {
		cortexm.BKPT(1)
	}
}
