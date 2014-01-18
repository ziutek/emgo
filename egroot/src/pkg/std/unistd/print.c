#include <unistd.h>

#include "runtime/types.h"

void unistd_Print(string s) {
	write(1, s.str, s.len);
}