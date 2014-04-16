#include <unistd.h>

#include "runtime/types.h"

void libc$stdio$Print(string s) {
	write(1, s.str, s.len);
}
