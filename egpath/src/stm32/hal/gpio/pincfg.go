package gpio

// Mode specifies operation mode.
type Mode byte

const (
	In    Mode = 0     // Logical input.
	Out   Mode = out   // Logical output.
	Alt   Mode = alt   // Alternate function output or bidirectional i/o.
	AltIn Mode = altIn // Alternate function input.
	Ana   Mode = ana   // Analog mode.
)

// Driver specifies output driver mode.
type Driver byte

const (
	PushPull  Driver = 0
	OpenDrain Driver = openDrain
)

// Speed specifies speed of output driver. Concrete STM32 series implements 3
// or 4 speed levels. Slow, Medium and Fast levels are always implemented. Actual
// speed at given level depends on supply voltage and load capacity. See the
// datasheet for more info.
type Speed int8

const (
	VeryLow  Speed = veryLow  // Typically < 1 MHz.
	Low      Speed = low      // Typically < 2 MHz.
	Medium   Speed = 0        // Typically < 25 MHz.
	High     Speed = high     // Typically < 50 MHz.
	VeryHigh Speed = veryHigh // Typically < 100 MHz.
)

// Pull specifies pull-up/pull-down configuration.
type Pull byte

const (
	NoPull   Pull = 0        // No pull-up/pull-down
	PullUp   Pull = pullUp   // Activate internal pull-up resistor.
	PullDown Pull = pullDown // Activate internal pull-down resistor.
)

// Config contains parameters used to setup GPIO pin.
type Config struct {
	Mode   Mode   // Mode: input, output, analog, alternate function.
	Driver Driver // Output driver type: push-pull or open-drain.
	Speed  Speed  // Output speed.
	Pull   Pull   // Pull-up/pull-down resistors.
}

// SetupPin configures n-th pin.
func (p Port) SetupPin(n int, cfg *Config) {
	setup(p, n, cfg)
}

// Setup configures pins.
func (p Port) Setup(pins Pins, cfg *Config) {
	for n := 0; n < 16; n++ {
		if pins&(1<<uint(n)) != 0 {
			setup(p, n, cfg)
		}
	}
}

// Lock locks configuration of n-th pin. Locked configuration can not be modified
// until reset.
func (p Port) Lock(pins Pins) {
	pins1 := pins | 0x10000
	p.lckr.Store(uint32(pins1))
	p.lckr.Store(uint32(pins))
	p.lckr.Store(uint32(pins1))
	p.lckr.Load()
	p.lckr.Load()
}
