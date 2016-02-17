package hdc

func (d *Display) ClearDisplay() error {
	return writeCmd(d, 0x01)
}

func (d *Display) ReturnHome() error {
	return writeCmd(d, 0x02)
}

type EntryMode byte

const (
	ShiftOff EntryMode = 0
	ShiftOn  EntryMode = 1 << 0
	Decr     EntryMode = 0
	Incr     EntryMode = 1 << 1
)

func (d *Display) SetEntryMode(f EntryMode) error {
	return writeCmd(d, byte(0x04|f&0x03))
}

type DisplayMode byte

const (
	DisplayOff DisplayMode = 0
	DisplayOn  DisplayMode = 1 << 2
	CursorOff  DisplayMode = 0
	CursorOn   DisplayMode = 1 << 1
	BlinkOff   DisplayMode = 0
	BlinkOn    DisplayMode = 1
)

func (d *Display) SetDisplayMode(f DisplayMode) error {
	return writeCmd(d, byte(0x08|f&7))
}

type Shift byte

const (
	ShiftCuror  Shift = 0
	ShiftScreen Shift = 1 << 3
	ShiftLeft   Shift = 0
	ShiftRight  Shift = 1 << 2
)

func (d *Display) SetShift(f Shift) error {
	return writeCmd(d, byte(0x10|f&0xc))
}

type Function byte

const (
	OneLine  Function = 0
	TwoLines Function = 1 << 3
	Font5x8  Function = 0
	Font5x10 Function = 1 << 2
)

func (d *Display) SetFunction(f Function) error {
	return writeCmd(d, byte(0x20|f&0x0f))
}

func (d *Display) SetAddrCGRAM(addr int) error {
	return writeCmd(d, byte(0x40|addr&0x3f))
}

func (d *Display) SetAddrDDRAM(addr int) error {
	return writeCmd(d, byte(0x80|addr&0x7f))
}

// Flush calls bufio.Writer.Flush if bufio.Writer was used as io.Write
func (d *Display) MoveCursor(col, row int) error {
	var addr int
	switch row {
	case 0:
		addr = col
	case 1:
		addr = 0x40 + col
	case 2:
		addr = int(d.Cols) + col
	case 3:
		addr = 0x40 + int(d.Cols) + col
	}
	return d.SetAddrDDRAM(addr)
}

func (d *Display) Write(data []byte) (int, error) {
	for i, b := range data {
		if err := writeData(d, b); err != nil {
			return i, err
		}
	}
	return len(data), nil
}

func (d *Display) WriteString(s string) (int, error) {
	for i, n := 0, len(s); i < n; i++ {
		if err := writeData(d, s[i]); err != nil {
			return i, err
		}
	}
	return len(s), nil
}
