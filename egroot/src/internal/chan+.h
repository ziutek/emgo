typedef internal$Chan chan;

#define NILCHAN (chan){}

#define MAKECHAN(typ, cap) internal$MakeChan(cap, sizeof(typ), __alignof__(typ))

#define SEND(cx, typ, val) do {                     \
	chan c = cx;                                    \
	typeof(typ) v = val;                            \
	unsafe$Pointer$$uintptr r = c.M->Send(c.C, &v); \
	if (r._0 != nil) {                              \
		*(typ*)r._0 = v;                            \
		c.M->Done(c.C, r._1);                       \
	}                                               \
} while(0)

#define RECV(typ, cx, zero) ({                      \
	chan c = cx;                                    \
	typeof(typ) v = zero;                           \
	unsafe$Pointer$$uintptr r = c.M->Recv(c.C, &v); \
	if (r._0 != nil) {                              \
		v = *(typ*)r._0;                            \
		c.M->Done(c.C, r._1);                       \
	}                                               \
	v;                                              \
})

#define RECVOK(tt, cx) ({                                \
	chan c = cx;                                         \
	tt vok = {};                                         \
	unsafe$Pointer$$uintptr r = c.M->Recv(c.C, &vok._0); \
	if (r._0 != nil) {                                   \
		vok._0 = *(typeof(&vok._0))r._0;                 \
		c.M->Done(c.C, r._1);                            \
		vok._1 = true;                                   \
	} else {                                             \
		vok._1 = (r._1 == internal$ChanOK);               \
	}                                                    \
	vok;                                                 \
})
 
#define SENDINIT(i, c, typ, v) \
	chan chan##i = c;          \
	typeof(typ) val##i = v
	
#define RECVINIT(i, c, typ) \
	chan chan##i = c;       \
	typeof(typ) val##i 
 
#define SENDCOMM(i) {               \
	.Case   = &&case##i,            \
	.C      = chan##i.C,            \
	.E      = &val##i,              \
	.Try    = chan##i.M->TrySend,   \
	.Cancel = chan##i.M->CancelSend \
}

#define RECVCOMM(i) {               \
	.Case   = &&case##i,            \
	.C      = chan##i.C,            \
	.E      = &val##i,              \
	.Try    = chan##i.M->TryRecv,   \
	.Cancel = chan##i.M->CancelRecv \
}

#define _SELECT(dflt, commList...)                              \
	internal$Comm arr[] = {commList};                            \
	int_ n = sizeof(arr)/sizeof(arr[0]);                         \
	internal$Comm *comms[n];                                     \
	int_ i = n;                                                  \
	while (i--) {                                               \
		comms[i] = &arr[i];                                     \
	}                                                           \
	unsafe$Pointer$$unsafe$Pointer$$uintptr r = internal$Select( \
		(slice){comms, n, n}, dflt                              \
	);                                                          \
	goto *r._0

#define SELECT(commList...) _SELECT(nil, commList)

#define NBSELECT(commList...) _SELECT(&&dflt, commList)

#define SELSEND(i) do {                	  \
	if (r._1 != nil) {                    \
		*(typeof(&val##i))r._1 = val##i;  \
		chan##i.M->Done(chan##i.C, r._2); \
	}                                     \
} while(0)

#define SELRECV(i) ({                     \
	if (r._1 != nil) {                    \
		val##i = *(typeof(&val##i))r._1;  \
		chan##i.M->Done(chan##i.C, r._2); \
	}                                     \
	val##i;                               \
})

#define SELRECVOK(i) ({                        \
	if (r._1 != nil) {                         \
		val##i._0 = *(typeof(&val##i._0))r._1; \
		chan##i.M->Done(chan##i.C, r._2);      \
		val##i._1 = true;                      \
	} else {                                   \
		val##i._1 = (r._2 == internal$ChanOK);  \
	}                                          \
	val##i;                                    \
})

static inline
void close(chan c) {
	c.M->Close(c.C);
}

static inline
int_ clen(chan c) {
	return c.M->Len(c.C);
}

static inline
int_ ccap(chan c) {
	return c.M->Cap(c.C);
}
