typedef __INT8_TYPE__ int8;
typedef __INT16_TYPE__ int16;
typedef __INT32_TYPE__ int32;
typedef __INT64_TYPE__ int64;

typedef __UINT8_TYPE__ uint8;
typedef __UINT16_TYPE__ uint16;
typedef __UINT32_TYPE__ uint32;
typedef __UINT64_TYPE__ uint64;

typedef __UINTPTR_TYPE__ uintptr;
typedef void *unsafe$Pointer;

typedef unsigned long uint;
typedef signed long int_;

typedef float float32;
typedef double float64;

typedef float _Complex complex64;
typedef double _Complex complex128;

#ifndef EG_ALLOW_C_TYPES

// Forget about C types
#define int			XXint / / /
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
typedef uint8 byte;
typedef int32 rune;

#define true (1)
#define false (0)

typedef struct {
	byte _;
} structE;

#define complex(re, im) ((re)+1.i*(im))
#define real(c) (__real__(c))
#define imag(c) (__imag__(c))

#define nil (0)

#define len(v) ((v).len)
#define cap(v) ((v).cap)
