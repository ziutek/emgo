package fence

// Compiler is the full compiler memory fence. It does nothing, but optimizer
// treats it as function that can modify any variable. Do not use this function
// for  synchronisation. Use it only if you want to avoid optimisation at some
// point in code.
//
//c:static inline
func Compiler()

// RW  ensures that any memory access (normal or I/O) after it, in program
// order, do not execute until all explicit memory accesses before it complete.
//
//c:static inline
func RW()

// RW_SMP works like Compiler in uniprocessor system. In multiprocessor system
// RW_SMP ensures that any normal memory access after it, in program order, do
// not execute until all explicit normal memory accesses before it complete.
//
//c:static inline
func RW_SMP()

//c:static inline
func R()

//c:static inline
func R_SMP()

//c:static inline
func W()

//c:static inline
func W_SMP()

//c:static inline
func RDP()

//c:static inline
func RDP_SMP()
