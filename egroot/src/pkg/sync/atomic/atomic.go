// Package atomic provides low-level atomic memory primitives.
//
// These functions do only what specified and don't provide any implicit memory
// barriers nor "happens-before" edges for different gorutines.
package atomic

import "unsafe"

func CompareAndSwapInt32(addr *int32, old, new int32) (swapped bool)
func CompareAndSwapInt64(addr *int64, old, new int64) (swapped bool)
func CompareAndSwapPointer(addr *unsafe.Pointer, old, new unsafe.Pointer) (swapped bool)
func CompareAndSwapUint32(addr *uint32, old, new uint32) (swapped bool)
func CompareAndSwapUint64(addr *uint64, old, new uint64) (swapped bool)
func CompareAndSwapUintptr(addr *uintptr, old, new uintptr) (swapped bool)

func AddInt32(addr *int32, delta int32) (new int32)
func AddUint32(addr *uint32, delta uint32) (new uint32)
func AddInt64(addr *int64, delta int64) (new int64)
func AddUint64(addr *uint64, delta uint64) (new uint64)
func AddUintptr(addr *uintptr, delta uintptr) (new uintptr)

func OrInt32(addr *int32, mask int32) (new int32)
func OrUint32(addr *uint32, mask uint32) (new uint32)
func OrInt64(addr *int64, mask int64) (new int64)
func OrUint64(addr *uint64, mask uint64) (new uint64)
func OrUintptr(addr *uintptr, mask uintptr) (new uintptr)

func AndInt32(addr *int32, mask int32) (new int32)
func AndUint32(addr *uint32, mask uint32) (new uint32)
func AndInt64(addr *int64, mask int64) (new int64)
func AndUint64(addr *uint64, mask uint64) (new uint64)
func AndUintptr(addr *uintptr, mask uintptr) (new uintptr)

func SwapInt32(addr *int32, new int32) (old int32)
func SwapInt64(addr *int64, new int64) (old int64)
func SwapPointer(addr *unsafe.Pointer, new unsafe.Pointer) (old unsafe.Pointer)
func SwapUint32(addr *uint32, new uint32) (old uint32)
func SwapUint64(addr *uint64, new uint64) (old uint64)
func SwapUintptr(addr *uintptr, new uintptr) (old uintptr)

func LoadInt32(addr *int32) (val int32)
func LoadInt64(addr *int64) (val int64)
func LoadUint32(addr *uint32) (val uint32)
func LoadUint64(addr *uint64) (val uint64)
func LoadUintptr(addr *uintptr) (val uintptr)
func LoadPointer(addr *unsafe.Pointer) (val unsafe.Pointer)

func StoreInt32(addr *int32, val int32)
func StoreInt64(addr *int64, val int64)
func StoreUint32(addr *uint32, val uint32)
func StoreUint64(addr *uint64, val uint64)
func StoreUintptr(addr *uintptr, val uintptr)
func StorePointer(addr *unsafe.Pointer, val unsafe.Pointer)
