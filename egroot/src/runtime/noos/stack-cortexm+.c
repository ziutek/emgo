// +build cortexm0 cortexm3 cortexm4 cortexm4f cortexm7f cortexm7d

extern byte StacksBegin, ISRStack, MainStack, TaskStack, StacksEnd;

inline __attribute__((always_inline))
uintptr
runtime$noos$stacksBegin() {
	return (uintptr) & StacksBegin;
}

inline __attribute__((always_inline))
uintptr
runtime$noos$isrStackSize() {
	return (uintptr) & ISRStack;
}

inline __attribute__((always_inline))
uintptr
runtime$noos$mainStackSize() {
	return (uintptr) & MainStack;
}

inline __attribute__((always_inline))
uintptr
runtime$noos$taskStackSize() {
	return (uintptr) & TaskStack;
}

inline __attribute__((always_inline))
uintptr
runtime$noos$stacksEnd() {
	return (uintptr) & StacksEnd;
}
