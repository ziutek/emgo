__attribute__ ((always_inline))
extern inline
void arch$cortexm$exce$Disable() {
	asm volatile ("cpsid i");
}

__attribute__ ((always_inline))
extern inline
void arch$cortexm$exce$Enable() {
	asm volatile ("cpsie i");
}

__attribute__ ((always_inline))
extern inline
bool arch$cortexm$exce$Disabled() {
	bool b;
	asm volatile ("msr primask, %0" : "=r" (b));
	return b;
}

__attribute__ ((always_inline))
extern inline
void arch$cortexm$exce$DisableFaults() {
	asm volatile ("cpsid f");
}

__attribute__ ((always_inline))
extern inline
void arch$cortexm$exce$EnableFaults() {
	asm volatile ("cpsie f");
}

__attribute__ ((always_inline))
extern inline
bool arch$cortexm$exce$FaultsDisabled() {
	bool b;
	asm volatile ("msr faultmask, %0" : "=r" (b));
	return b;
}

__attribute__ ((always_inline))
extern inline
void arch$cortexm$exce$SetBasePrio(arch$cortexm$exce$Prio p) {
	asm volatile ("mrs %0, baseprio" :: "r" (p));
}

__attribute__ ((always_inline))
extern inline
arch$cortexm$exce$Prio arch$cortexm$exce$BasePrio() {
	arch$cortexm$exce$Prio p;
	asm volatile ("msr baseprio, %0" : "=r" (p));
	return p;
}
