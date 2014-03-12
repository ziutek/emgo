// +build cortexm0 cortexm3 cortexm4 cortexm4f

#include "types/types.h"
#include "runtime.h"

int main_init();
int main_main();

void runtime_Start() {
	runtime_Copy(&DataRAM, &DataLoad, (uint)(&DataSize));
	runtime_Memset(&BSSRAM, 0, (uint)(&BSSSize));

	runtime_init();
	main_init();
	main_main();

	for (;;);
}

extern uint32 _MainStack;
void runtime_defaultHandler() {
	for (;;) {
	}
}

uint32 *cortexm_startup_vectors[4] __attribute__ ((section(".vectors"))) = {
		(uint32 *) &FreeEnd,               // MSP
		(uint32 *) runtime_Start,          // entry point
		(uint32 *) runtime_defaultHandler, // NMI
		(uint32 *) runtime_defaultHandler, // hard fault
};
