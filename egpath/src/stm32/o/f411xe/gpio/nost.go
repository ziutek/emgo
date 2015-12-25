package gpio

const (
	MODER_IN  = 0 // Input (reset state).
	MODER_OUT = 1 // General purpose output mode.
	MODER_ALT = 2 // Alternate function mode.
	MODER_ANA = 3 // Analog mode.
)

const (
	OT_PP = 0 // Output push-pull (reset state).
	OT_OD = 1 // Output open-drain.
)

const (
	PUPDR_FLOAT = 0 // No pull-up, no pull-down.
	PUPDR_PUP   = 1 // Pull-up.
	PUPDR_PDOWN = 2 // Pull-down.
)
