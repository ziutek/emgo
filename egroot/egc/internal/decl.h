struct tinfo;
typedef struct tinfo tinfo;

struct minfo;
typedef struct minfo minfo;

__attribute__ ((noreturn))
void panic(interface i);

__attribute__ ((noreturn))
void panicIndex();
