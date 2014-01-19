#include <unistd.h>

#include "runtime/types.h"

void libc_stdio_Print(string s) {
	write(1, s.str, s.len);
}
