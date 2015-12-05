// +build cortexm0 cortexm3 cortexm4 cortexm4f

.syntax unified

.global runtime$noos$faultHandler

.thumb_func
runtime$noos$faultHandler:
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
// 3: HardFault  - see HFSR: x/xw 0xE000ED2C
// 4: MemManage  - see MMSR: x/xb 0xE000ED28, MMAR: x 0xE000ED34
// 5: BusFault   - see BFSR: x/xb 0xE000ED29, BFAR: x 0xE000ED38
// 6: UsageFault - see UFSR: x/xh 0xE000ED2A

// R1 should contain pointer to the exception stack frame:
// (R1) -> [R0, R1, R2, R3, IP, LR, PC, PSR]
// If R1 points to valid memory examine:
// 1. Where PC points.
// 2. Thumb bit in PSR
// 3. IPSR in PSR

// To print stack frame in gdb use: p /x *(arch$cortexm$StackFrame*)($r1)
// To see line where PC points to: b * ((arch$cortexm$StackFrame*)($r1))->PC
