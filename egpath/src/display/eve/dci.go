package eve

type DCI interface {
	Read(s []byte)        // Read reads len(s) bytes into s.
	Write(s []byte)       // Write writes len(s) bytes from s.
	End()                 // End finishes current read/write transaction.
	Err(clear bool) error // Err returns and clears internal error status.
	IRQ() <-chan struct{} // IRQ allows to wait for IRQ by reading from channel.
	SetPDN(pdn int)       // SetPDN sets PD_N pin up/down.
}
