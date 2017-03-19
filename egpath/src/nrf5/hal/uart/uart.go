package uart

import (
	"mmio"
	"unsafe"

	"nrf5/hal/internal"
	"nrf5/hal/te"
)

type Periph struct {
	te.Regs

	_        [31]mmio.U32
	errorsrc mmio.U32
	_        [7]mmio.U32
	enable   mmio.U32
	_        mmio.U32
	pcelrts  mmio.U32
	pseltxd  mmio.U32
	pselcts  mmio.U32
	pselrxd  mmio.U32
	rxd      mmio.U32
	txd      mmio.U32
	_        mmio.U32
	baudrate mmio.U32
	_        [17]mmio.U32
	config   mmio.U32
}

//emgo:const
var UART0 = (*Periph)(unsafe.Pointer(internal.BaseAPB + 0x02000))

type Task byte

const (
	STARTRX Task = 0  // Start UART receiver.
	STOPRX  Task = 1  // Stop UART receiver.
	STARTTX Task = 2  // Start UART transmitter.
	STOPTX  Task = 3  // Stop UART transmitter.
	SUSPEND Task = 19 // Suspend UART.
)

type Event byte

const (
	CTS    Task = 0  // CTS is activated (set low). Clear To Send.
	NCTS   Task = 1  // CTS is deactivated (set high). Not Clear To Send.
	RXDRDY Task = 2  // Data received in RXD.
	TXDRDY Task = 3  // Data sent from TXD.
	ERROR  Task = 11 // Error detected.
	RXTO   Task = 43 // Receiver timeout.
)
