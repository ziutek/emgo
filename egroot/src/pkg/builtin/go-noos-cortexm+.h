// +build noos
// +build cortexm0 cortexm3 cortexm4 cortexm4f

#define go(call) do {						\
	void func() {							\
		call;								\
		asm volatile ("svc 1");				\
	}										\
	register void (*r0)() asm("r0") = func;	\
	asm volatile ("svc 0" :: "r" (r0));		\
} while(0)

/*
Jesli wolana funkcja f wymaga parametrow to call() musi być wywołaniem funkcji
typu wrapper. np:

void f(int i, byte b);

void call(int _0, byte _1) {
	__schedOn();  // asm volatile ("svc 3"); odblokowanie schedulera
	f(_0, _1);
}

Jesli funkcja nie wymaga parametrow to call moze bys sama funkcja
ale nalezy nie blokowac schedulera.

go powinno miec dodatkowy parametr wybierajacy svc 0/1

*/	
