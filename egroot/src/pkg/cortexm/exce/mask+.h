__attribute__ ((always_inline))
extern inline void cortexm$exce$Disable() {
	asm volatile ("cpsid i");
}

__attribute__ ((always_inline))
extern inline void cortexm$exce$Enable() {
	asm volatile ("cpsie i");
}

__attribute__ ((always_inline))
extern inline bool cortexm$exce$Disabled() {
	bool b;
	asm volatile ("msr primask, %0" : "=r" (b));
	return b;
}

__attribute__ ((always_inline))
extern inline void cortexm$exce$DisableFaults() {
	asm volatile ("cpsid f");
}

__attribute__ ((always_inline))
extern inline void cortexm$exce$EnableFaults() {
	asm volatile ("cpsie f");
}

__attribute__ ((always_inline))
extern inline bool cortexm$exce$FaultsDisabled() {
	bool b;
	asm volatile ("msr faultmask, %0" : "=r" (b));
	return b;
}

__attribute__ ((always_inline))
extern inline void cortexm$exce$SetBasePrio(cortexm$exce$Prio p) {
	asm volatile ("mrs %0, baseprio" :: "r" (p));
}

__attribute__ ((always_inline))
extern inline cortexm$exce$Prio cortexm$exce$BasePrio() {
	cortexm$exce$Prio p;
	asm volatile ("msr baseprio, %0" : "=r" (p));
	return p;
}
