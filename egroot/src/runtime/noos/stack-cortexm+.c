// +build cortexm0 cortexm3 cortexm4 cortexm4f

extern byte StacksBegin, ISRStack, MainStack, TaskStack, StacksEnd;

static inline uintptr
runtime$noos$stacksBegin() {
	return (uintptr) & StacksBegin;
}

static inline uintptr
runtime$noos$isrStackSize() {
	return (uintptr) & ISRStack;
}

static inline uintptr
runtime$noos$mainStackSize() {
	return (uintptr) & MainStack;
}

static inline uintptr
runtime$noos$taskStackSize() {
	return (uintptr) & TaskStack;
}

static inline uintptr
runtime$noos$stacksEnd() {
	return (uintptr) & StacksEnd;
}
