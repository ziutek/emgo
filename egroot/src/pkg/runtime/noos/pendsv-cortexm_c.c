// +build cortexm0 cortexm3 cortexm4 cortexm4f

__attribute__ ((naked))
static void runtime_noos_pendSVHandler() {
	asm volatile (
		// Call nextTask with PSP used by current task.
		"mrs	r0, psp\n\t"
		"bl		runtime_noos_nextTask\n\t"
		// Check wheater a context switch is need.
		"cbz	r0, 1f\n\t"
		
		"mrs	r1, psp\n\t"
		"stmdb	r1, {r4-r11}\n\t"
		// nextTask returns PSP before stacking.
		"msr	psp, r0\n\t"
		"subs	r0, #32\n\t" 
		"ldmia	r0, {r4-r11}\n\t"
		
		"1: bx	lr"
	);
}