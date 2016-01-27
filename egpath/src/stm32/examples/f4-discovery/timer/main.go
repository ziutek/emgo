package main

import (
	"fmt"

	"stm32/hal/osclk/systick"
	"stm32/hal/system"

	"stm32/hal/raw/rcc"
	"stm32/hal/raw/tim"
)

func init() {
	system.Setup168(8)
	systick.Setup()
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
