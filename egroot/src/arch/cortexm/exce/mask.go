package exce

// Disable disables all exceptions other than NMI and faults. Internally it
// sets Cortex-M PRIMASK to 1. Atomic primitives on Cortex-M0 always enable
// exceptions after atomic operation. If you need this functions on Cortex-M0
// don't use channels, mutexes, don't allocate memory and maybe don't do more
// things!
func DisablePri()

// Enable reverts Disable. If you modified any data that can be used by enabled
// interrupt handlers you probably need to call sync.Memory() before.
func EnablePri()

// Disabled returns true if excepions are disabled (PRIMASK != 0).
func DisabledPri() bool

// Disable disables all exceptions other than NMI. Internally it sets
// FAULTMASK to 1. Not supported by Cortex-M0.
func Disable()

// Enable reverts DisableFaults. If you modified any data that can be used
// by enabled interrupt handlers you probably need to call sync.Memory() before.
// Not supported by Cortex-M0.
func Enable()

// Disabled returns true if all excepions are disabled(FAULTMASK != 0).
// Not supported by Cortex-M0.
func Disabled() bool

// SetBasePrio sets BASEPRIO register. It prevents the activation of exceptions
// with the same or lower as p. Not supported by Cortex-M0.
func SetBasePrio(p Prio)

// BasePrio returns current value of BASEPRIO register.
func BasePrio() Prio