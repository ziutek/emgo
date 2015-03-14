__attribute__ ((always_inline))
extern inline
reflect$Value reflect$valueOf(interface i) {
	union {interface in; reflect$Value out;} cast = {};
	cast.in = i;
	return cast.out;
}