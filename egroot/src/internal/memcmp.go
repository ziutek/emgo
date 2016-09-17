package internal

import "unsafe"

// Memcmp compares two byte strings str1 and str2 each of length n, pointed by
// pointers p1 and p2 respectively. It returns value:
//	- less than zero if str1 < str2,
//	- equal zero if str1 == str2,
//  - greater than zero is string1 > string2.
func Memcmp(p1, p2 unsafe.Pointer, n uintptr) int