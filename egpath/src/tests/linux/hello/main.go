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

func main() {

	s := "Hello world!"
	fmt.Println(s)

	buf := make([]byte, 80)
	n := copy(buf, s)
	fmt.Printf("%s\n", buf[:n])

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
