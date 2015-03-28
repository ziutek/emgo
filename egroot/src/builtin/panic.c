#include <internal/types.h>
#include <builtin.h>

void panic(interface i) {
	string *s = (string*)(&i.val$); // Usefull if i contains string.
	for (;;);
}

void panicIC() {
	panic(INTERFACE(EGSTR("interface conversion"), &string$$));
}