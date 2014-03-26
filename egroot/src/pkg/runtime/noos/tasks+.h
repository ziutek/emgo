extern byte MaxTasks;

__attribute__ ((always_inline))
extern inline int runtime_noos_MaxTasks() {
	return (int)&MaxTasks;
}