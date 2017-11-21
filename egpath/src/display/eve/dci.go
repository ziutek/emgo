package eve

type DCI interface {
	Begin() // Begin begins transaction (in case of SPI it sets CSN pin to low).
	End()   // End ends transaction (in case of SPI is sets CSN pin to high).

	Read(s []byte)        // Read reads len(s) bytes into s.
	Write(s []byte)       // Write writes len(s) bytes form s.
	WriteString(s string) // WriteString writes len(s) bytes form s.

	Err() error // Err returns and resets internal error status.
}
