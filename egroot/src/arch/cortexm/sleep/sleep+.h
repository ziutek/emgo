static inline
void arch$cortexm$sleep$WFE() {
	asm volatile ("wfe");
}

static inline
void arch$cortexm$sleep$WFI() {
	asm volatile ("wfi");
}
