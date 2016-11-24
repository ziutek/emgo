### Configuration options for noos/cortexm target.

#### ISRStack, MainStack, TaskStack (no default value)

Defines the size of the stack for ISRs, main task and other tasks. It is rounding up to a multiple of 32 bytes.

#### MaxTasks

Defines the maximum number of tasks.

MaxTasks == 0 means:

- no tasks (no gorutines),
- one stack (uses MSP),
- runtime doesn't touch any peripherals,
- program runs in privileged mode.

MaxTasks > 0 means:

- separate stack for any task + separate stack for ISRs,
- MPU (if avilable) is used to implement stack guards,,
- all tasks runs in user mode.

#### Remarks

All configuration options should be set at the beginning of the linker script.

They are visible for C code as external symbols. In runtime code they are always declared as `byte` to prevent compiler to optimize any runtime align checks.