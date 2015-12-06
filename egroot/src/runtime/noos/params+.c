extern byte HeapBegin, HeapEnd, MaxTasks;

uintptr
runtime$noos$heapBegin() {
	return (uintptr) (&HeapBegin);
}

uintptr
runtime$noos$heapEnd() {
	return (uintptr) (&HeapEnd);
}

int
runtime$noos$maxTasks() {
	int mt = (int)&MaxTasks;

	// 0 is valid value of &MaxTasks. GCC during optimization assumes that
	// &MaxTask can't be 0 and removes any code that should be executed only
	// when &MaxTask is 0. Following line tells GCC that mt can be any value.
	asm volatile ("":"=r" (mt):"0"(mt));

	return mt;
}
