__attribute__ ((always_inline)) extern inline 
bool sync$atomic$CompareAndSwapInt32(int32 *addr, int32 old, int32 new) {
	return __sync_bool_compare_and_swap(addr, old, new);
}

__attribute__ ((always_inline)) extern inline 
bool sync$atomic$CompareAndSwapInt64(int64 *addr, int64 old, int64 new) {
	return __sync_bool_compare_and_swap(addr, old, new);
}

__attribute__ ((always_inline)) extern inline 
bool sync$atomic$CompareAndSwapPointer(unsafe$Pointer *addr, unsafe$Pointer old, unsafe$Pointer new) {
	return __sync_bool_compare_and_swap(addr, old, new);
}

__attribute__ ((always_inline)) extern inline 
bool sync$atomic$CompareAndSwapUint32(uint32 *addr, uint32 old, uint32 new) {
	return __sync_bool_compare_and_swap(addr, old, new);
}

__attribute__ ((always_inline)) extern inline 
bool sync$atomic$CompareAndSwapUint64(uint64 *addr, uint64 old, uint64 new) {
	return __sync_bool_compare_and_swap(addr, old, new);
}


__attribute__ ((always_inline)) extern inline 
bool sync$atomic$CompareAndSwapUintptr(uintptr *addr, uintptr old, uintptr new) {
	return __sync_bool_compare_and_swap(addr, old, new);
}