package eve

type DCI interface {
	Read(s []byte)        // Read reads len(s) bytes into s.
	Write(s []byte)       // Write writes len(s) bytes from s.
	End()                 // End ends command.
	Err(clear bool) error // Err returns and clears the internal error status.
}
