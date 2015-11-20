package systick

type SysTick struct {
	SYST_CSR   U32
	SYST_RVR   U32
	SYST_CVR   U32
	SYST_CALIB U32
}

var R = (*SysTick)(unsafe.Pointer(uintptr(0xe000e010)))

// SYST_CSR
const (
	ENABLE    Bit = 0
	TICKINT   Bit = 1
	CLKSOURCE Bit = 2
	COUNTFLAG Bit = 16
)

// SYST_RVR
const (
	RELOAD Field = 24<<o + 0
)

// SYST_CVR
const (
	CURRENT Field = 24<<o + 0
)

// SYST_CALIB
const (
	TENMS Field = 24<<o + 0
	SKEW  Bit   = 30
	NOREF Bit   = 31
)
