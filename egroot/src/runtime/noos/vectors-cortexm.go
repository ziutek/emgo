// +build cortexm0 cortexm3 cortexm4 cortexm4f

package noos

func reset()

// Vectors for system exceptions.
//emgo:export
var (
	ev1  = reset          //c:__attribute__((section(".Reset")))
	ev2  = nmiHandler     //c:__attribute__((section(".NMI")))
	ev3  = faultHandler   //c:__attribute__((section(".HardFault")))
	ev4  = faultHandler   //c:__attribute__((section(".MemManage")))
	ev5  = faultHandler   //c:__attribute__((section(".BusFault")))
	ev6  = faultHandler   //c:__attribute__((section(".UsageFault")))
	ev7  = faultHandler   //c:__attribute__((section(".Reserved7")))
	ev8  = faultHandler   //c:__attribute__((section(".Reserved8")))
	ev9  = faultHandler   //c:__attribute__((section(".Reserved9")))
	ev10 = faultHandler   //c:__attribute__((section(".Reserved10")))
	ev11 = svcHandler     //c:__attribute__((section(".SVCall")))
	ev12 = faultHandler   //c:__attribute__((section(".DebugMon")))
	ev13 = faultHandler   //c:__attribute__((section(".Reserved13")))
	ev14 = pendSVHandler  //c:__attribute__((section(".PendSV")))
)

	//ev15 = sysTickHandler //c:__attribute__((section(".SysTick")))
