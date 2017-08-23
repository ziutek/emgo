package main

import (
	"bluetooth/att"
	"bluetooth/gatt"
	"bluetooth/uuid"

	"fmt"
)

//emgo:const
var (
	srvNordicUART   = uuid.UUID{0x6E400001B5A3F393, 0xE0A9E50E24DCCA9E}
	chrNordicUARTTx = uuid.UUID{0x6E400002B5A3F393, 0xE0A9E50E24DCCA9E}
	chrNordicUARTRx = uuid.UUID{0x6E400003B5A3F393, 0xE0A9E50E24DCCA9E}
)

type chr struct {
	Handle    uint16
	Prop      byte
	ValHandle uint16
	UUID      uuid.UUID
}

type service struct {
	Handle         uint16
	GroupEndHandle uint16
	UUID           uuid.UUID
	Chrs           []chr
}

type gattServer struct {
	services []service
}

func (gs *gattServer) ServeATT(w *att.ResponseWriter, r *att.Request) {
	fmt.Printf("\r\nMethod/Cmd: %v/%v\r\n", r.Method, r.Cmd)
	fmt.Printf("Handle-End: %04x-%04x\r\n", r.Handle, r.EndHandle)
	fmt.Printf("Other:      %04x\r\n", r.Other)
	fmt.Printf("UUID:       %v\r\n\r\n", r.UUID)

	switch r.Method {
	case att.Read:
		gs.read(w, r)
	case att.ReadByType:
		gs.readByType(w, r)
	case att.ReadByGroupType:
		gs.readByGroupType(w, r)
	case att.FindInformation:
		gs.findInformation(w, r)
	default:
		w.SetError(att.RequestNotSupported, r)
	}
	if err := w.Send(); err != nil {
		fmt.Printf("Can't send response: %v\r\n", err)
	}
}

func (gs *gattServer) readByType(w *att.ResponseWriter, r *att.Request) {
	fmt.Printf("readByType\r\n")
	if r.UUID.CanShorten(32) {
		switch r.UUID.Short32() {
		case 0x2803: // GATT Characteristic Declaration.
			fieldSize := 0
		loop:
			for i := range gs.services {
				srv := &gs.services[i]
				for k := range srv.Chrs {
					chr := &srv.Chrs[k]
					if chr.ValHandle == 0 {
						continue
					}
					if chr.Handle >= r.Handle && chr.Handle <= r.EndHandle {
						short := chr.UUID.CanShorten(16)
						if fieldSize == 0 {
							if short {
								fieldSize = 5 + 2
							} else {
								fieldSize = 5 + 16
							}
							w.SetReadByType(fieldSize)
						} else if fieldSize == 5+2 && !short ||
							fieldSize == 5+16 && short {
							break loop
						}
						w.AppendWord16(chr.Handle)
						w.AppendByte(chr.Prop)
						w.AppendWord16(chr.ValHandle)
						if short {
							w.AppendUUID16(chr.UUID.Short16())
						} else {
							w.AppendUUID(chr.UUID)
						}
						if !w.Commit() {
							break loop // MTU reached.
						}
						fmt.Printf(
							"Chr: %x %v %x %v\r\n",
							chr.Handle, chr.Prop, chr.ValHandle, chr.UUID,
						)
					}
				}
			}
			if fieldSize != 0 {
				return
			}
		}
	}
	w.SetError(att.AttributeNotFound, r)
}

func (gs *gattServer) readByGroupType(w *att.ResponseWriter, r *att.Request) {
	fmt.Printf("readByGroupType\r\n")
	if r.UUID.CanShorten(32) {
		switch r.UUID.Short32() {
		case 0x2800: // GATT Primary Service Declaration
			fieldSize := 0
		loop:
			for i := range gs.services {
				srv := &gs.services[i]
				if srv.Handle >= r.Handle && srv.Handle <= r.EndHandle {
					short := srv.UUID.CanShorten(16)
					if fieldSize == 0 {
						if short {
							fieldSize = 4 + 2
						} else {
							fieldSize = 4 + 16
						}
						w.SetReadByGroupType(fieldSize)
					} else if fieldSize == 4+2 && !short ||
						fieldSize == 4+16 && short {
						break loop
					}
					w.AppendWord16(srv.Handle)
					w.AppendWord16(srv.GroupEndHandle)
					if short {
						w.AppendUUID16(srv.UUID.Short16())
					} else {
						w.AppendUUID(srv.UUID)
					}
					if !w.Commit() {
						break loop // MTU reached.
					}
					fmt.Printf(
						"Srv: %x %x %v\r\n",
						srv.Handle, srv.GroupEndHandle, srv.UUID,
					)
				}
			}
			if fieldSize != 0 {
				return
			}
		}
	}
	w.SetError(att.AttributeNotFound, r)
}

func (gs *gattServer) findInformation(w *att.ResponseWriter, r *att.Request) {
	fmt.Printf("findInformation\r\n")
	var format att.FindInformationFormat
loop:
	for i := range gs.services {
		srv := &gs.services[i]
		// BUG: Information about services not returned.
		for k := range srv.Chrs {
			chr := &srv.Chrs[k]
			if chr.Handle >= r.Handle && chr.Handle <= r.EndHandle {
				short := chr.UUID.CanShorten(16)
				if format == 0 {
					if short {
						format = att.HandleUUID16
					} else {
						format = att.HandleUUID
					}
					w.SetFindInformation(format)
				} else if format == att.HandleUUID16 && !short ||
					format == att.HandleUUID && short {
					break loop
				}
				w.AppendWord16(chr.Handle)
				if short {
					w.AppendUUID16(chr.UUID.Short16())
				} else {
					w.AppendUUID(chr.UUID)
				}
				if !w.Commit() {
					break loop // MTU reached.
				}
				fmt.Printf(
					"Chr: %x %v %x %v\r\n",
					chr.Handle, chr.Prop, chr.ValHandle, chr.UUID,
				)
			}
		}
	}
	if format != 0 {
		return
	}
	w.SetError(att.AttributeNotFound, r)
}

func  (gs *gattServer) read(w *att.ResponseWriter, r *att.Request) {
	w.SetRead()
	w.AppendString("Hello world!")
	w.Commit()
}

const (
	srvGenericAccessProfile    uuid.UUID16 = 0x1800
	srvGenericAttributeProfile uuid.UUID16 = 0x1801
)

var gattSrv = &gattServer{
	services: []service{
		{
			0x0001, 0x0007, srvGenericAccessProfile.Full(),
			[]chr{
				{
					0x0002, gatt.Write | gatt.Read, 0x0003,
					gatt.DeviceName.Full(),
				}, {
					0x0004, gatt.Read, 0x0005,
					gatt.Apperance.Full(),
				}, {
					0x0006, gatt.Read, 0x0007,
					gatt.PeriphPrefConnParams.Full(),
				},
			},
		},
		{
			0x0008, 0x000B, srvGenericAttributeProfile.Full(),
			[]chr{
				{
					0x0009, gatt.Indicate, 0x000A,
					gatt.ServiceChanged.Full(),
				}, {
					0x000B, 0, 0,
					gatt.ClientChrConfig.Full(),
				},
			},
		},
		{
			0x000C, 0xFFFF, srvNordicUART,
			[]chr{
				{
					0x000D, gatt.Notify, 0x000E,
					chrNordicUARTRx,
				}, {
					0x000F, 0, 0,
					gatt.ClientChrConfig.Full(),
				}, {
					0x0010, gatt.Write | gatt.WriteWithoutResp, 0x0011,
					chrNordicUARTTx,
				},
			},
		},
	},
}
