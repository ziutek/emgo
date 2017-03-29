package gpio

type Mode byte

const (
	In     Mode = 0
	InOut  Mode = 1
	Discon Mode = 2
	Out    Mode = 3
)

type Pull byte

const (
	NoPull   Pull = 0
	PullDown Pull = 1
	PullUp   Pull = 3
)

type Drive byte

const (
	S0S1 Drive = 0 // Standard 0, standard 1.
	H0S1 Drive = 1 // High drive 0, standard 1.
	S0H1 Drive = 2 // Standard 0, high drive 1.
	H0H1 Drive = 3 // High drive 0, high drive 1.
	D0S1 Drive = 4 // Disconnect 0, standard 1.
	D0H1 Drive = 5 // Disconnect 0, high drive 1.
	S0D1 Drive = 6 // Standard 0, disconnect 1.
	H0D1 Drive = 7 // High drive 0, disconnect 1.
)

type Sense byte

const (
	NoSense   Sense = 0
	SenseHigh Sense = 2
	SenseLow  Sense = 3
)

type Config struct {
	Mode  Mode
	Pull  Pull
	Drive Drive
	Sense Sense
}

// Setup configures n-th pin in port p.
func (p *Port) SetupPin(n int, cfg Config) {
	p.pincnf[n].Store(
		uint32(cfg.Sense)<<16 | uint32(cfg.Drive)<<8 | uint32(cfg.Pull)<<2 |
			uint32(cfg.Mode),
	)
}

// PinConfig returns current configuration of n-th pin in port p.
func (p *Port) PinConfig(n int) Config {
	c := p.pincnf[n].Load()
	return Config{
		Mode:  Mode(c & 3),
		Pull:  Pull(c >> 2 & 3),
		Drive: Drive(c >> 8 & 7),
		Sense: Sense(c >> 16 & 3),
	}
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
