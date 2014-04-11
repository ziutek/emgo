// +build cortexm3 cortexm4 cortexm4f

__attribute__ ((naked)) static
void runtime_noos_svcHandler() {
	asm volatile (
		"tst	lr, 4\n\t"
		"ite 	eq\n\t"
		"mrseq 	r0, msp\n\t"
		"mrsne	r0, psp\n\t"
		"b		runtime_noos_sv"
		:: "X" (runtime_noos_sv)
	);
}