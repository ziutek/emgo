// +build cortexm0 cortexm3 cortexm4 cortexm4f

package noos

import "arch/cortexm"

func Reset()

//emgo:export
//c:const
//c:__attribute__((section(".SystemVectors")))
var sysvectors = [...]func(){
	cortexm.Reset - 1:      Reset,
	cortexm.NMI - 1:        nmiHandler,
	cortexm.HardFault - 1:  FaultHandler,
	cortexm.MemManage - 1:  FaultHandler,
	cortexm.BusFault - 1:   FaultHandler,
	cortexm.UsageFault - 1: FaultHandler,
	cortexm.SVCall - 1:     svcHandler,
	cortexm.PendSV - 1:     pendSVHandler,
	cortexm.SysTick - 1:    sysTickHandler,
}
