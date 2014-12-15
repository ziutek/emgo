__attribute__ ((always_inline))
extern inline
void arch$cortexm$sleep$WFE() {
	asm volatile ("wfe");
}

__attribute__ ((always_inline))
extern inline
void arch$cortexm$sleep$WFI() {
	asm volatile ("wfi");
}
