// Package sync provides synchronisation primitives.
package barrier

// Compiler is a full compiler memory barrier. All compiler optimizations
// that can reorder explicit memory accesses can't cross this barrier.
func Compiler()

// Memory is a hardware full memory barrier. All explicit memory accesses
// before Memory should be finished before first subsequent explicit memory
// access. Memory implies Compiler.
func Memory()

// Sync is similar to Memory but stronger: it ensures that all explicit memory
// accesses before Sync should be finished before execution of first subsequent
// instruction. Sync implies Compiler.
func Sync()