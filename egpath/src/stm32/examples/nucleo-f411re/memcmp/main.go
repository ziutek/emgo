package main

import (
	"bytes"
	"delay"
	"fmt"

	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
)

func init() {
	system.Setup96(8)
	systick.Setup()
}

// emgo
var a1 = [...]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}

var a2 [len(a1)]byte

func main() {
	delay.Millisec(250) // Wait for OpenOCD (press reset if you see nothing).

	copy(a2[:0], a1[:0])
	copy(a2[:1], a1[:1])
	copy(a2[:2], a1[:2])
	copy(a2[:3], a1[:3])
	copy(a2[:4], a1[:4])
	a2 = a1

	copy(a2[1:], a1[1:])
	a2[11] = 0

	fmt.Println(a2[:])
	fmt.Println(bytes.Equal(a1[1:], a2[1:]))

	s1 := "babaa"
	s2 := "abbaa"
	fmt.Printf("s1='%s' s2='%s\n", s1, s2)
	fmt.Printf("s1 < s2:  %t\n", s1 < s2)
	fmt.Printf("s1 <= s2: %t\n", s1 <= s2)
	fmt.Printf("s1 == s2: %t\n", s1 == s2)
	fmt.Printf("s1 != s2: %t\n", s1 != s2)
	fmt.Printf("s1 >= s2: %t\n", s1 >= s2)
	fmt.Printf("s1 > s2:  %t\n", s1 > s2)
}
