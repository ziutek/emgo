// +build cortexm4f

__attribute__ ((naked)) static
void runtime$noos$pendSVHandler() {
	asm volatile (
		"push	{lr}\n\t"
		"mrs	r0, psp\n\t"
		
		// Set FPU bit in SP.
		"tst	lr, #0x10\n\t"
		"it		eq\n\t"
		"addeq	r0, #1\n\t"
		
		// Call nextTask with SP used by current task.
		"bl		runtime$noos$nextTask\n\t"
		
		"pop	{r2}\n\t"
		
		// Check whether the context switch is need (r0 contains taskInfo.sp
		// for next task or 0 if context switch isn't need).
		"cbz	r0, 1f\n\t"
		
		// Save the second part of the CPU context of the current task.
		"mrs	r1, psp\n\t"
		"stmdb	r1!, {r4-r11}\n\t" // Don't put this in IT block!
		
		"tst	r2, #0x10\n\t"
		"bne	0f\n\t"
		
		// Save the second part of the FPU context of the current task.
		"vstmdb	r1!, {d8-d15}\n\t" // Don't put this in IT block!
		
		"0:\n\t"
		
		// Adjust new SP and EXC_RETURN according to the FPU bit in taskInfo.sp.
		"tst	r0, #1\n\t"
		"itte	ne\n\t"
		"subne	r0, #1\n\t"
		"bicne	r2, #0x10\n\t"
		"orreq	r2, #0x10\n\t"
		
		"msr	psp, r0\n\t"
		
		// Restore the second part of the CPU context of the next task.
		"ldmdb	r0!, {r4-r11}\n\t" // Don't put this in IT block!
		
		"tst	r2, #0x10\n\t"
		"bne	1f\n\t"
		
		// Restore the second part of the FPU context of the next task.
		"vldmdb r0!, {d8-d15}\n\t" // Don't put this in IT block!
		
		"1:\n\t"
					
		"bx		r2"
		:: "X" (runtime$noos$nextTask)
	);
}

