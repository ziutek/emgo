// This example shows how to use USART as serial console.
// It uses PA2, PA3 pins as Tx, Rx that are by default connected to ST-LINK
// Virtual Com Port.
package main

import (
	"bits"
	"fmt"
	"io"
	"reflect"
	"rtos"
	"strconv"
	"unsafe"

	"stm32/serial"

	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
	"stm32/hal/usart"
)

var con *serial.Dev

func init() {
	system.Setup96(8)
	systick.Setup()

	// GPIO

	gpio.A.EnableClock(true)
	port, tx, rx := gpio.A, gpio.Pin2, gpio.Pin3

	// USART

	port.Setup(tx, &gpio.Config{Mode: gpio.Alt})
	port.Setup(rx, &gpio.Config{Mode: gpio.AltIn, Pull: gpio.PullUp})
	port.SetAltFunc(tx|rx, gpio.USART2)

	s := usart.USART2

	s.EnableClock(true)
	s.SetBaudRate(115200)
	s.SetConf(usart.RxEna | usart.TxEna)
	s.EnableIRQ(usart.RxNotEmptyIRQ)
	s.Enable()

	con = serial.New(s, 80, 8)
	con.SetUnix(true)
	fmt.DefaultWriter = con

	rtos.IRQ(irq.USART2).Enable()
}

func conISR() {
	con.IRQ()
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.USART2: conISR,
}

type Bool bool

func (b Bool) Format(st fmt.State, c rune) {
	s := "fałsz"
	if b {
		s = "prawda"
	}
	io.WriteString(st, s)
}

func nle(n int, err error) {
	con.WriteByte('|')
	strconv.WriteInt(con, n, 10, 0)
	if err != nil {
		con.WriteString(" Err: ")
		con.WriteString(err.Error())
	}
	con.WriteByte('\n')
}

