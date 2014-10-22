// +build cortexm0 cortexm3 cortexm4 cortexm4f

__attribute__ ((naked)) static
void runtime$noos$hardFaultHandler() {
	// At this point a lot of things can be broken so don't touch stack nor
	// memory. Do only few things that helps debuging. 
	asm volatile (
		"tst	lr, 4\n\t"
		"ite 	eq\n\t"
		"mrseq 	r0, msp\n\t"
		"mrsne	r0, psp\n\t"
		"bkpt	1"
	);
	// Now R0 should contain pointer to the exception stack frame:
	// (R0) -> [R0, R1, R2, R3, IP, LR, PC, xPSR]
	// If R0 points to valid memory examine:
	// 1. Where PC points.
	// 2. Thumb bit in xPSR
	// 3. IPSR in xPSR
	// See also two HFSR bits at address 0xE000ED2C:
	//  1 - BusFault on vector table read,
	// 30 - escalated, see *runtime$noos$cfsr for more info.
}
