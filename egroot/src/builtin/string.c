#include "types.h"

bool _stringEq(string s1, string s2) {
	if (s1.len != s2.len) {
		return false;
	}
	int i = s1.len;
	while (i != 0) {
		--i;
		if (s1.str[i] != s2.str[i]) {
			return false;
		}
	}
	return true;
}