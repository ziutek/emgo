package gatt

import (
	"bluetooth/uuid"
)

const (
	DeviceName           uuid.UUID16 = 0x2A00
	Apperance            uuid.UUID16 = 0x2A01
	PeriphPrefConnParams uuid.UUID16 = 0x2A04
	ServiceChanged       uuid.UUID16 = 0x2A05
	ClientChrConfig      uuid.UUID16 = 0x2902
)
