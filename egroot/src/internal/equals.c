#include <internal/types.h>

bool
equals(string s1, string s2) {
	if (s1.len != s2.len) {
		return false;
	}
	int_ i;
	for (i = 0; i < s1.len; ++i) {
		if (s1.str[i] != s2.str[i]) {
			return false;
		}
	}
	return true;
}
