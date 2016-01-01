package main

import (
	"fmt"

	"stm32/hal/raw/rcc"
	"stm32/hal/raw/tim"
	"stm32/hal/setup"
)

func init() {
	setup.Performance168(8)
	rcc.RCC.TIM10EN().Set()
}

func main() {
	tim := tim.TIM10
	tim.ARR.Store(65000)
	tim.CNT.Store(0)
	tim.CEN().Set()
	for {
		fmt.Println(tim.CNT.Load())
	}
}
