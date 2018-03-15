package main

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
)

func init() {
	system.Setup168(8)
	systick.Setup(2e6)
	initConsole()
}

func main() {
	t := time.Now()
	fmt.Println(t)
	fmt.Println(true, false)
	fmt.Println(10, -10, 1234567890, -123456789)
	fmt.Println(int64(1234567890123), int64(-1234567890123))
	fmt.Println(123.456e-20, -123.456e2)

	fmt.Printf("|%11s|\n", "abc")
	fmt.Printf("|%011s|\n", "abc")
	fmt.Printf("|%-11s|\n", "abc")
	fmt.Printf("|%-011s|\n", "abc")
	fmt.Printf("|%11d|\n", 123)
	fmt.Printf("|%011d|\n", 123)
	fmt.Printf("|%-11d|\n", 123)
	fmt.Printf("|%-011d|\n", 123)
	fmt.Printf("|%11x|\n", 123)
	fmt.Printf("|%011x|\n", 123)
	fmt.Printf("|%-11x|\n", 123)
	fmt.Printf("|%-011x|\n", 123)
	fmt.Printf("|%11X|\n", 123)
	fmt.Printf("|%011X|\n", 123)
	fmt.Printf("|%-11X|\n", 123)
	fmt.Printf("|%-011X|\n", 123)
	fmt.Printf("|%11.2f|\n", 12.499)
	fmt.Printf("|%011.2f|\n", 12.499)
	fmt.Printf("|%-11.2f|\n", 12.499)
	fmt.Printf("|%-011.2f|\n", 12.499)

	w := fmt.DefaultWriter
	io.WriteString(w, "\n|")
	strconv.WriteString(w, "ABC", 11, ' ')
	io.WriteString(w, "|\n|")
	strconv.WriteString(w, "ABC", -11, ' ')
	io.WriteString(w, "|\n|")
	strconv.WriteString(w, "ABC", 11, '0')
	io.WriteString(w, "|\n|")
	strconv.WriteString(w, "ABC", -11, '0')
	io.WriteString(w, "|\n|")
	strconv.WriteString(w, "ABC", 11, '.')
	io.WriteString(w, "|\n|")
	strconv.WriteString(w, "ABC", -11, '.')
	io.WriteString(w, "|\n|")

	strconv.WriteInt(w, 456, 10, 11, ' ')
	io.WriteString(w, "|\n|")
	strconv.WriteInt(w, 456, 10, -11, ' ')
	io.WriteString(w, "|\n|")
	strconv.WriteInt(w, 456, 10, 11, '_')
	io.WriteString(w, "|\n|")
	strconv.WriteInt(w, 456, 10, -11, '_')
	io.WriteString(w, "|\n|")

}
