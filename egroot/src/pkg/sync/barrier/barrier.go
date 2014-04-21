// Package barrier provides memory barriers.
// They are intended for use to implement higher level synchronisation
// primitives or perform interaction with hardware that need defined order of
// memory accesses.
package barrier

// Compiler is a full compiler memory barrier. I does nothing but optimizations
// treat it as function that can modify any variable.
func Compiler()

// Memory is a hardware full memory barrier. All explicit memory accesses
// performed by CPU before Memory should be finished before first subsequent
// explicit memory access. Memory implies Compiler.
func Memory()

// Sync is similar to Memory but stronger: it ensures that all explicit memory
// accesses performed by CPU before Sync should be finished before execution
// of the first subsequent instruction. Sync implies Compiler.
func Sync()