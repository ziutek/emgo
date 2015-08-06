// This example shows how to use USART as serial console.
// It uses PA2, PA3 pins as Tx, Rx that are by default connected to ST-LINK
// Virtual Com Port.
package main

import (
	"delay"
	"fmt"
	"io"
	"reflect"
	"rtos"
	"strconv"
	"unsafe"

	"stm32/f4/gpio"
	"stm32/f4/irqs"
	"stm32/f4/periph"
	"stm32/f4/setup"
	"stm32/f4/usarts"
	"stm32/serial"
	"stm32/usart"
)

const (
	Green = 5
)

var (
	leds = gpio.A
	udev = usarts.USART2
	s    = serial.New(udev, 80, 8)
)

func init() {
	setup.Performance84(8)

	periph.AHB1ClockEnable(periph.GPIOA)
	periph.AHB1Reset(periph.GPIOA)

	leds.SetMode(Green, gpio.Out)

	periph.APB1ClockEnable(periph.USART2)
	periph.APB1Reset(periph.USART2)

	port, tx, rx := gpio.A, uint(2), uint(3)

	port.SetMode(tx, gpio.Alt)
	port.SetOutType(tx, gpio.PushPull)
	port.SetPull(tx, gpio.PullUp)
	port.SetOutSpeed(tx, gpio.Fast)
	port.SetAltFunc(tx, gpio.USART2)
	port.SetMode(rx, gpio.Alt)
	port.SetAltFunc(rx, gpio.USART2)

	udev.SetBaudRate(115200, setup.APB1Clk)
	udev.SetWordLen(usart.Bits8)
	udev.SetParity(usart.None)
	udev.SetStopBits(usart.Stop1b)
	udev.SetMode(usart.Tx | usart.Rx)
	udev.EnableIRQs(usart.RxNotEmptyIRQ)
	udev.Enable()

	rtos.IRQ(irqs.USART2).UseHandler(sirq)
	rtos.IRQ(irqs.USART2).Enable()

	s.SetUnix(true)

	fmt.DefaultWriter = s
}

func blink(c uint, d int) {
	leds.SetBit(c)
	if d > 0 {
		delay.Millisec(d)
	} else {
		delay.Loop(-1e4 * d)
	}
	leds.ClearBit(c)
}

func sirq() {
	s.IRQ()
}

func checkErr(err error) {
	if err != nil {
		s.WriteString("\nError: ")
		s.WriteString(err.Error())
		s.WriteByte('\n')
	}
}

type Bool bool

func (b Bool) Format(st fmt.State, c rune) {
	s := "fałsz"
	if b {
		s = "prawda"
	}
	io.WriteString(st, s)
}

func main() {
	const (
		SmallestNormal         = 2.2250738585072014e-308
		SmallestNonzeroFloat64 = 4.940656458412465441765687928682213723651e-324
	)

	var buf [40]byte

	strconv.FormatBool(buf[:], true, -2)
	s.Write(buf[:])
	s.WriteByte('\n')

	strconv.FormatInt(buf[:], 123456, -10)
	s.Write(buf[:])
	s.WriteByte('\n')

	strconv.FormatFloat(buf[:], 0, -'e', 0)
	s.Write(buf[:])
	s.WriteByte('\n')

	fmt.Println()

	for i := 0; i < 20; i++ {
		strconv.FormatFloat(buf[:], SmallestNonzeroFloat64, -'e', i)
		s.Write(buf[:])
		s.WriteByte('\n')
	}

	var n int
	n, _ = fmt.Println("b =", true, "a =", 12)
	fmt.Println(n)

	n, _ = fmt.Println("v =", -6.4321e-3)
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
	// BUG: this doesn't compile
	//fmt.Println(slisiz())
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
	fmt.Println(&a, a[:])
}

func slisiz() (uintptr, uintptr)
func int64siz() (uintptr, uintptr)
func cplx128siz() (uintptr, uintptr)
