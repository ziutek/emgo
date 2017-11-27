package eve

type DCI interface {
	Read(s []byte)  // Read reads len(s) bytes into s.
	Write(s []byte) // Write writes len(s) bytes form s.
	End()           // End ends command.
	Err() error     // Err returns and resets internal error status.
}
