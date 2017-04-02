inline int_
debug$semihosting$hostIO(int_ cmd, unsafe$Pointer p) {
	register int_ r0 asm("r0") = cmd;
	register unsafe$Pointer r1 asm("r1") = p;
	asm volatile (
		"bkpt 0xAB"
		:"+r" (r0)
		:"r" (r1)
		:"memory"
	);
	return r0;
}
