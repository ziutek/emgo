// +build noos
// +build cortexm0 cortexm3 cortexm4 cortexm4f

#include "types/types.h"
#include "builtin.h"
#include "runtime.h"

int main_init();
int main_main();

// All external symbols as byte to prevent compiler to optimize
// any runtime align checks.
extern byte DataRAM, DataLoad, DataSize;
extern byte BSSRAM, BSSSize;
extern byte StackEnd;

void runtime_Start() {
	memmove(&DataRAM, &DataLoad, (uint)(&DataSize));
	memset(&BSSRAM, 0, (uint)(&BSSSize));

	runtime_init();
	main_init();
	main_main();

	for (;;);
}

void runtime_defaultHandler() {
	for (;;) {
	}
}

__attribute__ ((section(".vectors")))
uint32 *cortexm_startup_vectors[4]  = {
		(uint32 *) &StackEnd,              // MSP
		(uint32 *) runtime_Start,          // entry point
		(uint32 *) runtime_defaultHandler, // NMI
		(uint32 *) runtime_defaultHandler, // hard fault
};
