package exce

// StackFrame represents Cortex-M exception stack frame (without floating-point state).
type StackFrame struct {
	R   [4]uintptr
	IP  uintptr
	LR  uintptr
	PC  uintptr
	PSR uint32
}

// StackFrameFP represents Cortex-M stack frame including floating-point state.
type StackFrameFP struct {
	StackFrame
	S     [16]float32
	FPSCR uint32
}
