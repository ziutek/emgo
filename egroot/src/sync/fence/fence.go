package fence

// Compiler is a full compiler memory fence. It does nothing, but optimizer
// treats it as function that can modify any variable. Do not use this function
// for  synchronisation. Use it only if you want to avoid optimisation at some
// point in code.
//
//c:static inline
func Compiler()

// Memory issues a full memory barrier. 
//
//c:static inline
func Memory()