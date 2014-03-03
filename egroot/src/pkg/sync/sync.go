// Package sync provides synchronisation primitives.
package sync

// Barrier is a full compiler memory barrier. All compiler optimizations
// that can reorder explicit memory accesses (at source code level) can't
// cross this barrier.
func Barrier()

// Memory is a hardware full memory barrier. All explicit memory accesses
// (at instruction level) before Memory should be finished before first
// subsequent explicit memory access. Memory implies Barrier.
func Memory()

// Sync is similar to Memory but stronger: it ensures that all explicit memory
// accesses before Sync should be finished before first subsequent instruction.
// Sync implies Barrier.
func Sync()