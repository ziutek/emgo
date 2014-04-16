// +build cortexm0 cortexm3 cortexm4 cortexm4f

.syntax unified

.global __clear_cache

// __clear_cache is need by gcc.
// void __clear_cache(char *begin, char *end)

.thumb_func
__clear_cache:
	dsb  
	isb  
	bx   lr
