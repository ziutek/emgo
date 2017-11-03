// +build l476xx

#include <internal/types.h>
#include <internal.h>

#include <stm32/o/l476xx/mmap.h>
// type decl
// var  decl
// func decl
// const decl
// type def
// var  def
// func def
// init
void stm32$o$l476xx$mmap$init() {
	static bool called = false;
	if (called) {
		return;
	}
	called = true;
	internal$init();
}
