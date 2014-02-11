#include "runtime/types.h"
#include "cortexm/startup.h"

extern uint32 __dataStart, __dataEnd, __dataStartFlash;
extern uint32 __bssStart, __bssEnd;

int main_init();
int main_main();

void cortexm_startup_Start() {
	uint32 *src = &__dataStartFlash;
	uint32 *dst = &__dataStart;
	uint32 *end = &__dataEnd;

	while (dst < end) {
		*dst = *src;
		++dst;
		++src;
	}

	dst = (uint32 *) & __bssStart;
	end = (uint32 *) & __bssEnd;

	while (dst < end) {
		*dst = 0;
		++dst;
	}
	
	main_init();
	main_main();
	
	for(;;);
}

extern uint32 stackptr;

uint32 *cortexm_startup_vectors[4] __attribute__ ((section(".vectors"))) = {
	&mainstack,
	(uint32 *) cortexm_startup_Start,          // entry point
	(uint32 *) cortexm_startup_DefaultHandler, // NMI
	(uint32 *) cortexm_startup_DefaultHandler, // hard fault
};
