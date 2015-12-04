// +build cortexm0

__attribute__ ((naked)) static
void runtime$noos$pendSVHandler() {
	asm volatile (
		"push	{lr}\n\t"
		"mrs	r0, psp\n\t"
		
		// Call nextTask with SP used by current task.
		"bl		runtime$noos$nextTask\n\t"
		
		// Check wheater the context switch is need (r0 contains taskInfo.sp
		// for next task or 0 if context switch isn't need).
		"tst	r0, r0\n\t"
		"beq	1f\n\t"
		
		// Save the second part of the context of the current task.
		"mrs	r1, psp\n\t"
		"sub	r1, #8*4\n\t"
		"stmia	r1!, {r4-r7}\n\t"
		"mov    r4, r8\n\t"
		"mov    r5, r9\n\t"
		"mov    r6, r10\n\t"
		"mov    r7, r11\n\t"
		"stmia	r1!, {r4-r7}\n\t"	
		
		// Restore the second part of the context of the next task.
		"msr	psp, r0\n\t"
		"sub	r0, #4*4\n\t" 
		"ldmia	r0!, {r4-r7}\n\t"
		"mov    r8, r4\n\t"
		"mov    r9, r5\n\t"
		"mov    r10, r6\n\t"
		"mov    r11, r7\n\t"
		"sub	r0, #8*4\n\t"
		"ldmia	r0!, {r4-r7}\n\t"	
		
		"1:\n\t"
	
		"pop	{pc}"
		:: "X" (runtime$noos$nextTask)
	);
}
