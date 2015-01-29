// +build cortexm0 cortexm3 cortexm4 cortexm4f

.syntax unified

.global runtime$noos$FaultHandler

.thumb_func
runtime$noos$FaultHandler:
	// At this point a lot of things can be broken so don't touch
	// stack nor memory. Do only few things that helps debuging.
	mov   r0, lr
	movs  r1, #4
	tst   r0, r1
	bne   0f

	mrs  r1, msp
	b    1f
0:
	mrs  r1, psp
1:
	mrs   r0, ipsr
2:  bkpt  1
	b     2b

// Now R0 and R1 contain useful information.

// R0 contains exception number:
// 3: HardFault  - see arch$cortexm$exce$FSR->HFS
// 4: MemManage  - see arch$cortexm$exce$FSR->MMS, arch$cortexm$exce$FSR->MMA
// 5: BusFault   - see arch$cortexm$exce$FSR->BFS, arch$cortexm$exce$FSR->BFA
// 6: UsageFault - see arch$cortexm$exce$FSR->UFS

// R1 should contain pointer to the exception stack frame:
// (R1) -> [R0, R1, R2, R3, IP, LR, PC, PSR]
// If R1 points to valid memory examine:
// 1. Where PC points.
// 2. Thumb bit in PSR
// 3. IPSR in PSR

// To print stack frame in gdb use: p /x *(arch$cortexm$exce$StackFrame*)($r1)
// To see line where PC points to: b * ((arch$cortexm$exce$StackFrame*)($r1))->PC
