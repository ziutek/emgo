package main

type Menu struct {
	disp Display
	g    [2]Gauge
	sel  byte
	cnt  byte
}

func (m *Menu) Display() *Display {
	return &m.disp
}

const (
	SetInRPM = iota
	SetOutRPM
	ShowRPM
	DispOff
	Init
)

func (m *Menu) SetMaxRPM(rpm int) {
	for i := range m.g {
		m.g[i].SetMax(rpm)
	}
}

func (m *Menu) clearDisp() {
	if m.sel == DispOff {
		for i := 0; i < 8; i++ {
			m.disp.Clear(i)
		}
	} else {
		m.disp.WriteString(0, 0, 8, "")
	}
}

func (m *Menu) Select(item byte) {
	m.sel = item
	m.clearDisp()
}

func (m *Menu) Next() {
	m.sel = (m.sel + 1) & 3
	m.clearDisp()
}

func (m *Menu) printDec(row, val int) {
	addr := 4 * row
	m.disp.WriteDec(addr, addr+3, 4, val)
}

func (m *Menu) printStr(row int, s string) {
	addr := 4 * row
	m.disp.WriteString(addr, addr, 4, s)
}

func (m *Menu) clearRow(row int) {
	addr := 4 * row
	m.disp.WriteString(addr, addr, 4, "")
}

func (m *Menu) RTCISR() {
	if m.disp.RTCISR() != 7 {
		return // Wait for the frame to be completed.
	}
	switch {
	case m.sel <= SetOutRPM:
		set := int(m.sel & 1)
		if m.cnt += 8; m.cnt > 80 {
			if rpm := fc.TargetRPM(set); rpm < 0 {
				m.printStr(set, "IErr") // Identification error.
			} else {
				m.printDec(set, rpm)
			}
		} else {
			m.clearRow(set)
		}
		show := set ^ 1
		m.printDec(show, fc.RPM(show))
	case m.sel == ShowRPM:
		m.printDec(0, fc.RPM(0))
		m.printDec(1, fc.RPM(1))
	case m.sel == Init:
		if m.cnt += 8; m.cnt >= 80 {
			m.printStr(0, "Idnt") // Identification.
			m.printDec(1, fc.IdentProgress())
		} else {
			m.clearRow(0)
		}
	}
}

func (m *Menu) HandleEncoder(change int) (n, rpm int) {
	if m.sel > SetOutRPM {
		return -1, 0
	}
	g := &m.g[m.sel]
	g.AddCube(change)
	return int(m.sel), g.Val()
}
