extern byte MaxTasks;

__attribute__ ((always_inline))
extern inline
int runtime$noos$MaxTasks() {
	int mt = (int)&MaxTasks;
	
	// 0 is valid value of &MaxTasks. GCC during optimization assumes that
	// &MaxTask can't be 0 and removes any code that should be executed only
	// when &MaxTask is 0. Following line tells GCC that mt can be any value.
	asm volatile ("" : "=r" (mt): "0" (mt));
	
	return mt;
}