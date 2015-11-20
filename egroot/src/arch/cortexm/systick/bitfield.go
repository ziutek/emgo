package systick

import "mmio"

type Field uint16
type Bit byte

type U8 struct{ r mmio.U8 }

func (u *U8) Bit(n Bit) bool            { return u.r.Bit(int(n)) }
func (u *U8) SetBit(n Bit)              { u.r.SetBit(int(n)) }
func (u *U8) ClearBit(n Bit)            { u.r.ClearBit(int(n)) }
func (u *U8) Field(f Field) uint8       { return u.r.Field(uint16(f)) }
func (u *U8) SetField(f Field, v uint8) { u.r.SetField(uint16(f), v) }

type U16 struct{ r mmio.U16 }

func (u *U16) Bit(n Bit) bool             { return u.r.Bit(int(n)) }
func (u *U16) SetBit(n Bit)               { u.r.SetBit(int(n)) }
func (u *U16) ClearBit(n Bit)             { u.r.ClearBit(int(n)) }
func (u *U16) Field(f Field) uint16       { return u.r.Field(uint16(f)) }
func (u *U16) SetField(f Field, v uint16) { u.r.SetField(uint16(f), v) }

type U32 struct{ r mmio.U32 }

func (u *U32) Bit(n Bit) bool             { return u.r.Bit(int(n)) }
func (u *U32) SetBit(n Bit)               { u.r.SetBit(int(n)) }
func (u *U32) ClearBit(n Bit)             { u.r.ClearBit(int(n)) }
func (u *U32) Field(f Field) uint32       { return u.r.Field(uint16(f)) }
func (u *U32) SetField(f Field, v uint32) { u.r.SetField(uint16(f), v) }

const o = 8
