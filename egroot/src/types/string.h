typedef struct {
	byte *str;
	uint len;
} string;

#define __EGSTR(s) ((string){(byte *)s, sizeof(s)-1})