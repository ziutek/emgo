package memory

// Type is a bitfield that describes limitations of some memory type.
type Type byte

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

// ObjectDescr specifies type and size for a dynamically allocated object.
// It also specifies location of pointers for garbage collector.
type ObjectDescr struct {
	Size uintptr    // object length in bytes
	Ptrs []uint     // bitfield for pointers 
	Type MemoryType // type of memory
}

// Allocator is type of function that aloocates continous block of memory that
// can hold n objects descibed by d.
type Allocator func(n int, d *ObjectDescr) uintptr

// SetAllocator sets allocator used by new() and make() to allocate memory.
func SetAllocator(a Allocator)
