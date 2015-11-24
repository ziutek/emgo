package cortexm

// Cortex-M exception numbers.
const (
	Reset      = 1
	NMI        = 2
	HardFault  = 3
	MemManage  = 4
	BusFault   = 5
	UsageFault = 6
	SVCall     = 11
	DebugMon   = 12
	PendSV     = 14
	SysTick    = 15
)

// Exception number for first external interrupt.
const IRQ0 = 16

// Lowest and highest exception priority.
const (
	PrioLowest  = 255
	PrioHighest = 0
)
