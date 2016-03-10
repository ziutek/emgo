### Emgo follows Go specification with exception for memory allocation.

Think about Emgo as C with: Go syntax, Go packages, Go building philosophy.

In Go, variable declared in function can be allocated on the stack or on the heap - escaping analisis is used for decision. In Emgo (like in C), all local variables are stack allocated. Dynamic allocation and garbage collection can occurs only when:

1. new or make builtin function is used.
2. Non-empty strings are concatanated. 
3. Builtin append function is called and there is not enough space in destination.
4. An element to map is added or removed.

Current allocator is trivial and doesn't contain GC. This isn't big disadventage for many embeded applications that need allocation only during startup. They often run on MCU that has only few kilobytes of SRAM and in many cases they must respond in realtime. The simple "stop the world GC" can't be used and much sophisticated one consumes to much Flash/SRAM.

The target is to allow chose between multiple allocators (simply import it as package) and use one that best fits application needs.

Examples:

The following function is correct Go function:

func F() ([]byte, *int) {
	i := 4
	return []byte{1, 2, 3}, &i
}

but it isn't correct Emgo function - you need to rewrite it this way:

func F() ([]byte, *int) {
	i := new(int)
	*i = 4
	b := append([]byte{}, []byte{1, 2, 3}...)
	return b, i
}

### Unexported methods

By default Emgo does not include information about unexported methods in typeinfo. Use minfo pragma to disable this "feature"..

### Generated C code.

Current gotoc generates C code that relies on many GCC extensions.

### Standard library.

Emgo standard library doesn't follow Go standard library. When porting some Go package to Emgo there is always prefered to preserve original interface and package name.

### Not yet implemented:

Maps.
Defer.
String concatanation.
Append.
Unnamed structs.
Closures.