// +build cortexm0 cortexm3 cortexm4 cortexm4f cortexm7f cortexm7d

extern byte VectorsBegin, VectorsEnd;

inline __attribute__((always_inline))
uintptr
runtime$noos$vectorsSize() {
	return ((uintptr) & VectorsEnd) - ((uintptr) & VectorsBegin);
}
