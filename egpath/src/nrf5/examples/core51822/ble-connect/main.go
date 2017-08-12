package main

import (
	"fmt"
	"rtos"

	"bluetooth/att"
	"bluetooth/ble"
	"bluetooth/l2cap"
	"bluetooth/uuid"

	"nrf5/blec"
	"nrf5/hal/clock"
	"nrf5/hal/gpio"
	"nrf5/hal/irq"
	"nrf5/hal/rtc"
	"nrf5/hal/system"
	"nrf5/hal/system/timer/rtcst"
	"nrf5/hal/uart"
)

var (
	leds [5]gpio.Pin
	bctr *blec.Ctrl
	udrv *uart.Driver
)

func init() {
	system.Setup(clock.XTAL, clock.XTAL, true)
	rtcst.Setup(rtc.RTC1, 0)

	p0 := gpio.P0

	for i := range leds {
		led := p0.Pin(18 + i)
		led.Setup(gpio.ModeOut)
		leds[i] = led
	}

	bctr = blec.NewCtrl(ble.MaxDataPay, 3, 3)
	bctr.InitPhy()
	bctr.LEDs = &leds

	udrv = uart.NewDriver(uart.UART0, make([]byte, 80))
	udrv.P.StorePSEL(uart.SignalRXD, p0.Pin(11))
	udrv.P.StorePSEL(uart.SignalTXD, p0.Pin(9))
	udrv.P.StoreBAUDRATE(uart.Baud115200)
	udrv.P.StoreENABLE(true)
	udrv.EnableTx()
	rtos.IRQ(udrv.P.IRQ()).Enable()
	fmt.DefaultWriter = udrv
}

type pduLogger struct {
	*blec.Ctrl
}

func (p pduLogger) Recv() (ble.DataPDU, error) {
	pdu, err := p.Ctrl.Recv()
	i := p.Ctrl.Iter
	if err == nil {
		fmt.Printf(
			"R LL#%d LLID=%x P=%02x\r\n",
			i, pdu.Header()&ble.LLID, pdu.Payload(),
		)
	}
	return pdu, err
}

func main() {
	fmt.Printf("\r\nDevAddr: %08x\r\n", uint64(bctr.DevAddr()))

	u, err := uuid.Parse([]byte("abcdef01-1234-1651-8765-43805F9B34F1"))
	ub := make([]byte, 16)
	u.Encode(ub)
	u = uuid.DecodeLong(ub)

	fmt.Printf("%x %x %v %v\r\n", u.H, u.L, u, err)

	pdu := ble.MakeAdvPDU(ble.MaxAdvPay)
	pdu.SetType(ble.ScanRsp)
	pdu.SetTxAdd(bctr.DevAddr() < 0)
	pdu.AppendAddr(bctr.DevAddr())
	pdu.AppendString(ble.LocalName, "Emgo BLE")
	pdu.AppendUUIDs(ble.Services, serviceNordicUART)
	pdu.AppendBytes(ble.TxPower, 0)
	bctr.Advertise(pdu, 625)

	far := l2cap.NewBLEFAR(pduLogger{bctr})
	srv := att.NewServer(23)
	handler := new(nordicUART)
	srv.SetHandler(handler)
	for {
		cid, err := far.ReadHeader()
		if err != nil {
			fmt.Printf("ReadHeader: %v\r\n", err)
			continue
		}
		fmt.Printf("R L2CAP header: len=%d cid=%d\r\n", far.Len(), cid)
		switch cid {
		case 4: // ATT
			srv.HandleTransaction(far, cid)
		}
	}
}

