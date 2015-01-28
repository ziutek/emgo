#define ithead builtin$ITHead
#define tinfo  builtin$Type

typedef struct {
	ithead h$;
	string (*Error)(ival *);
} error;
