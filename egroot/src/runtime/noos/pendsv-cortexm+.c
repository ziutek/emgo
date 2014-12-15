// +build cortexm3 cortexm4

__attribute__ ((naked)) static
void runtime$noos$pendSVHandler() {
	asm volatile (
		"push	{lr}\n\t"
		"mrs	r0, psp\n\t"
		
		// Call nextTask with SP used by current task.
		"bl		runtime$noos$nextTask\n\t"
		
		// Check wheater the context switch is need (r0 contains taskInfo.sp
		// for next task or 0 if context switch isn't need).
		"cbz	r0, 1f\n\t"
		
		// Save the second part of the context of the current task.
		"mrs	r1, psp\n\t"
		"stmdb	r1, {r4-r11}\n\t"
		
		// Restore the second part of the context of the next task.
		"msr	psp, r0\n\t"
		"subs	r0, 8*4\n\t" 
		"ldm	r0, {r4-r11}\n\t"
		
		"1:\n\t"
	
		"pop	{pc}"
		:: "X" (runtime$noos$nextTask)
	);
}
