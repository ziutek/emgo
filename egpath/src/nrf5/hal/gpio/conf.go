package gpio

type Config uint32

const (
	ModeIn     Config = 0 // Input buffer connected, output disabled.
	ModeInOut  Config = 1 // Input buffer connected, output enabled.
	ModeDiscon Config = 2 // Input buffer disconnected, output disabled.
	ModeOut    Config = 3 // Input buffer disconnected, output enabled.

	PullNone Config = 0 << 2 // No pull.
	PullDown Config = 1 << 2 // Pull down on pin.
	PullUp   Config = 3 << 2 // Pull dup on pin.

	DriveS0S1 Config = 0 << 8 // Standard 0, standard 1.
	DriveH0S1 Config = 1 << 8 // High drive 0, standard 1.
	DriveS0H1 Config = 2 << 8 // Standard 0, high drive 1.
	DriveH0H1 Config = 3 << 8 // High drive 0, high drive 1.
	DriveD0S1 Config = 4 << 8 // Disconnect 0, standard 1.
	DriveD0H1 Config = 5 << 8 // Disconnect 0, high drive 1.
	DriveS0D1 Config = 6 << 8 // Standard 0, disconnect 1.
	DriveH0D1 Config = 7 << 8 // High drive 0, disconnect 1.

	SenseNone Config = 0 << 16 // Sense disabled.
	SenseHigh Config = 2 << 16 // Sense for high level.
	SenseLow  Config = 3 << 16 // Sense for low level.
)

// Setup configures n-th pin in port p.
func (p *Port) SetupPin(n int, cfg Config) {
	p.pincnf[n].Store(uint32(cfg))
}

// PinConfig returns current configuration of n-th pin in port p.
func (p *Port) PinConfig(n int) Config {
	return Config(p.pincnf[n].Load())
}

// Setup configures pins.
func (p *Port) Setup(pins Pins, cfg Config) {
	for n := 0; n < 32; n++ {
		if pins&(1<<uint(n)) != 0 {
			p.SetupPin(n, cfg)
		}
	}
}

// SetDirIn allows a fast change of direction to input for specified pins.
func (p *Port) SetDirIn(pins Pins) {
	p.dirset.Store(uint32(pins))
}

// SetDirOut allows afast change of direction to output for specified pins.
func (p *Port) SetDirOut(pins Pins) {
	p.dirclr.Store(uint32(pins))
}
