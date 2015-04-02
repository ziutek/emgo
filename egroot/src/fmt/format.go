package fmt

import "io"

// State is passed to the custom formatters. It provides io.Writer interface and
// information about format options like width, precision and flags.
type State interface {
	io.Writer
	Width() (width int, ok bool)
	Precision() (prec int, ok bool)
	Flag(c int) bool
}

// Stringer can be implemented by value to provide its string representation.
type Stringer interface {
    String() string
}

// Formatter can be implemented by value that knows how to format itself. 
type Formatter interface {
	Format(f State, c rune)
}