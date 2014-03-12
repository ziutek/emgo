// +build none

extern uint32 MaxTasks;

__attribute__ ((always_inline))
extern inline int runtime_MaxTasks() {
	return (int)&MaxTasks;
}
