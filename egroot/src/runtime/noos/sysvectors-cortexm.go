// +build cortexm0 cortexm3 cortexm4 cortexm4f

package noos

import "arch/cortexm/exce"

func Reset()

//emgo:export
//c:const
//c:__attribute__((section(".SystemVectors")))
var sysvectors = [...]func(){
	exce.Reset - 1:      Reset,
	exce.NMI - 1:        nmiHandler,
	exce.HardFault - 1:  FaultHandler,
	exce.MemManage - 1:  FaultHandler,
	exce.BusFault - 1:   FaultHandler,
	exce.UsageFault - 1: FaultHandler,
	exce.SVC - 1:        svcHandler,
	exce.PendSV - 1:     pendSVHandler,
	exce.SysTick - 1:    sysTickHandler,
}
