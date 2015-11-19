// Package fence provides memory fences.
// They are intended for use to implement higher level synchronisation
// primitives or perform interaction with hardware that need defined order of
// memory accesses.
package fence

// Compiler is a full compiler memory fence. I does nothing but optimizations
// treat it as function that can modify any variable.
//c:static inline
func Compiler()

// Memory is a hardware full memory fence. All explicit memory accesses
// performed by CPU before Memory should be finished before first subsequent
// explicit memory access. Memory implies Compiler.
//c:static inline
func Memory()

// Sync is similar to Memory but stronger: it ensures that all explicit memory
// accesses performed by CPU before Sync should be finished before execution
// of the first subsequent instruction. Sync implies Compiler.
//c:static inline
func Sync()
