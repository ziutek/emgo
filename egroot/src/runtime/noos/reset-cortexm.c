// +build cortexm0 cortexm3 cortexm4 cortexm4f

#include <internal/types.h>
#include <builtin.h>
#include <runtime/noos.h>

int runtime$init();
int main$init();
int main$main();

// All external symbols as byte to prevent compiler to optimize any runtime
// align checks.
extern byte DataRAM, DataLoad, DataSize;
extern byte BSSRAM, BSSSize;

void runtime$noos$Reset() {
	memmove(&DataRAM, &DataLoad, (uint) (&DataSize));
	memset(&BSSRAM, 0, (uint) (&BSSSize));

	runtime$init();
	main$init();
	main$main();

	for (;;);
}

/*__attribute__((section(".vectors")))
uint32 *cortexm$startup$vectors[4] = {
	(uint32 *) runtime$noos$Reset,        // entry point
	(uint32 *) runtime$noos$NMIHandler,   // NMI
	(uint32 *) runtime$noos$FaultHandler, // HardFault
};*/
