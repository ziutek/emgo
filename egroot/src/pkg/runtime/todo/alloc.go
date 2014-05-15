package memory

// Type is a bitfield that describes limitations of some memory type.
type Type int

const (
	Universal Type = 0
	NoDMA     Type = 1
)

// Typer is an interface that can be implemented by any dynamically allocated
// object. It specifies type of memory that satisfies object's needs. In case
// of new() and make(), whether type implements Typer is determined at
// compilation time, so there is no any additional runtime overhead for
// "ordinary" allocation and small (function call) overhead for allocation of
// memory of specific type.
type Typer interface {
	RuntimeMemoryType() Type
}