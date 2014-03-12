package runtime

// Panic is called by builtin panic function. If Panic returns panic executes
// infinite loop.
var Panic func(s string) = defaultPanic

func defaultPanic(s string) {
	// do nothing
}