/*
var (
	Nordic_UART = []byte{
		0x9E, 0xCA, 0xDC, 0x24, 0x0E, 0xE5,
		0xA9, 0xE0,
		0x93, 0xF3,
		0xA3, 0xB5,
		0x01, 0x00, 0x40, 0x6E,
	}
	Nordic_UART_Rx = []byte{
		0x9E, 0xCA, 0xDC, 0x24, 0x0E, 0xE5,
		0xA9, 0xE0,
		0x93, 0xF3,
		0xA3, 0xB5,
		0x03, 0x00, 0x40, 0x6E,
	}
	Nordic_UART_Tx = []byte{
		0x9E, 0xCA, 0xDC, 0x24, 0x0E, 0xE5,
		0xA9, 0xE0,
		0x93, 0xF3,
		0xA3, 0xB5,
		0x02, 0x00, 0x40, 0x6E,
	}
)

func decodeU16(b []byte) uint16 {
	return uint16(int(b[0]) | int(b[1])<<8)
}

func printBadLength(length int) {
	fmt.Printf(" error: bad length %d\r\n", length)
}

func writeError(far *l2cap.LEFAR, m att.Method, h uint16, e att.ErrorCode) {
	resp := []byte{
		byte(att.ErrorRsp), byte(m), byte(h), byte(h) >> 8, byte(e),
	}
	far.WriteHeader(len(resp), 4)
	far.Write(resp)
}

func writeResp_Decl(far *l2cap.LEFAR, hDecl uint16, prop gatt.Prop, hVal uint16, uuid []byte) {
	resp := []byte{
		byte(att.ReadByTypeRsp), byte(5 + len(uuid)),
		byte(hDecl), byte(hDecl) >> 8, byte(prop), byte(hVal), byte(hVal) >> 8,
	}
	far.WriteHeader(len(resp)+len(uuid), 4)
	far.Write(resp)
	far.Write(uuid)
}

func parseATT(far *l2cap.LEFAR, req []byte) {
	opcode := req[0]
	method := att.Method(opcode & 0x3F)
	cmd := opcode&0x40 != 0
	authSig := opcode&0x80 != 0

	fmt.Printf("R ATT: method=%02x cmd=%t authSig=%t\r\n", method, cmd, authSig)

	switch method {
	case att.ReqdByTypeReq:
		fmt.Printf("Read by Type Request\r\n")
		if len(req) != 7 && len(req) != 21 {
			writeError(far, method, 0, att.UnlikelyError)
			break
		}
		startHandle := decodeU16(req[1:3])
		endHandle := decodeU16(req[3:5])
		groupType := decodeU16(req[5:7])
		fmt.Printf(
			" startHandle=%04x endHandle=%04x groupType=%04x\r\n",
			startHandle, endHandle, groupType,
		)
		if groupType == 0x2802 {
			// GATT Include Declaration
			writeError(far, method, startHandle, att.AttributeNotFound)
			break
		}
		if groupType == 0x2803 {
			// GATT Characteristic Declaration.
			if startHandle <= 0x0D {
				writeResp_Decl(
					far, 0x0D,
					gatt.Notify, 0x0E, Nordic_UART_Rx,
				)
			} else if startHandle <= 0x10 {
				writeResp_Decl(
					far, 0x10,
					gatt.Write|gatt.WriteWithoutResp, 0x11, Nordic_UART_Tx,
				)
			} else if startHandle <= 0x11 {
				writeError(far, method, startHandle, att.AttributeNotFound)
			}
		}
	case att.ReadByGroupReq:
		fmt.Printf("Read by Group Type Request\r\n")
		if len(req) != 7 && len(req) != 21 {
			writeError(far, method, 0, att.UnlikelyError)
			break
		}
		startHandle := decodeU16(req[1:3]) // 0x0001
		endHandle := decodeU16(req[3:5])   // 0xFFFF
		groupType := decodeU16(req[5:7])
		fmt.Printf(
			" startHandle=%04x endHandle=%04x groupType=%04x\r\n",
			startHandle, endHandle, groupType,
		)
		resp := []byte{
			0x11,       // Opcode: Read by Group Type Response
			6,          // Length

			0x01, 0x00, // Attribute handle
			0x07, 0x00, // End group handle
			0x00, 0x18, // Attribute: Generic Access Profile

			0x08, 0x00, // Attribute handle
			0x0B, 0x00, // End group handle
			0x01, 0x18, // Attribute: Generic Attribute Profile
		}
		resp := []byte{
			0x11, // Opcode: Read by Group Type Response
			4 + byte(len(Nordic_UART)), // Length

			0x0C, 0x00, // Attribute handle
			0xFF, 0xFF, // End group handle
		}
		far.WriteHeader(len(resp)+len(Nordic_UART), 4)
		far.Write(resp)
		far.Write(Nordic_UART)
	}
}
*/

func radioISR() {
	bctr.RadioISR()
}

func uartISR() {
	udrv.ISR()
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTC1: rtcst.ISR,

	irq.RADIO: radioISR,

	irq.UART0: uartISR,
}