func main() {
	const (
		SmallestNormal         = 2.2250738585072014e-308
		SmallestNonzeroFloat64 = 4.940656458412465441765687928682213723651e-324
	)

	nle(strconv.WriteBool(con, true, 't', 0))
	nle(strconv.WriteBool(con, true, 't', -10))
	nle(strconv.WriteBool(con, true, 't', 10))
	nle(strconv.WriteBool(con, true, -'1', -10))
	nle(strconv.WriteBool(con, true, -'1', 10))

	nle(strconv.WriteUint32(con, 0, 10, 0))
	nle(strconv.WriteUint32(con, 1234567890, 10, 0))
	nle(strconv.WriteUint32(con, 1234567890, 10, -20))
	nle(strconv.WriteUint32(con, 1234567890, 10, 20))
	nle(strconv.WriteUint32(con, 1234567890, -10, -20))
	nle(strconv.WriteUint32(con, 1234567890, -10, 20))
	nle(strconv.WriteUint32(con, 0xf0f0f0f0, 2, 0))
	nle(strconv.WriteUint32(con, 0x12345678, 16, 0))
	nle(strconv.WriteUint32(con, 0x12345678, 16, -20))
	nle(strconv.WriteUint32(con, 0x12345678, 16, 20))
	nle(strconv.WriteUint32(con, 0x12345678, -16, -20))
	nle(strconv.WriteUint32(con, 0x12345678, -16, 20))

	nle(strconv.WriteInt32(con, 0, 10, 0))
	nle(strconv.WriteInt32(con, -1234567890, 10, 0))
	nle(strconv.WriteInt32(con, -1234567890, 10, -20))
	nle(strconv.WriteInt32(con, -1234567890, 10, 20))
	nle(strconv.WriteInt32(con, -1234567890, -10, -20))
	nle(strconv.WriteInt32(con, -1234567890, -10, 20))
	nle(strconv.WriteInt32(con, -0x10f0f0f0, 2, 0))
	nle(strconv.WriteInt32(con, -0x12345678, 16, 0))
	nle(strconv.WriteInt32(con, -0x12345678, 16, -20))
	nle(strconv.WriteInt32(con, -0x12345678, 16, 20))
	nle(strconv.WriteInt32(con, -0x12345678, -16, -20))
	nle(strconv.WriteInt32(con, -0x12345678, -16, 20))

	nle(strconv.WriteUint64(con, 0, 10, 0))
	nle(strconv.WriteUint64(con, 12345678900987654321, 10, 0))
	nle(strconv.WriteUint64(con, 12345678900987654321, 10, -20))
	nle(strconv.WriteUint64(con, 12345678900987654321, 10, 20))
	nle(strconv.WriteUint64(con, 12345678900987654321, -10, -20))
	nle(strconv.WriteUint64(con, 12345678900987654321, -10, 20))
	nle(strconv.WriteUint64(con, 0xf0f0f0f00f0f0f0f, 2, 0))
	nle(strconv.WriteUint64(con, 0x1234567887654321, 16, 0))
	nle(strconv.WriteUint64(con, 0x1234567887654321, 16, -20))
	nle(strconv.WriteUint64(con, 0x1234567887654321, 16, 20))
	nle(strconv.WriteUint64(con, 0x1234567887654321, -16, -20))
	nle(strconv.WriteUint64(con, 0x1234567887654321, -16, 20))

	nle(strconv.WriteInt64(con, 0, 10, 0))
	nle(strconv.WriteInt64(con, -1234567890987654321, 10, 0))
	nle(strconv.WriteInt64(con, -1234567890987654321, 10, -20))
	nle(strconv.WriteInt64(con, -1234567890987654321, 10, 20))
	nle(strconv.WriteInt64(con, -1234567890987654321, -10, -20))
	nle(strconv.WriteInt64(con, -1234567890987654321, -10, 20))
	nle(strconv.WriteInt64(con, -0x10f0f0f00f0f0f0f, 2, 0))
	nle(strconv.WriteInt64(con, -0x1234567887654321, 16, 0))
	nle(strconv.WriteInt64(con, -0x1234567887654321, 16, -20))
	nle(strconv.WriteInt64(con, -0x1234567887654321, 16, 20))
	nle(strconv.WriteInt64(con, -0x1234567887654321, -16, -20))
	nle(strconv.WriteInt64(con, -0x1234567887654321, -16, 20))

	nle(strconv.WriteFloat(con, 1.23e45, 'b', 24, 2, 32))
	nle(strconv.WriteFloat(con, 1.23e45, 'b', -24, 2, 32))
	nle(strconv.WriteFloat(con, 1.23e45, -'b', -24, 2, 32))
	nle(strconv.WriteFloat(con, 1.23e45, 'b', 24, 2, 64))
	nle(strconv.WriteFloat(con, 1.23e45, 'b', -24, 2, 64))
	nle(strconv.WriteFloat(con, 1.23e45, -'b', -24, 2, 64))
	nle(strconv.WriteFloat(con, -1.23e45, 'b', 24, 2, 64))
	nle(strconv.WriteFloat(con, -1.23e45, 'b', -24, 2, 64))
	nle(strconv.WriteFloat(con, -1.23e45, -'b', -24, 2, 64))

	nle(strconv.WriteFloat(con, 1.23456e9, 'f', 24, 4, 64))
	nle(strconv.WriteFloat(con, 1.23456e9, 'f', -24, 4, 64))
	nle(strconv.WriteFloat(con, 1.23456e9, -'f', -24, 4, 64))
	nle(strconv.WriteFloat(con, -1.23456e9, 'f', 24, 4, 64))
	nle(strconv.WriteFloat(con, -1.23456e9, 'f', -24, 4, 64))
	nle(strconv.WriteFloat(con, -1.23456e9, -'f', -24, 4, 64))

	nle(strconv.WriteFloat(con, 1.23456e-6, 'f', 24, 4, 64))
	nle(strconv.WriteFloat(con, 1.23456e-6, 'f', -24, 4, 64))
	nle(strconv.WriteFloat(con, 1.23456e-6, -'f', -24, 4, 64))
	nle(strconv.WriteFloat(con, -1.23456e-6, 'f', 24, 4, 64))
	nle(strconv.WriteFloat(con, -1.23456e-6, 'f', -24, 4, 64))
	nle(strconv.WriteFloat(con, -1.23456e-6, -'f', -24, 4, 64))

	nle(strconv.WriteFloat(con, 1.23456e-6, 'f', 24, 11, 64))
	nle(strconv.WriteFloat(con, 1.23456e-6, 'f', -24, 11, 64))
	nle(strconv.WriteFloat(con, 1.23456e-6, -'f', -24, 11, 64))
	nle(strconv.WriteFloat(con, -1.23456e-6, 'f', 24, 11, 64))
	nle(strconv.WriteFloat(con, -1.23456e-6, 'f', -24, 11, 64))
	nle(strconv.WriteFloat(con, -1.23456e-6, -'f', -24, 11, 64))

	nle(strconv.WriteFloat(con, 1.235e45, 'e', 24, 3, 64))
	nle(strconv.WriteFloat(con, 1.235e45, 'e', -24, 3, 64))
	nle(strconv.WriteFloat(con, 1.235e45, -'e', -24, 3, 64))
	nle(strconv.WriteFloat(con, -1.235e45, 'e', 24, 3, 64))
	nle(strconv.WriteFloat(con, -1.235e45, 'e', -24, 3, 64))
	nle(strconv.WriteFloat(con, -1.235e45, -'e', -24, 3, 64))

	nle(strconv.WriteFloat(con, 12340, 'g', 0, 4, 64))
	nle(strconv.WriteFloat(con, 1234, 'g', 0, 4, 64))
	nle(strconv.WriteFloat(con, 123.4, 'g', 0, 4, 64))
	nle(strconv.WriteFloat(con, 12.34, 'g', 0, 4, 64))
	nle(strconv.WriteFloat(con, 1.234, 'g', 0, 4, 64))
	nle(strconv.WriteFloat(con, 0.1234, 'g', 0, 4, 64))
	nle(strconv.WriteFloat(con, 0.01234, 'g', 0, 4, 64))
	nle(strconv.WriteFloat(con, 0.001234, 'g', 0, 4, 64))
	nle(strconv.WriteFloat(con, 0.0001234, 'g', 0, 4, 64))
	nle(strconv.WriteFloat(con, 0.00001234, 'g', 0, 4, 64))

	for i := 0; i < 20; i++ {
		const f = SmallestNonzeroFloat64
		nle(strconv.WriteFloat(con, f, 'e', -40, i, 64))
	}

	n, _ := fmt.Println("v =", -6.4321e-3)
	fmt.Println(n)

	n, _ = fmt.Println("cplx =", 3-3i)
	fmt.Println(n)

	x := []int{1, 3}
	n, _ = fmt.Println("ptr =", &x)
	fmt.Println(n)

	n, _ = fmt.Println("main =", main)
	fmt.Println(n)

	n, _ = fmt.Println("serial.Dev.Write =", (*serial.Dev).Write)
	fmt.Println(n)

	var c chan int
	n, _ = fmt.Println("chan =", c)
	fmt.Println(n)

	var b Bool
	n, _ = fmt.Println("b =", b)
	fmt.Println(n)

	b = true
	n, _ = fmt.Println("b =", b)
	fmt.Println(n)

	fmt.Println(
		"slice",
		unsafe.Sizeof([]byte(nil)), unsafe.Alignof([]byte(nil)),
	)

	siz, ali := slisiz()
	fmt.Println("C slice", siz, ali)

	fmt.Println(
		"int64",
		unsafe.Sizeof(int64(0)), unsafe.Alignof(int64(0)),
	)
	siz, ali = int64siz()
	fmt.Println("C int64", siz, ali)

	fmt.Println(
		"complex128",
		unsafe.Sizeof(complex128(0)), unsafe.Alignof(complex128(0)),
	)
	siz, ali = cplx128siz()
	fmt.Println("C complex128", siz, ali)

	type S struct {
		b []byte
		u uint32
	}
	fmt.Println(
		"S",
		unsafe.Sizeof(S{}), unsafe.Alignof(S{}),
	)

	a := [...]Bool{true, true, false, true}
	v := reflect.ValueOf(a)
	t := v.Type()
	fmt.Println(t, t.Size(), t.Align())
	fmt.Println(&a, a[:], a)
	fmt.Printf("Jeden\n")
	fmt.Printf("Type: %T Value: %v\n", a, a)
	i := -20
	fmt.Printf("%b %o %d %x\n", i, i, i, i)
	f := -1.234567890123456789e20
	fmt.Printf("|%-20.11g| |%20.10e|\n", f, f)
	var s = "STM32F411"
	fmt.Printf("|%15s|%015s|%-15s|%0-15s|\n", s, s, s, s)

	u := uint64(0xffff0000ffff000f)
	for i := uint(0); i <= 64; i++ {
		fmt.Println(bits.LeadingZeros64(u >> i))
	}
}

type II struct {
	c128 complex128
	u32  uint32
}

type SS struct {
	ii1 II
	ii2 II
	u   uint32
}

func slisiz() (uintptr, uintptr)
func int64siz() (uintptr, uintptr)
func cplx128siz() (uintptr, uintptr)
func interfacesiz() (uintptr, uintptr)
func sssize() uintptr
