
typedef __INT8_TYPE__ int8;
typedef __INT16_TYPE__ int16;
typedef __INT32_TYPE__ int32;
typedef __INT64_TYPE__ int64;

typedef __UINT8_TYPE__ byte;
typedef __UINT8_TYPE__ uint8;
typedef __UINT16_TYPE__ uint16;
typedef __UINT32_TYPE__ uint32;
typedef __UINT64_TYPE__ uint64;
typedef unsigned int uint;

typedef __UINTPTR_TYPE__ uintptr;

typedef void* unsafe_Pointer;

typedef float float32;
typedef double float64;

typedef float _Complex complex64;
typedef double _Complex complex128;

// Forget about C types
#ifndef EG_ALLOW_C_TYPES

#define	unsigned	XXunsigned / / /
#define	signed		XXsigned / / /
#define	char		XXchar / / /
#define	short		XXshort / / /
#define	long		XXlong / / /
#define	float		XXfloat / / /
#define	double		XXdouble / / /
#define _Complex	XX_Complex / / /

#endif

typedef uint8 bool;

#define true (1)
#define false (0)

typedef struct {
	byte *str;
	int len;
} string;

#define _GOSTR(s) (string) {(byte *)s, sizeof(s)-1}

bool _string_eq(string s1, string s2);
