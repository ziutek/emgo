__attribute__ ((always_inline))
extern inline void cortexm$irq$Disable() {
	asm volatile ("cpsid i");
}

__attribute__ ((always_inline))
extern inline void cortexm$irq$Enable() {
	asm volatile ("cpsie i");
}

__attribute__ ((always_inline))
extern inline bool cortexm$irq$Disabled() {
	bool b;
	asm volatile ("msr primask, %0" : "=r" (b));
	return b;
}

__attribute__ ((always_inline))
extern inline void cortexm$irq$DisableFaults() {
	asm volatile ("cpsid f");
}

__attribute__ ((always_inline))
extern inline void cortexm$irq$EnableFaults() {
	asm volatile ("cpsie f");
}

__attribute__ ((always_inline))
extern inline bool cortexm$irq$FaultsDisabled() {
	bool b;
	asm volatile ("msr faultmask, %0" : "=r" (b));
	return b;
}

__attribute__ ((always_inline))
extern inline void cortexm$irq$SetBasePrio(cortexm$irq$Prio p) {
	asm volatile ("mrs %0, baseprio" :: "r" (p));
}

__attribute__ ((always_inline))
extern inline cortexm$irq$Prio cortexm$irq$BasePrio() {
	cortexm$irq$Prio p;
	asm volatile ("msr baseprio, %0" : "=r" (p));
	return p;
}
