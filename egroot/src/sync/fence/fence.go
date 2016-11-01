package fence

// Compiler is the full compiler memory fence. It does nothing, but optimizer
// treats it as function that can modify any variable. Do not use this function
// for  synchronisation. Use it only if you want to avoid optimisation at some
// point in code.
//
//c:static inline
func Compiler()

////c:static inline
//func Memory()


// Sync ensures that any instruction after it, in program order, do not execute
// until all explicit memory accesses before it complete.
//
//c:static inline
func Sync()

// Memory ensures that any memory access after it, in program order, do not
// execute until all explicit memory accesses before it complete.
//
//c:static inline
func Memory() 

// SMP works like Memory in multiprocessor system and like Compiler in
// uniprocessor system.
//
//c:static inline
func SMP() 