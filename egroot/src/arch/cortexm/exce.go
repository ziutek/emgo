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
	PrioStep    = -1
	PrioNum     = 256
)

// EXC_RETURN
const (
	// ReturnMask: selects bits that can be compared with Return* constants.
	ReturnMask = 0xf

	// ReturnHandler: exception will return to handler mode.
	ReturnHandler = 0x1

	// ReturnMSP: exception will return to thread mode using MSP .
	ReturnMSP = 0x9

	// ReturnPSP: exception will return to thread mode using PSP.
	ReturnPSP = 0xd

	// BasicFrame: if set there is no floating-point state saved on stack.
	BasicFrame = 0x10
)
