__attribute__ ((always_inline))
extern inline void cortexm_irq_Disable() {
	asm volatile ("cpsid i");
}

__attribute__ ((always_inline))
extern inline void cortexm_irq_Enable() {
	asm volatile ("cpsie i");
}

__attribute__ ((always_inline))
extern inline bool cortexm_irq_Disabled() {
	bool b;
	asm volatile ("msr primask, %0" : "=r" (b));
	return b;
}

__attribute__ ((always_inline))
extern inline void cortexm_irq_DisableFaults() {
	asm volatile ("cpsid f");
}

__attribute__ ((always_inline))
extern inline void cortexm_irq_EnableFaults() {
	asm volatile ("cpsie f");
}

__attribute__ ((always_inline))
extern inline bool cortexm_irq_FaultsDisabled() {
	bool b;
	asm volatile ("msr faultmask, %0" : "=r" (b));
	return b;
}

__attribute__ ((always_inline))
extern inline void cortexm_irq_SetBasePrio(cortexm_irq_Prio p) {
	asm volatile ("mrs %0, baseprio" :: "r" (p));
}

__attribute__ ((always_inline))
extern inline cortexm_irq_Prio cortexm_irq_BasePrio() {
	cortexm_irq_Prio p;
	asm volatile ("msr baseprio, %0" : "=r" (p));
	return p;
}
