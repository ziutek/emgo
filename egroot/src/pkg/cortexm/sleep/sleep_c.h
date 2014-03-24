__attribute__ ((always_inline))
extern inline void cortexm_sleep_WFE() {
	asm volatile ("wfe");
}

__attribute__ ((always_inline))
extern inline void cortexm_sleep_WFI() {
	asm volatile ("wfi");
}
