package main

import (
	"fmt"

	"stm32/hal/system"
	"stm32/hal/system/timer/systick"

	"stm32/hal/raw/rcc"
	"stm32/hal/raw/tim"
)

func init() {
	system.Setup168(8)
	systick.Setup(2e6)
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
