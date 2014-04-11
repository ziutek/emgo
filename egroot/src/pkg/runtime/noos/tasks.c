#include <types/types.h>

extern byte MaxTasks;

// 0 is valid value of &MaxTasks. GCC during optimization assumes that &MaxTask
// can't be 0 and removes any code that should be executed when &MaxTask is 0.
// There are some other methods to avoid this problem but place MaxTasks in
// separate object file seems to be the safest solution.
int runtime_noos_MaxTasks() {
	return (int)&MaxTasks;
}