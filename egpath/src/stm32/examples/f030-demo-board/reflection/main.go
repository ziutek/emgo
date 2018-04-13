package main

import (
	"debug/semihosting"
	"reflect"
	"strconv"

	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
)

var stdout semihosting.File

func init() {
	system.SetupPLL(8, 1, 48/8)
	systick.Setup(2e6)

	var err error
	stdout, err = semihosting.OpenFile(":tt", semihosting.W)
	for err != nil {
	}
}

type stringer interface {
	String() string
}

func println(args ...interface{}) {
	for i, a := range args {
		if i > 0 {
			stdout.WriteString(" ")
		}
		switch v := a.(type) {
		case string:
			stdout.WriteString(v)
		case int:
			strconv.WriteInt(stdout, v, 10, 0, 0)
		case bool:
			strconv.WriteBool(stdout, v, 't', 0, 0)
		case stringer:
			stdout.WriteString(v.String())
		default:
			stdout.WriteString("%unknown")
		}
	}
	stdout.WriteString("\r\n")
}

type S struct {
	A int
	B bool
}

func main() {
	p := &S{-123, true}

	v := reflect.ValueOf(p)

	println("kind(p) =", v.Kind())
	println("kind(*p) =", v.Elem().Kind())
	println("type(*p) =", v.Elem().Type().Name())

	v = v.Elem()

	println("*p = {")
	for i := 0; i < v.NumField(); i++ {
		ft := v.Type().Field(i)
		fv := v.Field(i)
		println("  ", ft.Name(), ":", fv.Interface())
	}
	println("}")
}
