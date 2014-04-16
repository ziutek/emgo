package atomic

import "unsafe"

func CompareAndSwapInt32(addr *int32, old, new int32) (swapped bool)

func CompareAndSwapInt64(addr *int64, old, new int64) (swapped bool)

func CompareAndSwapPointer(addr *unsafe.Pointer, old, new unsafe.Pointer) (swapped bool)

func CompareAndSwapUint32(addr *uint32, old, new uint32) (swapped bool)

func CompareAndSwapUint64(addr *uint64, old, new uint64) (swapped bool)

func CompareAndSwapUintptr(addr *uintptr, old, new uintptr) (swapped bool)
