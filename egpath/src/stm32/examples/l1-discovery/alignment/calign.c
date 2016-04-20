#include "__noos_cortexm3_l1xx_md.h"

uintptr
main$SizeByte() {
	return sizeof(byte);
}

uintptr
main$SizeInt() {
	return sizeof(int_);
}

uintptr
main$SizeInt16() {
	return sizeof(int16);
}

uintptr
main$SizeInt32() {
	return sizeof(int32);
}

uintptr
main$SizeInt64() {
	return sizeof(int64);
}

uintptr
main$SizeS16() {
	return sizeof(main$S16);
}

uintptr
main$SizeS32() {
	return sizeof(main$S32);
}

uintptr
main$SizeS64() {
	return sizeof(main$S64);
}

uintptr
main$AlignByte() {
	return __alignof__(byte);
}

uintptr
main$AlignInt() {
	return __alignof__(int_);
}

uintptr
main$AlignInt16() {
	return __alignof__(int16);
}

uintptr
main$AlignInt32() {
	return __alignof__(int32);
}

uintptr
main$AlignInt64() {
	return __alignof__(int64);
}

uintptr
main$AlignS16() {
	return __alignof__(main$S16);
}

uintptr
main$AlignS32() {
	return __alignof__(main$S32);
}

uintptr
main$AlignS64() {
	return __alignof__(main$S64);
}
