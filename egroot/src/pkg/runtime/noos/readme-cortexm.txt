### Configuration options for noos/cortexm target.

#### StackExp, StackFrac (no default value)

Defines the size of the stack for one task as follows: 

	(1 << StackExp) * StackFrac / 8.

Such a little complicated definition allows to easy use MPU (if available) for stack protection.


#### MaxTasks

Defines the maximum number of tasks.

MaxTasks == 0 means:

- no tasks (gorutines),
- one stack (uses MSP),
- runtime doesn't touch any peripherals,
- program runs in privileged mode.

MaxTasks > 0 means:

- separate stack for any task + separate stack for ISRs,
- MPU (if avilable) is used to protect the stacks and other areas,
- all tasks runs in user mode.

#### IRTExp

Defines the length of the table for ISR vectors. IRTExp should be >= 5 (2^5 allows to handle 16 system exceptions and 16 external intnerrupts.) IRTExp is used only if MaxTasks > 0.

#### Remarks

All configuration options should be set at the beginning of the linker script.

They are visible for C code as external symbols. In runtime code they are always declared as `byte` to prevent compiler to optimize any runtime align checks.