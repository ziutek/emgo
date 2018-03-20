package clock

import (
	"nrf5/hal/te"
)

type Task byte

const (
	HFCLKSTART Task = 0 // Start high frequency crystal oscilator.
	HFCLKSTOP  Task = 1 // Stop high frequency crystal oscilator.
	LFCLKSTART Task = 2 // Start low frequency source.
	LFCLKSTOP  Task = 3 // Stop low frequency source.
	CAL        Task = 4 // Start calibration of low freq. RC oscilator.
	CTSTART    Task = 5 // Start calibration timer.
	CTSTOP     Task = 6 // Stop calibration timer.
)

func (t Task) Task() *te.Task { return r().Regs.Task(int(t)) }

type Event byte

const (
	HFCLKSTARTED Event = 0 // High frequency crystal oscilator started.
	LFCLKSTARTED Event = 1 // Low frequency source started.
	DONE         Event = 3 // Calibration of low freq. RC osc. complete.
	CTTO         Event = 4 // Calibration timer timeout.
)

func (e Event) Event() *te.Event { return r().Regs.Event(int(e)) }

// LoadHFCLKRUN returns true if HFCLKSTART task was triggered.
func LoadHFCLKRUN() bool {
	return r().hfclkrun.Load() != 0
}

type Source byte

const (
	RC    Source = 0
	XTAL  Source = 1
	SYNTH Source = 2
)

// LoadHFCLKStat returns information about HFCLK status (running or not) and
// clock source.
func LoadHFCLKSTAT() (src Source, running bool) {
	s := r().hfclkstat.Load()
	return Source(s & 1), s&(1<<16) != 0
}

// LoadLFCLKRUN returns true if LFCLKSTART task was triggered.
func LoadLFCLKRUN() bool {
	return r().lfclkrun.Bit(0) != 0
}

// LoadLFCLKSTAT returns information about LFCLK status (running or not) and
// clock source.
func LoadLFCLKSTAT() (src Source, running bool) {
	s := r().lfclkstat.Load()
	return Source(s & 1), s&(1<<16) != 0
}

// LoadLFCLKSRCCOPY returns clock source for LFCLK from time when LFCLKSTART
// task has been triggered.
func LoadLFCLKSRCCOPY() Source {
	return Source(r().lfclksrccopy.Bits(3))
}

// LoadLFCLKSRC returns clock source for LFCLK.
func LoadLFCLKSRC() Source {
	return Source(r().lfclksrc.Bits(3))
}

// StoreLFCLKSRC sets clock source for LFCLK. It can only be modified when
// LFCLK is not running.
func StoreLFCLKSRC(src Source) {
	r().lfclksrc.Store(uint32(src))
}

// LoadCTIV returns calibration timer interval in milliseconds.
func LoadCTIV() int {
	return int(r().ctiv.Bits(0x7f) * 250)
}

// StoreCTIV sets calibration timer interval as number of milliseconds
// (range: 250 ms to 31750 ms).
func StoreCTIV(ctiv int) {
	r().ctiv.Store(uint32(ctiv+125) / 250)
}

type XtalFreq byte

const (
	X16MHz XtalFreq = 0xff
	X32MHz XtalFreq = 0x00
)

// LoadXTALFREQ returns selected frequency of external crystal for HFCLK. nRF51.
func LoadXTALFREQ() XtalFreq {
	return XtalFreq(r().xtalfreq.Bits(0xff))
}

// StoreXTALFREQ selects frequency of external crystal for HFCLK. nRF51.
func StoreXTALFREQ(f XtalFreq) {
	r().xtalfreq.Store(uint32(f))
}

// TraceSpeed represents speed of Trace Port clock.
type TraceSpeed byte

const (
	T32MHz TraceSpeed = 0 // 32 MHz Trace Port clock (TRACECLK = 16 MHz).
	T16MHz TraceSpeed = 1 // 16 MHz Trace Port clock (TRACECLK = 8 MHz).
	T8MHz  TraceSpeed = 2 // 8 MHz Trace Port clock (TRACECLK = 4 MHz).
	T4MHz  TraceSpeed = 3 // 4 MHz Trace Port clock (TRACECLK = 2 MHz).
)

// TraceMux represents trace pins multiplexing configuration.
type TraceMux byte

const (
	GPIO     TraceMux = 0 // GPIOs multiplexed onto all trace pins.
	Serial   TraceMux = 1 // SWO onto P0.18, GPIO onto other trace pins.
	Parallel TraceMux = 2 // TRACECLK and TRACEDATA onto P0.20,18,16,15,14.
)

// LoadTRACECONFIG returns current speed of Trace Port clock and pin
// multiplexing of trace signals. nRF52.
func LoadTRACECONFIG() (TraceSpeed, TraceMux) {
	tc := r().traceconfig.Load()
	return TraceSpeed(tc & 3), TraceMux(tc >> 16 & 3)
}

// StoreTRACECONFIG sets speed of Trace Port clock and pin multiplexing of
// trace signals. nRF52.
func StoreTRACECONFIG(s TraceSpeed, m TraceMux) {
	r().traceconfig.Store(uint32(s) | uint32(m)<<16)
}
