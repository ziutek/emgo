#include "types/types.h"
#include "cortexm/startup.h"

extern uint32 dataStart, dataEnd, dataStartFlash;
extern uint32 bssStart, bssEnd;

int main_init();
int main_main();

void cortexm_startup_Start() {
	uint32 *src = &dataStartFlash;
	uint32 *dst = &dataStart;
	uint32 *end = &dataEnd;

	while (dst < end) {
		*dst = *src;
		++dst;
		++src;
	}

	dst = (uint32 *) & bssStart;
	end = (uint32 *) & bssEnd;

	while (dst < end) {
		*dst = 0;
		++dst;
	}
	
	main_init();
	main_main();
	
	for(;;);
}

extern uint32 mainStack;

uint32 *cortexm_startup_vectors[4] __attribute__ ((section(".vectors"))) = {
	&mainStack,
	(uint32 *) cortexm_startup_Start,          // entry point
	(uint32 *) cortexm_startup_DefaultHandler, // NMI
	(uint32 *) cortexm_startup_DefaultHandler, // hard fault
};
