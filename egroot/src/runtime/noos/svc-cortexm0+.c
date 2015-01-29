// +build cortexm0

__attribute__ ((naked)) static
void runtime$noos$svcHandler() {
	asm volatile (
		"mov    r0, lr\n\t"
		"movs   r1, #4\n\t"
		"tst	r0, r1\n\t"
		"bne    0f\n\t"

		"mrs 	r0, msp\n\t"
		"b		runtime$noos$sv\n\t"

		"0:\n\t"

		"mrs	r0, psp\n\t"
		"b		runtime$noos$sv"

		:: "X" (runtime$noos$sv)
	);
}
