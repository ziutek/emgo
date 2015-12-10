package sync

// Fence is a full compiler memory fence. I does nothing, but optimizations
// treat it as function that can modify any variable. Do not use this function
// for synchronisation. Use it only to avoid optimisation at some place.
//c:static inline
func Fence()