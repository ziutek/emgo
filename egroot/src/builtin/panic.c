#include <internal/types.h>
#include <builtin.h>

void panic(interface i) {
	for (;;) {
		if (builtin$Panic != nil) {
			builtin$Panic(i);
		}
	}
}

void panicIC() {
	panic(INTERFACE(EGSTR("interface conversion"), &string$$));
}

void panicIndex() {
	panic(INTERFACE(EGSTR("index out of range"), &string$$));
}
