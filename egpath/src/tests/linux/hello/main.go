package main

import (
	"fmt"
	"os"
)

func checkErr(err error) {
	if err != nil {
		os.Stderr.WriteString("Fatal error: ")
		os.Stderr.WriteString(err.Error())
		os.Stderr.WriteString(".\n")
		os.Exit(1)
	}
}

func strlen(s *[2<<31 - 1]byte) int {
	for n, c := range s {
		if c == 0 {
			return n
		}
	}
	panic("strlen overflow")
}

func main() {
	fmt.Println("Args:")
	for i, a := range os.Args {
		fmt.Printf("%d: %s\n", i, a)
	}

	fmt.Println("Env:")
	for i, e := range os.Env {
		fmt.Printf("%d: %s\n", i, e)
	}

	buf := make([]byte, 80)
	n := copy(buf, os.Args[0])
	fmt.Printf("Program name: %s\n", buf[:n])

	n, err := os.Stdin.Read(buf)
	checkErr(err)

	f, err := os.OpenFile(
		"file.txt",
		os.O_CREATE|os.O_RDWR|os.O_EXCL,
		0660,
	)
	checkErr(err)

	n, err = f.Write(buf[:n])
	checkErr(err)

	f.Close()
	fmt.Printf("%d bytes written.\n", n)

	/*
		sd, e := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
		checkErrno(e)
		e = syscall.SetsockoptInt(sd, syscall.SOL_SOCKET, syscall.SO_RCVBUF, 524288)
		checkErrno(e)

		sa := syscall.RawSockaddrInet4{
			Family: syscall.AF_INET,
			Port:   0xd204,
		}
		e = syscall.Bind(sd, &sa)
		checkErrno(e)

		var buf [2048]byte
		for {
			_, e := syscall.Read(sd, buf[:])
			checkErrno(e)
			syscall.WriteString(1, "udp pkt\n")
		}
	*/
}
