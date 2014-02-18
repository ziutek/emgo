#include "types/types.h"
#include "runtime.h"

#include "cortexm/startup.h"

extern uint32 _DataRAM, _DataLoad, _DataLen;
extern uint32 _BSSRAM, _BSSLen;

int main_init();
int main_main();

void cortexm_startup_Start() {
	runtime_Copy(&_DataRAM, &_DataLoad, (uint)(&_DataLen));
	runtime_Memset(&_BSSRAM, 0, (uint)(&_BSSLen));
	
	main_init();
	main_main();
	
	for(;;);
}

extern uint32 _MainStack;

uint32 *cortexm_startup_vectors[4] __attribute__ ((section(".vectors"))) = {
	&_MainStack,
	(uint32 *) cortexm_startup_Start,          // entry point
	(uint32 *) cortexm_startup_DefaultHandler, // NMI
	(uint32 *) cortexm_startup_DefaultHandler, // hard fault
};
