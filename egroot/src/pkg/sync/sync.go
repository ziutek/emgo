// Package sync provides synchronisation primitives.
package sync

// Barrier is a full compiler memory barrier. All compiler optimizations
// that can reorder explicit memory accesses (at source code level) can't
// cross this barrier.
func Barrier()

// Memory is a hardware full memory barier. All optimizations performed by
// CPU that can reorder explicit memory accesses (at instruction level)
// can't cross this barrier. Memory implies Barrier.
func Memory()
