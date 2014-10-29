// +build cortexm0 cortexm3 cortexm4 cortexm4f

.syntax unified

.global runtime$noos$FaultHandler

.thumb_func
runtime$noos$FaultHandler:
	// At this point a lot of things can be broken so don't touch
	// stack nor memory. Do only few things that helps debuging.
	mrs    r0, ipsr
	tst    lr, 4
	ite    eq
	mrseq  r1, msp
	mrsne  r1, psp
	bkpt   1

// Now R0 and R1 contain useful information.

// R0 contains exception number:
// 3: HardFault  - see cortexm$exce$FSR->HFS
// 4: MemManage  - see cortexm$exce$FSR->MMS, cortexm$exce$FSR->MMA
// 5: BusFault   - see cortexm$exce$FSR->BFS, cortexm$exce$FSR->BFA
// 6: UsageFault - see cortexm$exce$FSR->UFS

// R1 should contain pointer to the exception stack frame:
// (R1) -> [R0, R1, R2, R3, IP, LR, PC, PSR]
// If R1 points to valid memory examine:
// 1. Where PC points.
// 2. Thumb bit in PSR
// 3. IPSR in PSR

// To print stack frame in gdb use: p /x *(cortexm$exce$StackFrame*)($r1)
// To see line where PC points to: b * ((cortexm$exce$StackFrame*)($r1))->PC
