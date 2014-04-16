// +build noos
// +build cortexm0 cortexm3 cortexm4 cortexm4f

#include "builtin.h"

int main$init();
int main$main();

// All external symbols as byte to prevent compiler to optimize
// any runtime align checks.
extern byte DataRAM, DataLoad, DataSize;
extern byte BSSRAM, BSSSize;
extern byte StackEnd;

void runtime$Start() {
	memmove(&DataRAM, &DataLoad, (uint)(&DataSize));
	memset(&BSSRAM, 0, (uint)(&BSSSize));

	runtime$init();
	main$init();
	main$main();

	for (;;);
}

static
void defaultHandler() {
	for (;;) {
	}
}

__attribute__ ((section(".vectors")))
uint32 *cortexm$startup$vectors[4]  = {
		(uint32 *) &StackEnd,      // MSP
		(uint32 *) runtime$Start,  // entry point
		(uint32 *) defaultHandler, // NMI
		(uint32 *) defaultHandler, // hard fault
};
