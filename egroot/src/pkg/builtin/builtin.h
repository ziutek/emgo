void memmove(unsafe_Pointer dst, unsafe_Pointer src, uint n);
void memcpy(unsafe_Pointer dst, unsafe_Pointer src, uint n);
void memset(unsafe_Pointer s, byte b, uint n);

__attribute__ ((noreturn))
void panic(string s);