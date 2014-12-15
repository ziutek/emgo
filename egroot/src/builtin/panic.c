#include "internal/types.h"

void panic(interface i) {
	string *s = (string*)(&i.val$); // Usefull if i contains string.
	for (;;);
}