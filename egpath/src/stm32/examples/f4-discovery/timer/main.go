package main

import (
	"fmt"

	"stm32/f4/periph"
	"stm32/f4/setup"
	"stm32/timer"
)

func init() {
	setup.Performance168(8)
	periph.APB2ClockEnable(periph.TIM10)
	periph.APB2Reset(periph.TIM10)
}

func main() {
	tim := timer.TIM10
	tim.ARR_Store(65000)
	tim.CNT_Store(0)
	tim.CEN().Set()
	for {
		fmt.Println(tim.CNT_Load())
	}
}
