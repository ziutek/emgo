typedef builtin$Chan chan;

#define MAKECHAN(typ, cap) builtin$MakeChan(cap, sizeof(typ), __alignof__(typ))

#define SEND(cx, typ, val) do {                     \
	chan c = cx;                                    \
	typeof(typ) v = val;                            \
	unsafe$Pointer$$uintptr r = c.M->Send(c.P, &v); \
	if (r._0 != nil) {                              \
		*(typ*)r._0 = v;                            \
		c.M->Done(c.P, r._1);                       \
	}                                               \
} while(0)

#define RECV(typ, cx) ({                            \
	chan c = cx;                                    \
	typeof(typ) v = {0};                            \
	unsafe$Pointer$$uintptr r = c.M->Recv(c.P, &v); \
	if (r._0 != nil) {                              \
		v = *(typ*)r._0;                            \
		c.M->Done(c.P, r._1);                       \
	}                                               \
	v;                                              \
})

#define RECVOK(tt, cx) ({                                \
	chan c = cx;                                         \
	tt vok = {0};                                        \
	unsafe$Pointer$$uintptr r = c.M->Recv(c.P, &vok._0); \
	if (r._0 != nil) {                                   \
		vok._0 = *(typeof(&vok._0))r._0;                 \
		c.M->Done(c.P, r._1);                            \
		vok._1 = true;                                   \
	} else if (r._1 == 0) {                              \
		vok._1 = true;                                   \
	}                                                    \
	vok;                                                 \
})
	