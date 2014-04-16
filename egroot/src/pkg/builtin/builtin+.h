void memmove(unsafe$Pointer dst, unsafe$Pointer src, uint n);
void memcpy(unsafe$Pointer dst, unsafe$Pointer src, uint n);
void memset(unsafe$Pointer s, byte b, uint n);

__attribute__ ((noreturn))
void panic(string s);