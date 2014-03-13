package noos

import "unsafe"

const minStack = 128

var stacks struct{ first, last, minCap uintptr }

// stackeSetup setups len(tasks)+1 stacks in mem, down, from end of mem to
// beginning. First (number zero) stack is located at end of mem. It usually
// has greater capacity than other stacks and is intended to used by initial
// task. Last stack (number len(task)) is intended to be used by exception
// handlers.
func stacksSetup(mem []byte) {
	begin := uintptr(unsafe.Pointer(&mem[0]))
	end := begin + uintptr(len(mem))
	begin = alignUp(begin, 4)
	end = alignDown(end, 4)

	n := uintptr(len(tasks) + 1)

	if begin+n*minStack > end {
		panicMemory()
	}
	stacks.minCap = (end - begin) / n
	stacks.minCap = alignDown(stacks.minCap, 4)
	if stacks.minCap < minStack {
		panicMemory()
	}
	stacks.first = end
	stacks.last = begin + stacks.minCap
}

// stackInitSP returns initial stack pointer for n-th task.
func stackInitSP(n int) uintptr {
	if n == 0 {
		return stacks.first
	}
	return stacks.last + uintptr(len(tasks)-n)*stacks.minCap
}

// stackCap returns stack capacity for n-th task.
func stackCap(n int) uintptr {
	if n == 0 {
		return stacks.first - stacks.last - uintptr(len(tasks))*stacks.minCap
	}
	return stacks.minCap
}