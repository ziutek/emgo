package eve

type DCI interface {
	Read(s []byte)        // Read reads len(s) bytes into s.
	Write32(s []uint32)   // Write writes len(s) words from s.
	End()                 // End ends command.
	Err(clear bool) error // Err returns and clears internal error status.
}
